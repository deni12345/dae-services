package order

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *orderRepo) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if order.ID == "" {
		return nil, fmt.Errorf("order ID is required")
	}

	docRef := r.collection.Doc(order.ID)
	_, err := docRef.Create(ctx, order)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, ErrOrderExists
		}
		return nil, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}
