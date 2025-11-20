package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const MDKeyIdem = "idempotency-key"

func IdemInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		var key string
		if v := md.Get(MDKeyIdem); len(v) > 0 {
			key = strings.TrimSpace(v[0])
		}
		if key == "" {
			return nil, status.Error(codes.InvalidArgument, "missing idempotency key")
		}
		ctx = context.WithValue(ctx, MDKeyIdem, key)
		return handler(ctx, req)
	}
}
