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
	corev1 "github.com/deni12345/dae-services/proto/gen"
	redisgo "github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := configs.LoadWithRetry(ctx, configs.WithYamlFile("configs.yml"))
	if err != nil {
		slog.Error("load config failed", "error", err)
		os.Exit(1)
	}

	logShutdown, err := observability.InitLogger(ctx, observability.LoggerConfig{
		Level:        cfg.LogLevel,
		ServiceName:  cfg.ServiceName,
		OTLPEndpoint: cfg.OtelCol,
		Insecure:     cfg.Insecure,
	})
	if err != nil {
		observability.Fatal(ctx, "failed to initialize logger", "error", err)
	}
	defer func() { _ = logShutdown(ctx) }()

	slog.InfoContext(ctx, "logger initialized", "environment", cfg.Environment, "level", cfg.LogLevel)

	fsClient, err := frstore.NewFirestoreClient(ctx, cfg.Firestore.ProjectID)
	if err != nil {
		observability.Fatal(ctx, "failed to initialize firestore", "error", err)
	}
	defer func() { _ = fsClient.Close() }()

	userRepo, orderRepo, sheetRepo := initRepos(fsClient, cfg)

	redisClient, idemStore := initIdempotencyStore(ctx, cfg)
	defer func() {
		if redisClient != nil {
			_ = redisClient.Close()
		}
	}()

	traceShutdown, metrics, metricShutdown, err := initObservability(ctx, cfg)
	if err != nil {
		observability.Fatal(ctx, "failed to initialize observability", "error", err)
	}

	userUC := user.NewUsecase(userRepo)
	orderUC := order.NewUsecase(orderRepo, sheetRepo, idemStore)
	sheetUC := sheet.NewUsecase(sheetRepo, idemStore)
	healthUC := health.NewUsecase(fsClient, redisClient)

	grpcServer := createGRPCServer(metrics, userUC, orderUC, sheetUC, healthUC)
	_, err = startGRPCServer(grpcServer, cfg.GRPCAddress)
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
	traceShutdown, err := observability.InitTracer(ctx, observability.TracerConfig{
		ServiceName:  cfg.ServiceName,
		OTLPEndpoint: cfg.OtelCol,
		Insecure:     cfg.Insecure,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	metrics, metricShutdown, err := observability.InitMeter(ctx, observability.MeterConfig{
		ServiceName:  cfg.ServiceName,
		OTLPEndpoint: cfg.OtelCol,
		Insecure:     cfg.Insecure,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize metrics: %w", err)
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
