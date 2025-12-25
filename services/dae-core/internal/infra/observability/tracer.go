package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type TracerConfig struct {
	ServiceName  string
	OTLPEndpoint string
	Insecure     bool
}

func InitTracer(ctx context.Context, cfg TracerConfig) (func(context.Context) error, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	exporterOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
	}
	if cfg.Insecure {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create trace resource: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider.Shutdown, nil
}

func (cfg *TracerConfig) validate() error {
	if cfg.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if cfg.OTLPEndpoint == "" {
		return fmt.Errorf("OTLP endpoint is required")
	}
	return nil
}
