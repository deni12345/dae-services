package order

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
)

func (r *orderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	ctx, span := tracer.Start(ctx, "OrderRepo.GetByID")
	defer span.End()

	doc, err := r.collection.Doc(id).Get(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("get order by id: %w", err)
	}
	var order domain.Order
	if err := doc.DataTo(&order); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("data to order: %w", err)
	}
	if order.ID == "" {
		order.ID = doc.Ref.ID
	}
	return &order, nil
}
