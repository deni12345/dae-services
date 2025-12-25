package order

import (
	"time"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
)

// Command DTOs - for write operations

type OrderLineOptionReq struct {
	GroupID  string
	OptionID string
	Quantity int
}

type OrderLineReq struct {
	MenuItemID string
	Options    []OrderLineOptionReq
	Quantity   int
	Note       string
}

type CreateOrderReq struct {
	SheetID string
	Lines   []OrderLineReq
	Note    string
	UserID  string
}

type UpdateOrderReq struct {
	ID             string
	IdempotencyKey string // Required for write operations
	Lines          []OrderLineReq
	Note           string
}

// Query DTOs - for read operations

type ListFilter struct {
	UserID *string
	Since  *time.Time
}

type ListOrdersReq struct {
	Limit   int32  `json:"limit,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
	SheetID string `json:"sheet_id,omitempty"`
	Filter  ListFilter
}

type ListOrdersResp struct {
	Orders     []*domain.Order
	NextCursor string
}
