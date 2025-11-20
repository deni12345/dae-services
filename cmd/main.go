package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/deni12345/dae-core/common/prettylog"
	"github.com/deni12345/dae-core/internal/app/order"
	"github.com/deni12345/dae-core/internal/app/sheet"
	"github.com/deni12345/dae-core/internal/app/user"
	"github.com/deni12345/dae-core/internal/configs"
	grpchandler "github.com/deni12345/dae-core/internal/grpc"
	"github.com/deni12345/dae-core/internal/infra/firestore"
	"github.com/deni12345/dae-core/internal/infra/redis"
	corev1 "github.com/deni12345/dae-core/proto/gen"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()
	handler := prettylog.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Load config
	if err := configs.LoadWithRetry(ctx,
		configs.WithYamlFile("configs.yml"),
	); err != nil {
		log.Fatalf("load config error: %v", err)
	}

	// Initialize Firestore
	fsClient, err := firestore.NewFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("failed to init firestore: %v", err)
	}
	defer fsClient.Close()

	// Repositories
	userRepo := firestore.NewUserRepo(fsClient)
	orderRepo := firestore.NewOrderRepo(fsClient, configs.Values.Firestore.ProjectID)
	sheetRepo := firestore.NewSheetRepo(fsClient)

	// Idempotency store (stub for now - need Redis client)
	// TODO: Initialize Redis client properly
	redisClient := redis.NewRedisClient(ctx)
	var idemStore = redis.NewIdempotencyStore(redisClient)

	// Usecases
	userUC := user.NewUsecase(userRepo)
	orderUC := order.NewUsecase(orderRepo, sheetRepo, idemStore)
	sheetUC := sheet.NewUsecase(sheetRepo, idemStore)

	// gRPC server setup
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	corev1.RegisterUsersServiceServer(grpcServer, grpchandler.NewUserHandler(userUC))
	corev1.RegisterOrdersServiceServer(grpcServer, grpchandler.NewOrderHandler(orderUC))
	corev1.RegisterSheetsServiceServer(grpcServer, grpchandler.NewSheetHandler(sheetUC))

	// HTTP health check
	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			slog.Error("failed to write health response: %v", err)
		}
	})

	// Start HTTP health server
	go func() {
		healthAddr := fmt.Sprintf(":%d", configs.Values.Port)
		slog.Info("health server running at", "addr", healthAddr)
		if err := http.ListenAndServe(healthAddr, r); err != nil {
			prettylog.Fatal(ctx, "health server error", "error", err)
		}
	}()

	// Start gRPC server
	grpcAddr := ":50051"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		prettylog.Fatal(ctx, "failed to listen", "addr", grpcAddr, "error", err)
	}

	slog.Info("gRPC server running", slog.Group("server",
		slog.String("component", "main"),
		"port", grpcAddr,
	))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			prettylog.Fatal(ctx, "grpc server error", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down gracefully...")
	grpcServer.GracefulStop()
	slog.Info("server stopped")
}
