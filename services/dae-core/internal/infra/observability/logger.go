package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/deni12345/dae-services/libs/prettylog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type LoggerConfig struct {
	Level        string
	ServiceName  string
	OTLPEndpoint string
	Insecure     bool
}

func InitLogger(ctx context.Context, cfg LoggerConfig) (func(context.Context) error, error) {
	level := parseLogLevel(cfg.Level)

	if cfg.OTLPEndpoint == "" {
		return nil, fmt.Errorf("OTLP endpoint is required for logger")
	}
	exporterOpts := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(cfg.OTLPEndpoint),
	}
	if cfg.Insecure {
		exporterOpts = append(exporterOpts, otlploggrpc.WithInsecure())
	}

	exporter, err := otlploggrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP log exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create log resource: %w", err)
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(exporter)),
		log.WithResource(res),
	)

	global.SetLoggerProvider(loggerProvider)
	shutdown := loggerProvider.Shutdown
	otelLogger := loggerProvider.Logger("otel-logger")

	handler := prettylog.NewPrettyHandler(
		os.Stdout,
		otelLogger,
		&slog.HandlerOptions{Level: level},
	)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return shutdown, nil
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}
