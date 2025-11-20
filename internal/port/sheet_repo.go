package port

import (
	"context"

	"github.com/deni12345/dae-core/internal/domain"
)

type ListSheetsQuery struct {
	Limit  int32
	Cursor string
}

type ListSheetsForUserQuery struct {
	UserID       string
	Limit        int32
	Cursor       string // last joined_at timestamp for pagination
	StatusFilter *domain.Status
}

type ListSheetsForUserResp struct {
	Sheets     []*domain.Sheet
	NextCursor string
}

// SheetRepo defines the interface for persisting and retrieving sheets
type SheetRepo interface {
	GetByID(ctx context.Context, id string) (*domain.Sheet, error)
	Create(ctx context.Context, sheet *domain.Sheet) (*domain.Sheet, error)
	Update(ctx context.Context, id string, fn func(sheet *domain.Sheet) error) (*domain.Sheet, error)
	List(ctx context.Context, query ListSheetsQuery) ([]*domain.Sheet, error)
	ListForUser(ctx context.Context, query ListSheetsForUserQuery) (*ListSheetsForUserResp, error)

	AddMember(ctx context.Context, sheetID string, userID string) error
	RemoveMember(ctx context.Context, sheetID string, userID string) error
	ListMemberIDs(ctx context.Context, sheetID string) ([]string, error)

	// Menu Items by Sheet ID
	GetMenuItems(ctx context.Context, sheetID string) ([]*domain.MenuItem, error)
	// Menu particular Item by ID
	GetMenuItemByID(ctx context.Context, sheetID string, id string) (*domain.MenuItem, error)

	AttachMenuItems(ctx context.Context, sheetID string, menuItems []*domain.MenuItem) error
}
