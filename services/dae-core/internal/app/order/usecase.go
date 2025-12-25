package order

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"go.opentelemetry.io/otel"
)

// Usecase defines all order operations
type Usecase interface {
	// Commands
	CreateOrder(ctx context.Context, req *CreateOrderReq) (*domain.Order, error)
	UpdateOrder(ctx context.Context, req *UpdateOrderReq) (*domain.Order, error)

	// Queries
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	ListOrders(ctx context.Context, req *ListOrdersReq) (*ListOrdersResp, error)
}

type usecase struct {
	orderRepo port.OrdersRepo
	sheetRepo port.SheetRepo
	idemStore port.IdempotencyStore
}

// NewUsecase creates a new order usecase
func NewUsecase(orderRepo port.OrdersRepo, sheetRepo port.SheetRepo, idemStore port.IdempotencyStore) Usecase {
	return &usecase{
		orderRepo: orderRepo,
		sheetRepo: sheetRepo,
		idemStore: idemStore,
	}
}

var tracer = otel.Tracer("usecase/order")
