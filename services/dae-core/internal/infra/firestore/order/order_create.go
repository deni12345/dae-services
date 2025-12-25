package order

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *orderRepo) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	ctx, span := tracer.Start(ctx, "OrderRepo.Create")
	defer span.End()

	if order.ID == "" {
		err := fmt.Errorf("order ID is required")
		span.RecordError(err)
		return nil, err
	}

	docRef := r.collection.Doc(order.ID)
	_, err := docRef.Create(ctx, order)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			span.RecordError(ErrOrderExists)
			return nil, ErrOrderExists
		}
		span.RecordError(err)
		return nil, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}
