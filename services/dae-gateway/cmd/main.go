package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	server := http.Server{
		Addr:    ":8084",
		Handler: r,
	}

	go func() {
		slog.Info("gateway listening on :8084")
		if err := server.ListenAndServe(); err != nil {
			slog.Error("gateway server failed", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
		slog.Error("failed to shutdown gateway server", "error", err)
	}

	slog.Info("gateway server shutdown complete")
}
