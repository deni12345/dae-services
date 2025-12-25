package order

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
)

// GetOrderByID retrieves an order by ID
func (u *usecase) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	ctx, span := tracer.Start(ctx, "OrderUC.GetOrderByID")
	defer span.End()

	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return order, nil
}

// ListOrders retrieves a paginated list of orders
func (u *usecase) ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, error) {
	ctx, span := tracer.Start(ctx, "OrderUC.ListOrders")
	defer span.End()

	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Fetch one extra to determine if there are more results
	orders, err := u.orderRepo.List(ctx, port.ListOrdersQuery{
		Limit:  req.Limit + 1,
		Cursor: req.Cursor,
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Determine if there are more results
	var nextCursor string
	if int32(len(orders)) > req.Limit {
		// Trim to requested page size
		orders = orders[:req.Limit]
		// Use last item's ID as next cursor
		nextCursor = orders[len(orders)-1].ID
	}

	return &ListOrdersResp{
		Orders:     orders,
		NextCursor: nextCursor,
	}, nil
}
