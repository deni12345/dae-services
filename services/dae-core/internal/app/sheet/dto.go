package sheet

import "github.com/deni12345/dae-services/services/dae-core/internal/domain"

// Request DTOs for nested structures

type MenuOptionReq struct {
	ID     string
	Name   string
	Price  int64
	Active bool
}

type MenuOptionGroupReq struct {
	ID          string
	Name        string
	Required    bool
	MultiSelect bool
	MinSelect   int32
	MaxSelect   int32
	Options     []MenuOptionReq
}

type MenuItemReq struct {
	ID           string
	Name         string
	Description  string
	Active       bool
	Price        int64
	Currency     string
	OptionGroups []MenuOptionGroupReq
}

// Command DTOs

type CreateSheetReq struct {
	IdempotencyKey string
	Name           string
	HostUserID     string
	DeliveryFee    *domain.Money
	Discount       int32
	Description    string
	MemberIDs      []string
	MenuItems      []MenuItemReq // Clean request, not domain entities
}

type UpdateSheetReq struct {
	ID          string
	Name        *string
	Description *string
	Status      *domain.Status
	DeliveryFee *domain.Money
	Discount    *int32
}

// Query DTOs

type ListSheetsReq struct {
	Limit      int
	Cursor     string
	HostUserID *string
}

type ListSheetsResp struct {
	Sheets     []*domain.Sheet
	NextCursor string
}

type ListUserSheetsReq struct {
	UserID       string
	Limit        int32
	Cursor       string
	StatusFilter *domain.Status
}

type ListUserSheetsResp struct {
	Sheets     []*domain.Sheet
	NextCursor string
}

// Command DTOs

type JoinSheetReq struct {
	SheetID string
	UserID  string
}

type LeaveSheetReq struct {
	SheetID string
	UserID  string
}

type CloseSheetReq struct {
	SheetID     string
	ActorUserID string
}

type ReopenSheetReq struct {
	SheetID     string
	ActorUserID string
}
