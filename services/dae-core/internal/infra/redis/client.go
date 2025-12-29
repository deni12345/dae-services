package redis

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/configs"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg configs.Value) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
}
