package sheet

import (
	"context"

	"github.com/deni12345/dae-core/internal/domain"
	"github.com/deni12345/dae-core/internal/port"
)

// Usecase defines all sheet operations
type Usecase interface {
	// Commands
	CreateSheet(ctx context.Context, req *CreateSheetReq) (*domain.Sheet, error)
	UpdateSheet(ctx context.Context, req *UpdateSheetReq) (*domain.Sheet, error)
	JoinSheet(ctx context.Context, req *JoinSheetReq) error
	LeaveSheet(ctx context.Context, req *LeaveSheetReq) error
	CloseSheet(ctx context.Context, req *CloseSheetReq) (*domain.Sheet, error)
	ReopenSheet(ctx context.Context, req *ReopenSheetReq) (*domain.Sheet, error)

	// Queries
	GetSheet(ctx context.Context, id string) (*domain.Sheet, error)
	ListSheets(ctx context.Context, req *ListSheetsReq) (*ListSheetsResp, error)
	ListUserSheets(ctx context.Context, req *ListUserSheetsReq) (*ListUserSheetsResp, error)
	GetSheetMembers(ctx context.Context, sheetID string) ([]string, error)
}

type usecase struct {
	sheetRepo port.SheetRepo
	idemStore port.IdempotencyStore
}

// NewUsecase creates a new sheet usecase
func NewUsecase(sheetRepo port.SheetRepo, idemStore port.IdempotencyStore) Usecase {
	return &usecase{
		sheetRepo: sheetRepo,
		idemStore: idemStore,
	}
}
