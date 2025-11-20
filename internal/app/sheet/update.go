package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// UpdateSheet updates an existing sheet
func (u *usecase) UpdateSheet(ctx context.Context, req *UpdateSheetReq) (*domain.Sheet, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}

	// Use callback pattern with validation
	updatedSheet, err := u.sheetRepo.Update(ctx, req.ID, func(sheet *domain.Sheet) error {
		// Apply updates
		if req.Status != nil {
			if err := validateStatusTransition(sheet.Status, *req.Status); err != nil {
				return err
			}
			sheet.Status = *req.Status
		}

		if req.DeliveryFee != nil {
			sheet.DeliveryFee = *req.DeliveryFee
		}

		if req.Discount != nil {
			if *req.Discount < 0 || *req.Discount > 100 {
				return fmt.Errorf("discount must be between 0 and 100: %w", ErrInvalidArgument)
			}
			sheet.Discount = *req.Discount
		}

		if req.Description != nil {
			sheet.Description = *req.Description
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("update sheet: %w", err)
	}

	return updatedSheet, nil
}

// validateStatusTransition validates status transitions
func validateStatusTransition(from, to domain.Status) error {
	// Define allowed transitions
	allowedTransitions := map[domain.Status][]domain.Status{
		domain.Status_OPEN:    {domain.Status_PENDING, domain.Status_CLOSED},
		domain.Status_PENDING: {domain.Status_OPEN, domain.Status_CLOSED},
		domain.Status_CLOSED:  {}, // No transitions from closed
	}

	allowed, ok := allowedTransitions[from]
	if !ok {
		return fmt.Errorf("unknown status: %d: %w", from, ErrInvalidStatus)
	}

	for _, allowedStatus := range allowed {
		if to == allowedStatus {
			return nil
		}
	}

	return fmt.Errorf("cannot transition from %s to %s: %w",
		domain.Status_name[domain.Status(from)],
		domain.Status_name[domain.Status(to)],
		ErrInvalidTransition)
}
