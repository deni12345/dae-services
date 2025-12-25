package interceptor

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor logs all gRPC requests and responses
func LoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call handler
		resp, err := handler(ctx, req)

		// Log after completion
		code := codes.OK
		if err != nil {
			st, _ := status.FromError(err)
			code = st.Code()
		}

		traceID := ""
		if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
			traceID = sc.TraceID().String()
		}

		idemKey := idempotencyKeyFromContext(ctx)

		slog.Info("grpc_request",
			"method", info.FullMethod,
			"code", code.String(),
			"duration_ms", time.Since(start).Milliseconds(),
			"trace_id", traceID,
			"idempotency_key", idemKey,
		)

		return resp, err
	}
}
