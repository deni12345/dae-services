package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
	"github.com/deni12345/dae-core/internal/port"
)

// GetSheet retrieves a sheet by ID
func (u *usecase) GetSheet(ctx context.Context, id string) (*domain.Sheet, error) {
	if id == "" {
		return nil, fmt.Errorf("sheet_id is required: %w", ErrInvalidArgument)
	}

	sheet, err := u.sheetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get sheet: %w", err)
	}

	return sheet, nil
}

// ListSheets retrieves a paginated list of sheets
func (u *usecase) ListSheets(ctx context.Context, req *ListSheetsReq) (*ListSheetsResp, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	sheets, err := u.sheetRepo.List(ctx, port.ListSheetsQuery{
		Limit:  int32(req.Limit) + 1,
		Cursor: req.Cursor,
	})
	if err != nil {
		return nil, fmt.Errorf("list sheets: %w", err)
	}

	var nextCursor string
	if len(sheets) > req.Limit {
		sheets = sheets[:req.Limit]
		nextCursor = sheets[len(sheets)-1].ID
	}

	return &ListSheetsResp{
		Sheets:     sheets,
		NextCursor: nextCursor,
	}, nil
}

// ListUserSheets returns all sheets accessible by a user (as host or member)
func (u *usecase) ListUserSheets(ctx context.Context, req *ListUserSheetsReq) (*ListUserSheetsResp, error) {
	// Validate request
	if req.UserID == "" {
		return nil, fmt.Errorf("user_id is required: %w", ErrInvalidArgument)
	}

	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Call repository
	resp, err := u.sheetRepo.ListForUser(ctx, port.ListSheetsForUserQuery{
		UserID:       req.UserID,
		Limit:        req.Limit,
		Cursor:       req.Cursor,
		StatusFilter: req.StatusFilter,
	})
	if err != nil {
		return nil, fmt.Errorf("list user sheets: %w", err)
	}

	return &ListUserSheetsResp{
		Sheets:     resp.Sheets,
		NextCursor: resp.NextCursor,
	}, nil
}
