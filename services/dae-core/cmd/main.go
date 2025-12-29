package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	libconfigs "github.com/deni12345/dae-services/libs/configs"
	corev1 "github.com/deni12345/dae-services/proto/gen"
	"github.com/deni12345/dae-services/services/dae-core/internal/app/health"
	"github.com/deni12345/dae-services/services/dae-core/internal/app/order"
	"github.com/deni12345/dae-services/services/dae-core/internal/app/sheet"
	"github.com/deni12345/dae-services/services/dae-core/internal/app/user"
	"github.com/deni12345/dae-services/services/dae-core/internal/configs"
	grpchandler "github.com/deni12345/dae-services/services/dae-core/internal/grpc"
	"github.com/deni12345/dae-services/services/dae-core/internal/grpc/interceptor"
	frstore "github.com/deni12345/dae-services/services/dae-core/internal/infra/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/infra/observability"
	infraredis "github.com/deni12345/dae-services/services/dae-core/internal/infra/redis"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	redisgo "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var config = configs.Value{}
	if err := libconfigs.LoadWithEnvOptions(ctx, &config, "ENVIRONMENT", libconfigs.WithYamlFile("configs.yml")); err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	// Initialize logger if enabled
	var logShutdown func(context.Context) error
	if config.EnableLogging {
		var err error
		logShutdown, err = observability.InitLogger(ctx, observability.LoggerConfig{
			Level:        config.LogLevel,
			ServiceName:  config.ServiceName,
			OTLPEndpoint: config.OtelCol,
			Insecure:     config.Insecure,
		})
		if err != nil {
			observability.Fatal(ctx, "failed to initialize logger", "error", err)
		}
		defer func() { _ = logShutdown(ctx) }()
		slog.InfoContext(ctx, "logger initialized", "environment", config.Environment, "level", config.LogLevel)
	}

	slog.InfoContext(ctx, "logger initialized", "environment", config.Environment, "level", config.LogLevel)

	fsClient, err := frstore.NewFirestoreClient(ctx, config.FirestoreProjectID)
	if err != nil {
		observability.Fatal(ctx, "failed to initialize firestore", "error", err)
	}
	defer func() { _ = fsClient.Close() }()

	userRepo, orderRepo, sheetRepo := initRepos(fsClient, config)

	redisClient, idemStore := initIdempotencyStore(ctx, config)
	defer func() {
		if redisClient != nil {
			_ = redisClient.Close()
		}
	}()

	var traceShutdown, metricShutdown func(context.Context) error
	var metrics *observability.Metrics

	traceShutdown, metrics, metricShutdown, err = initObservability(ctx, config)
	if err != nil {
		observability.Fatal(ctx, "failed to initialize observability", "error", err)
	}

	userUC := user.NewUsecase(userRepo)
	orderUC := order.NewUsecase(orderRepo, sheetRepo, idemStore)
	sheetUC := sheet.NewUsecase(sheetRepo, idemStore)
	healthUC := health.NewUsecase(fsClient, redisClient)

	grpcServer := createGRPCServer(metrics, userUC, orderUC, sheetUC, healthUC)
	_, err = startGRPCServer(grpcServer, config.GRPCAddress)
	if err != nil {
		observability.Fatal(ctx, "failed to start gRPC server", "error", err)
	}

	// Setup signal handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		slog.Info("shutdown signal received, initiating graceful shutdown")
		cancel()
	}()

	// Wait on context cancel, then shutdown with timeout
	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	slog.Info("shutting down gracefully...")
	if err := traceShutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown tracer", "error", err)
	}
	if err := metricShutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown metrics", "error", err)
	}
	grpcServer.GracefulStop()
	slog.Info("server stopped")
}

func initRepos(fsClient *firestore.Client, cfg configs.Value) (port.UsersRepo, port.OrdersRepo, port.SheetRepo) {
	userRepo := frstore.NewUserRepo(fsClient, cfg.PageSize)
	orderRepo := frstore.NewOrderRepo(fsClient, cfg.PageSize)
	sheetRepo := frstore.NewSheetRepo(fsClient, cfg.PageSize)
	return userRepo, orderRepo, sheetRepo
}

func initIdempotencyStore(ctx context.Context, cfg configs.Value) (*redisgo.Client, port.IdempotencyStore) {
	redisClient := infraredis.NewRedisClient(ctx, cfg)
	idemStore := infraredis.NewIdempotencyStore(redisClient)
	return redisClient, idemStore
}

func initObservability(ctx context.Context, cfg configs.Value) (func(context.Context) error, *observability.Metrics, func(context.Context) error, error) {
	var traceShutdown func(context.Context) error
	var metricShutdown func(context.Context) error
	var metrics *observability.Metrics

	// Initialize tracer if enabled
	if cfg.EnableTracing {
		var err error
		traceShutdown, err = observability.InitTracer(ctx, observability.TracerConfig{
			ServiceName:  cfg.ServiceName,
			OTLPEndpoint: cfg.OtelCol,
			Insecure:     cfg.Insecure,
		})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to initialize tracer: %w", err)
		}
		slog.InfoContext(ctx, "tracer initialized")
	} else {
		traceShutdown = func(context.Context) error { return nil }
		slog.InfoContext(ctx, "tracer disabled via config")
	}

	// Initialize metrics if enabled
	if cfg.EnableMetrics {
		var err error
		metrics, metricShutdown, err = observability.InitMeter(ctx, observability.MeterConfig{
			ServiceName:  cfg.ServiceName,
			OTLPEndpoint: cfg.OtelCol,
			Insecure:     cfg.Insecure,
		})
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to initialize metrics: %w", err)
		}
		slog.InfoContext(ctx, "metrics initialized")
	} else {
		metricShutdown = func(context.Context) error { return nil }
		slog.InfoContext(ctx, "metrics disabled via config")
	}

	return traceShutdown, metrics, metricShutdown, nil
}

func createGRPCServer(
	metrics *observability.Metrics,
	userUC user.Usecase,
	orderUC order.Usecase,
	sheetUC sheet.Usecase,
	healthUC health.Usecase,
) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.IdemInterceptor(),
			interceptor.MetricsInterceptor(metrics),
			interceptor.ValidateRequestInterceptor(metrics),
			interceptor.LoggingInterceptor(),
		),
	)

	reflection.Register(grpcServer)
	corev1.RegisterUsersServiceServer(grpcServer, grpchandler.NewUserHandler(userUC))
	corev1.RegisterOrdersServiceServer(grpcServer, grpchandler.NewOrderHandler(orderUC))
	corev1.RegisterSheetsServiceServer(grpcServer, grpchandler.NewSheetHandler(sheetUC))
	corev1.RegisterHealthServiceServer(grpcServer, grpchandler.NewHealthHandler(healthUC))
	return grpcServer
}

func startGRPCServer(srv *grpc.Server, addr string) (net.Listener, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	go func() {
		slog.Info("gRPC server started", "address", addr)
		if err := srv.Serve(lis); err != nil {
			observability.Fatal(context.Background(), "grpc server error", "error", err)
		}
	}()
	return lis, nil
}
