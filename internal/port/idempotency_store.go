package port

import (
	"context"
	"time"
)

type IdempotencyStore interface {
	Do(ctx context.Context, key string, ttlSeconds time.Duration,
		fn func(ctx context.Context) ([]byte, error)) ([]byte, error)
}
