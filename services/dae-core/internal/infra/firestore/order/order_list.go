package order

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
)

func (r *orderRepo) List(ctx context.Context, query port.ListOrdersQuery) ([]*domain.Order, error) {
	ctx, span := tracer.Start(ctx, "OrderRepo.List")
	defer span.End()

	limit := query.Limit
	if limit <= 0 || limit > 1000 {
		limit = r.defaultPageSize
	}
	q := r.collection.Query
	q = q.Limit(int(limit))

	// If cursor is provided, start after that document
	if query.Cursor != "" {
		cursorSnap, err := r.collection.Doc(query.Cursor).Get(ctx)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("get cursor document: %w", err)
		}
		q = q.StartAfter(cursorSnap)
	}

	// Add ordering (most recent first)
	q = q.OrderBy("created_at", firestore.Desc)

	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("list orders: %w", err)
	}

	orders := make([]*domain.Order, 0, len(docs))
	for _, doc := range docs {
		var order domain.Order
		if err := doc.DataTo(&order); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("unmarshal order: %w", err)
		}
		if order.ID == "" {
			order.ID = doc.Ref.ID
		}
		orders = append(orders, &order)
	}

	return orders, nil
}
