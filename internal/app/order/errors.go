package order

import "fmt"

// Domain errors
var (
	ErrSheetNotFound     = fmt.Errorf("sheet not found")
	ErrSheetNotOpen      = fmt.Errorf("sheet is not open for orders")
	ErrNotFound          = fmt.Errorf("order not found")
	ErrInvalidOrderLines = fmt.Errorf("order must have at least one line")
	ErrInvalidMenuItemID = fmt.Errorf("invalid menu item id")
	ErrInvalidVariantID  = fmt.Errorf("invalid variant id")
	ErrInvalidOptionID   = fmt.Errorf("invalid option id")
)
