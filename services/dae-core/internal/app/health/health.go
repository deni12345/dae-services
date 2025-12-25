package health

import (
	"context"
	"fmt"
	"time"

	"sync"

	"cloud.google.com/go/firestore"
	redisgo "github.com/redis/go-redis/v9"
)

// Usecase defines health check behavior.
type Usecase interface {
	// Check returns nil when service is healthy or an error describing issues.
	Check(ctx context.Context) error
}

type healthUC struct {
	fs      *firestore.Client
	rdb     *redisgo.Client
	timeout time.Duration
}

// NewUsecase creates a health usecase that checks Firestore and Redis connectivity.
func NewUsecase(fs *firestore.Client, rdb *redisgo.Client) Usecase {
	return &healthUC{fs: fs, rdb: rdb, timeout: 2 * time.Second}
}

func (h *healthUC) Check(ctx context.Context) error {
	// use a short timeout for health checks
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	if h.fs != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := h.fs.Collection("__health_check").Doc("_ping").Get(ctx)
			if err != nil {
				errCh <- fmt.Errorf("firestore: %w", err)
			}
		}()
	}

	if h.rdb != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := h.rdb.Ping(ctx).Result(); err != nil {
				errCh <- fmt.Errorf("redis: %w", err)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return err
	}
	return nil
}
