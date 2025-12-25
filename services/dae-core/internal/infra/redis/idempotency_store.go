package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var luaUnlockIfMatch = redis.NewScript(`
if redis.call("Get", KEYS[1]) == ARGV[1] then
	return redis.call("Del", KEYS[1])
else
	return 0
end
`)

func unlockIfMatch(ctx context.Context, rdb *redis.Client, processingKey, taskID string) error {
	_, err := luaUnlockIfMatch.Run(ctx, rdb, []string{processingKey}, taskID).Result()
	return err
}

type idempotencyStore struct {
	client *redis.Client
}

func NewIdempotencyStore(client *redis.Client) port.IdempotencyStore {
	return &idempotencyStore{
		client: client,
	}
}

func (s *idempotencyStore) Do(
	ctx context.Context,
	key string,
	ttl time.Duration,
	fn func(ctx context.Context) ([]byte, error),
) ([]byte, error) {
	taskID := uuid.NewString()
	doneKey := fmt.Sprintf("idem:done:%s", key)

	// Fast path: check if already done
	if val, err := s.client.Get(ctx, doneKey).Result(); err == nil {
		return []byte(val), nil
	}

	// Try to acquire lock with SETNX
	processingKey := fmt.Sprintf("idem:processing:%s", key)
	acquired, err := s.client.SetNX(ctx, processingKey, taskID, ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("redis setnx: %w", err)
	}

	if !acquired {
		// Still processing — wait up to maxWait duration before returning error to caller
		backOff := 100 * time.Millisecond
		maxWait := 30 * time.Second
		waited := time.Duration(0)
		for {
			if val, err := s.client.Get(ctx, doneKey).Result(); err == nil {
				return []byte(val), nil
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backOff):
				waited += backOff
				if waited >= maxWait {
					// Give up waiting to avoid infinitely blocking the client
					return nil, fmt.Errorf("idem:processing:timeout")
				}
				if backOff < 500*time.Millisecond {
					backOff *= 2
				}
			}
		}
	}

	closeChan := make(chan struct{})
	defer close(closeChan)
	go func() {
		ticker := time.NewTicker(ttl / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-closeChan:
				return
			case <-ticker.C:
				// Extend processing key TTL
				if v, err := s.client.Get(ctx, processingKey).Result(); err == nil && v == taskID {
					s.client.PExpire(ctx, processingKey, ttl)
				}
			}
		}
	}()

	// Execute function
	result, err := fn(ctx)
	if err != nil {
		unlockIfMatch(ctx, s.client, processingKey, taskID)
		return nil, err
	}

	// Store result
	if err := s.client.Set(ctx, doneKey, result, 10*ttl).Err(); err != nil {
		unlockIfMatch(ctx, s.client, processingKey, taskID)
		return nil, fmt.Errorf("redis set done: %w", err)
	}

	unlockIfMatch(ctx, s.client, processingKey, taskID)
	return result, nil
}
