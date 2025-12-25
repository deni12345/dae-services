package interceptor

import (
	"context"
	"fmt"
	"path"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	mdKey                    = "idempotency-key"
	IdemKey       contextKey = "idempotency-key"
	IdemMethodKey contextKey = "idempotency-method"
)

func IdemInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, _ := metadata.FromIncomingContext(ctx)
		var key string

		if v := md.Get(mdKey); len(v) > 0 {
			key = strings.TrimSpace(v[0])
		}

		// Only enforce idempotency for write-like RPCs
		methodName := path.Base(info.FullMethod)

		ctx = context.WithValue(ctx, IdemMethodKey, methodName)
		if isWriteMethod(methodName) {
			if key == "" {
				return nil, status.Error(codes.InvalidArgument, "missing idempotency key")
			}
			ctx = context.WithValue(ctx, IdemKey, key)
		}

		return handler(ctx, req)
	}
}

// idempotencyKeyFromContext returns idempotency key stored in context, or empty string when absent.
func idempotencyKeyFromContext(ctx context.Context) string {
	if v := ctx.Value(IdemKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// methodNameFromContext returns the gRPC method name stored in context.
func methodNameFromContext(ctx context.Context) string {
	if v := ctx.Value(IdemMethodKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetOrCreateIdempotencyKeyWithHash returns the idempotency key if present in context; otherwise it builds one from method name, the provided parts, and payload hash.
func GetOrCreateIdempotencyKeyWithHash(ctx context.Context, payloadHash string, parts ...string) string {
	if k := idempotencyKeyFromContext(ctx); k != "" {
		return k
	}
	method := methodNameFromContext(ctx)
	if method == "" {
		method = "unknown"
	}
	if len(parts) == 0 && payloadHash == "" {
		return method
	}

	joined := ""
	if len(parts) > 0 {
		joined = parts[0]
		for _, p := range parts[1:] {
			joined = fmt.Sprintf("%s:%s", joined, p)
		}
	}

	if joined == "" {
		return fmt.Sprintf("%s:%s", method, payloadHash)
	}
	if payloadHash == "" {
		return fmt.Sprintf("%s:%s", method, joined)
	}
	return fmt.Sprintf("%s:%s:%s", method, joined, payloadHash)
}

// isWriteMethod determines whether a gRPC method should require idempotency key.
func isWriteMethod(methodName string) bool {
	// Simple heuristics: if name starts with or contains these prefixes.
	prefixes := []string{"Create", "Update", "Delete", "Set", "AdminSet", "Close", "Reopen", "Join", "Leave"}
	for _, p := range prefixes {
		if strings.HasPrefix(methodName, p) || strings.Contains(methodName, p) {
			return true
		}
	}
	return false
}
