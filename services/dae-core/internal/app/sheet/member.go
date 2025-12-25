package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/libs/apperror"
)

// JoinSheet adds a user to a sheet's member list
func (u *usecase) JoinSheet(ctx context.Context, req *JoinSheetReq) error {
	ctx, span := tracer.Start(ctx, "SheetUC.JoinSheet")
	defer span.End()

	// Validate request
	if req.SheetID == "" {
		err := apperror.InvalidInput("sheet_id is required")
		span.RecordError(err)
		return err
	}
	if req.UserID == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return err
	}

	// Verify sheet exists and is open
	sheet, err := u.sheetRepo.GetByID(ctx, req.SheetID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if !sheet.IsOpen() {
		err := apperror.InvalidInput(fmt.Sprintf("sheet %s is not open for joining", req.SheetID))
		span.RecordError(err)
		return err
	}

	// Add member (idempotent operation)
	if err := u.sheetRepo.AddMember(ctx, req.SheetID, req.UserID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

// LeaveSheet removes a user from a sheet's member list
func (u *usecase) LeaveSheet(ctx context.Context, req *LeaveSheetReq) error {
	ctx, span := tracer.Start(ctx, "SheetUC.LeaveSheet")
	defer span.End()

	// Validate request
	if req.SheetID == "" {
		err := apperror.InvalidInput("sheet_id is required")
		span.RecordError(err)
		return err
	}
	if req.UserID == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return err
	}

	// Verify user is not the host
	sheet, err := u.sheetRepo.GetByID(ctx, req.SheetID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	if sheet.HostUserID == req.UserID {
		err := apperror.Forbidden("host cannot leave sheet")
		span.RecordError(err)
		return err
	}

	// Remove member (idempotent operation)
	if err := u.sheetRepo.RemoveMember(ctx, req.SheetID, req.UserID); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

// GetSheetMembers returns list of user IDs who are members of a sheet
func (u *usecase) GetSheetMembers(ctx context.Context, sheetID string) ([]string, error) {

	memberIDs, err := u.sheetRepo.ListMemberIDs(ctx, sheetID)
	if err != nil {
		return nil, err
	}

	return memberIDs, nil
}

// CloseSheet closes a sheet (sets status to CLOSED)
func (u *usecase) CloseSheet(ctx context.Context, req *CloseSheetReq) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetUC.CloseSheet")
	defer span.End()

	// Validate request
	if req.SheetID == "" {
		err := apperror.InvalidInput("sheet_id is required")
		span.RecordError(err)
		return nil, err
	}
	if req.ActorUserID == "" {
		err := apperror.InvalidInput("actor_user_id is required")
		span.RecordError(err)
		return nil, err
	}

	// Use patch-in-transaction pattern
	updatedSheet, err := u.sheetRepo.Update(ctx, req.SheetID, func(sheet *domain.Sheet) error {
		// Business rule: only host can close
		if sheet.HostUserID != req.ActorUserID {
			return apperror.Forbidden("only host can close sheet")
		}

		// Business rule: cannot close already closed sheet
		if sheet.Status == domain.Status_CLOSED {
			return apperror.InvalidInput(fmt.Sprintf("sheet %s is already closed", req.SheetID))
		}

		// Apply change
		sheet.Status = domain.Status_CLOSED
		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return updatedSheet, nil
}

// ReopenSheet reopens a closed sheet (sets status to OPEN)
func (u *usecase) ReopenSheet(ctx context.Context, req *ReopenSheetReq) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetUC.ReopenSheet")
	defer span.End()

	// Validate request
	if req.SheetID == "" {
		err := apperror.InvalidInput("sheet_id is required")
		span.RecordError(err)
		return nil, err
	}
	if req.ActorUserID == "" {
		err := apperror.InvalidInput("actor_user_id is required")
		span.RecordError(err)
		return nil, err
	}

	// Use patch-in-transaction pattern
	updatedSheet, err := u.sheetRepo.Update(ctx, req.SheetID, func(sheet *domain.Sheet) error {
		// Business rule: only host can reopen
		if sheet.HostUserID != req.ActorUserID {
			return apperror.Forbidden("only host can reopen sheet")
		}

		// Business rule: can only reopen closed sheets
		if sheet.Status != domain.Status_CLOSED {
			return apperror.InvalidInput(fmt.Sprintf("can only reopen closed sheets, current status: %s", domain.Status_name[domain.Status(sheet.Status)]))
		}

		// Apply change
		sheet.Status = domain.Status_OPEN
		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return updatedSheet, nil
}
