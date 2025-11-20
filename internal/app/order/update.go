package order

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// UpdateOrder updates an existing order with re-pricing
func (u *usecase) UpdateOrder(ctx context.Context, req *UpdateOrderReq) (*domain.Order, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("order_id is required")
	}
	if len(req.Lines) == 0 {
		return nil, ErrInvalidOrderLines
	}

	// Use callback pattern to fetch, validate, and update
	updatedOrder, err := u.orderRepo.Update(ctx, req.ID, func(order *domain.Order) error {
		// Verify sheet is still open for updates
		sheet, err := u.sheetRepo.GetByID(ctx, order.SheetID)
		if err != nil {
			return fmt.Errorf("get sheet: %w", err)
		}

		if !sheet.IsOpen() {
			return ErrSheetNotOpen
		}

		// Re-price all lines with correct sheetID
		newLines := make([]domain.OrderLine, 0, len(req.Lines))
		for _, lineReq := range req.Lines {
			line, err := u.buildOrderLine(ctx, order.SheetID, lineReq)
			if err != nil {
				return fmt.Errorf("build order line: %w", err)
			}
			newLines = append(newLines, line)
		}

		// Apply changes to the current order
		order.Lines = newLines
		order.Note = req.Note
		domain.CalculateOrderTotals(order)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("update order: %w", err)
	}

	return updatedOrder, nil
}
