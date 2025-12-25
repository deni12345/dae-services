package order

import "github.com/deni12345/dae-services/libs/apperror"

// Domain errors using apperror for better gRPC mapping
var (
	ErrSheetNotFound     = apperror.NotFound("sheet not found")
	ErrSheetNotOpen      = apperror.InvalidInput("sheet is not open for orders")
	ErrNotFound          = apperror.NotFound("order not found")
	ErrInvalidOrderLines = apperror.InvalidInput("order must have at least one line")
	ErrInvalidMenuItemID = apperror.InvalidInput("invalid menu item id")
	ErrInvalidVariantID  = apperror.InvalidInput("invalid variant id")
	ErrInvalidOptionID   = apperror.InvalidInput("invalid option id")
)
