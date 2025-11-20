package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deni12345/dae-core/internal/domain"
	"github.com/deni12345/dae-core/internal/grpc/interceptor"
	"github.com/google/uuid"
)

const (
	idempotencyTTL = 24 * time.Hour // Cache successful creations for 24h
)

// CreateOrder creates a new order with pricing validation and idempotency protection
func (u *usecase) CreateOrder(ctx context.Context, req *CreateOrderReq) (*domain.Order, error) {
	// Validate early
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Use idempotency key from request, fallback to generated key
	idemKey := ctx.Value(interceptor.MDKeyIdem).(string)
	if idemKey == "" {
		idemKey = fmt.Sprintf("order:create:%s:%s", req.SheetID, req.UserID)
	}

	result, err := u.idemStore.Do(ctx, idemKey, idempotencyTTL, func(ctx context.Context) ([]byte, error) {
		order, err := u.createOrderInternal(ctx, req)
		if err != nil {
			return nil, err
		}
		return json.Marshal(order)
	})

	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	var order domain.Order
	if err := json.Unmarshal(result, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order: %w", err)
	}

	return &order, nil
}

// validateCreateRequest validates the create order request
func validateCreateRequest(req *CreateOrderReq) error {
	if req.SheetID == "" {
		return fmt.Errorf("sheet_id is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if len(req.Lines) == 0 {
		return ErrInvalidOrderLines
	}
	return nil
}

// createOrderInternal is the core order creation logic without idempotency
func (u *usecase) createOrderInternal(ctx context.Context, req *CreateOrderReq) (*domain.Order, error) {
	// Fetch and validate sheet
	sheet, err := u.sheetRepo.GetByID(ctx, req.SheetID)
	if err != nil {
		return nil, fmt.Errorf("get sheet: %w", err)
	}

	if !sheet.IsOpen() {
		return nil, ErrSheetNotOpen
	}

	// Build order lines with pricing
	orderLines, err := u.buildOrderLines(ctx, req.SheetID, req.Lines)
	if err != nil {
		return nil, err
	}

	// Create order entity
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

	// Persist
	createdOrder, err := u.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("persist order: %w", err)
	}

	return createdOrder, nil
}

// buildOrderLines builds all order lines with pricing validation
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
