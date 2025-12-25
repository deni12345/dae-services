package interceptor

import (
	"context"
	"time"

	"github.com/deni12345/dae-services/services/dae-core/internal/infra/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor(metrics *observability.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start).Seconds()

		// no metrics registered
		if metrics == nil {
			return resp, err
		}

		st := status.Code(err)

		attributes := []attribute.KeyValue{
			attribute.String("rpc.system", "grpc"),
			attribute.String("rpc.method", info.FullMethod),
			attribute.String("rpc.grpc.status_code", st.String()),
		}

		metrics.RPCCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
		metrics.RPCLatency.Record(ctx, elapsed, metric.WithAttributes(attributes...))

		return resp, err
	}
}
