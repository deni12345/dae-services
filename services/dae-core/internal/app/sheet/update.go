package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/libs/apperror"
)

// UpdateSheet updates an existing sheet
func (u *usecase) UpdateSheet(ctx context.Context, req *UpdateSheetReq) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetUC.UpdateSheet")
	defer span.End()

	if req.ID == "" {
		err := apperror.InvalidInput("sheet_id is required")
		span.RecordError(err)
		return nil, err
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
				return apperror.InvalidInput(fmt.Sprintf("discount must be between 0 and 100, got %d", *req.Discount))
			}
			sheet.Discount = *req.Discount
		}

		if req.Description != nil {
			sheet.Description = *req.Description
		}

		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
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
		return apperror.InvalidInput(fmt.Sprintf("unknown status: %d", from))
	}

	for _, allowedStatus := range allowed {
		if to == allowedStatus {
			return nil
		}
	}

	return apperror.InvalidInput(fmt.Sprintf("cannot transition from %s to %s",
		domain.Status_name[domain.Status(from)],
		domain.Status_name[domain.Status(to)]))
}
