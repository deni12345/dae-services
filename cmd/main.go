package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/deni12345/dae-core/common/logx"
	"github.com/deni12345/dae-core/internal/configs"
	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	logx.NewLoggerJSON(os.Stdout)

	if err := configs.LoadWithRetry(ctx,
		configs.WithYamlFile("configs.yml"),
	); err != nil {
		log.Fatalf("load config error: %v", err)
	}

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from dae service"))
		})
	})

	logx.Info("server running at port: %v", configs.Values.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", configs.Values.Port), r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
