package redis

import (
	"context"

	"github.com/deni12345/dae-core/internal/configs"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     configs.Values.Redis.Addr,
		Password: configs.Values.Redis.Password,
	})
}
