package sheet

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"github.com/deni12345/dae-services/libs/apperror"
)

// GetSheet retrieves a sheet by ID
func (u *usecase) GetSheet(ctx context.Context, id string) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetUC.GetSheet")
	defer span.End()

	if id == "" {
		err := apperror.InvalidInput("sheet_id is required")
		span.RecordError(err)
		return nil, err
	}

	sheet, err := u.sheetRepo.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return sheet, nil
}

// ListSheets retrieves a paginated list of sheets
func (u *usecase) ListSheets(ctx context.Context, req *ListSheetsReq) (*ListSheetsResp, error) {
	ctx, span := tracer.Start(ctx, "SheetUC.ListSheets")
	defer span.End()

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
		span.RecordError(err)
		return nil, err
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
	ctx, span := tracer.Start(ctx, "SheetUC.ListUserSheets")
	defer span.End()

	// Validate request
	if req.UserID == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return nil, err
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
		span.RecordError(err)
		return nil, err
	}

	return &ListUserSheetsResp{
		Sheets:     resp.Sheets,
		NextCursor: resp.NextCursor,
	}, nil
}
