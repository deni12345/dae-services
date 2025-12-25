package port

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
)

type ListOrdersQuery struct {
	Limit  int32
	Cursor string
}

// OrdersRepo defines the interface for persisting and retrieving orders
type OrdersRepo interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, error)
	Update(ctx context.Context, id string, fn func(o *domain.Order) error) (*domain.Order, error)
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	List(ctx context.Context, query ListOrdersQuery) ([]*domain.Order, error)
}
