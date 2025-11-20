package order

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
	"github.com/deni12345/dae-core/internal/port"
)

// GetOrderByID retrieves an order by ID
func (u *usecase) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order_id is required")
	}

	order, err := u.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	return order, nil
}

// ListOrders retrieves a paginated list of orders
func (u *usecase) ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, error) {
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
		return nil, fmt.Errorf("list orders: %w", err)
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
