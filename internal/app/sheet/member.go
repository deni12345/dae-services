package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// JoinSheet adds a user to a sheet's member list
func (u *usecase) JoinSheet(ctx context.Context, req *JoinSheetReq) error {
	// Validate request
	if req.SheetID == "" {
		return fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required: %w", ErrInvalidArgument)
	}

	// Verify sheet exists and is open
	sheet, err := u.sheetRepo.GetByID(ctx, req.SheetID)
	if err != nil {
		return fmt.Errorf("get sheet: %w", err)
	}

	if !sheet.IsOpen() {
		return fmt.Errorf("sheet is not open for joining: %w", ErrInvalidStatus)
	}

	// Add member (idempotent operation)
	if err := u.sheetRepo.AddMember(ctx, req.SheetID, req.UserID); err != nil {
		return fmt.Errorf("add member: %w", err)
	}

	return nil
}

// LeaveSheet removes a user from a sheet's member list
func (u *usecase) LeaveSheet(ctx context.Context, req *LeaveSheetReq) error {
	// Validate request
	if req.SheetID == "" {
		return fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}
	if req.UserID == "" {
		return fmt.Errorf("user_id is required: %w", ErrInvalidArgument)
	}

	// Verify user is not the host
	sheet, err := u.sheetRepo.GetByID(ctx, req.SheetID)
	if err != nil {
		return fmt.Errorf("get sheet: %w", err)
	}

	if sheet.HostUserID == req.UserID {
		return fmt.Errorf("host cannot leave sheet: %w", ErrUnauthorized)
	}

	// Remove member (idempotent operation)
	if err := u.sheetRepo.RemoveMember(ctx, req.SheetID, req.UserID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}

	return nil
}

// GetSheetMembers returns list of user IDs who are members of a sheet
func (u *usecase) GetSheetMembers(ctx context.Context, sheetID string) ([]string, error) {
	if sheetID == "" {
		return nil, fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}

	memberIDs, err := u.sheetRepo.ListMemberIDs(ctx, sheetID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	return memberIDs, nil
}

// CloseSheet closes a sheet (sets status to CLOSED)
func (u *usecase) CloseSheet(ctx context.Context, req *CloseSheetReq) (*domain.Sheet, error) {
	// Validate request
	if req.SheetID == "" {
		return nil, fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}
	if req.ActorUserID == "" {
		return nil, fmt.Errorf("actor_user_id is required: %w", ErrInvalidArgument)
	}

	// Use patch-in-transaction pattern
	updatedSheet, err := u.sheetRepo.Update(ctx, req.SheetID, func(sheet *domain.Sheet) error {
		// Business rule: only host can close
		if sheet.HostUserID != req.ActorUserID {
			return fmt.Errorf("only host can close sheet: %w", ErrUnauthorized)
		}

		// Business rule: cannot close already closed sheet
		if sheet.Status == domain.Status_CLOSED {
			return fmt.Errorf("sheet is already closed: %w", ErrInvalidStatus)
		}

		// Apply change
		sheet.Status = domain.Status_CLOSED
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("close sheet: %w", err)
	}

	return updatedSheet, nil
}

// ReopenSheet reopens a closed sheet (sets status to OPEN)
func (u *usecase) ReopenSheet(ctx context.Context, req *ReopenSheetReq) (*domain.Sheet, error) {
	// Validate request
	if req.SheetID == "" {
		return nil, fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}
	if req.ActorUserID == "" {
		return nil, fmt.Errorf("actor_user_id is required: %w", ErrInvalidArgument)
	}

	// Use patch-in-transaction pattern
	updatedSheet, err := u.sheetRepo.Update(ctx, req.SheetID, func(sheet *domain.Sheet) error {
		// Business rule: only host can reopen
		if sheet.HostUserID != req.ActorUserID {
			return fmt.Errorf("only host can reopen sheet: %w", ErrUnauthorized)
		}

		// Business rule: can only reopen closed sheets
		if sheet.Status != domain.Status_CLOSED {
			return fmt.Errorf("can only reopen closed sheets: %w", ErrInvalidStatus)
		}

		// Apply change
		sheet.Status = domain.Status_OPEN
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("reopen sheet: %w", err)
	}

	return updatedSheet, nil
}
