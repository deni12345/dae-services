package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/grpc/interceptor"
	"github.com/deni12345/dae-services/libs/apperror"
	"github.com/google/uuid"
)

const (
	idempotencyTTL = 24 * time.Hour
)

// CreateOrder creates an order.
func (u *usecase) CreateOrder(ctx context.Context, req *CreateOrderReq) (*domain.Order, error) {
	ctx, span := tracer.Start(ctx, "OrderUC.CreateOrder")
	defer span.End()

	idemKey := interceptor.GetOrCreateIdempotencyKeyWithHash(ctx, req.SheetID, req.UserID)

	result, err := u.idemStore.Do(ctx, idemKey, idempotencyTTL, func(ctx context.Context) ([]byte, error) {
		order, err := u.createOrderInternal(ctx, req)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		return json.Marshal(order)
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal(result, &order); err != nil {
		span.RecordError(err)
		return nil, apperror.Internal(fmt.Sprintf("unmarshal order: %v", err))
	}

	return &order, nil
}

func validateCreateRequest(req *CreateOrderReq) error {
	if req.SheetID == "" {
		return apperror.InvalidInput("sheet_id is required")
	}
	if req.UserID == "" {
		return apperror.InvalidInput("user_id is required")
	}
	if len(req.Lines) == 0 {
		return ErrInvalidOrderLines
	}
	return nil
}

func (u *usecase) createOrderInternal(ctx context.Context, req *CreateOrderReq) (*domain.Order, error) {
	sheet, err := u.sheetRepo.GetByID(ctx, req.SheetID)
	if err != nil {
		return nil, err
	}
	if !sheet.IsOpen() {
		return nil, ErrSheetNotOpen
	}

	orderLines, err := u.buildOrderLines(ctx, req.SheetID, req.Lines)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	order := &domain.Order{
		ID:        uuid.New().String(),
		SheetID:   req.SheetID,
		UserID:    req.UserID,
		Lines:     orderLines,
		Note:      req.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	domain.CalculateOrderTotals(order)

	createdOrder, err := u.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (u *usecase) buildOrderLines(ctx context.Context, sheetID string, lineReqs []OrderLineReq) ([]domain.OrderLine, error) {
	orderLines := make([]domain.OrderLine, 0, len(lineReqs))
	for i, lineReq := range lineReqs {
		line, err := u.buildOrderLine(ctx, sheetID, lineReq)
		if err != nil {
			return nil, fmt.Errorf("build line %d: %w", i, err)
		}
		orderLines = append(orderLines, line)
	}
	return orderLines, nil
}
