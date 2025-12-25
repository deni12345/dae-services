package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Metrics struct {
	RPCCounter         api.Int64Counter
	RPCLatency         api.Float64Histogram
	ValidationFailures api.Int64Counter
}

type MeterConfig struct {
	ServiceName  string
	OTLPEndpoint string
	Insecure     bool
}

func InitMeter(ctx context.Context, cfg MeterConfig) (*Metrics, func(context.Context) error, error) {
	if err := cfg.validate(); err != nil {
		return nil, nil, err
	}

	exporterOpts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint),
	}
	if cfg.Insecure {
		exporterOpts = append(exporterOpts, otlpmetricgrpc.WithInsecure())
	}

	exporter, err := otlpmetricgrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, nil, fmt.Errorf("create OTLP metric exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create metric resource: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)
	meter := meterProvider.Meter(cfg.ServiceName)

	rpcCounter, err := meter.Int64Counter(
		"rpc_server_requests_total",
		api.WithDescription("Total number of RPC requests received"))
	if err != nil {
		return nil, nil, fmt.Errorf("create rpc counter: %w", err)
	}

	rpcLatency, err := meter.Float64Histogram(
		"rpc_server_request_duration_seconds",
		api.WithDescription("RPC server request duration in seconds"),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("create rpc latency histogram: %w", err)
	}

	validationFailures, err := meter.Int64Counter(
		"rpc_server_validation_failures_total",
		api.WithDescription("Total number of RPC validation failures"))
	if err != nil {
		return nil, nil, fmt.Errorf("create validation failures counter: %w", err)
	}

	metrics := &Metrics{
		RPCCounter:         rpcCounter,
		RPCLatency:         rpcLatency,
		ValidationFailures: validationFailures,
	}

	return metrics, meterProvider.Shutdown, nil
}

func (cfg *MeterConfig) validate() error {
	if cfg.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if cfg.OTLPEndpoint == "" {
		return fmt.Errorf("OTLP endpoint is required")
	}
	return nil
}
