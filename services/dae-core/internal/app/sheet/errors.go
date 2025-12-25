package sheet

import "github.com/deni12345/dae-services/libs/apperror"

// Domain errors using apperror for better gRPC mapping
var (
	ErrNotFound          = apperror.NotFound("sheet not found")
	ErrInvalidArgument   = apperror.InvalidInput("invalid argument")
	ErrInvalidStatus     = apperror.InvalidInput("invalid sheet status")
	ErrAlreadyExists     = apperror.AlreadyExists("sheet already exists")
	ErrUnauthorized      = apperror.Unauthorized("unauthorized")
	ErrInvalidTransition = apperror.InvalidInput("invalid status transition")

	// Menu validation errors
	ErrMenuItemNameRequired        = apperror.InvalidInput("menu item name required")
	ErrDuplicateMenuItemName       = apperror.AlreadyExists("duplicate menu item name")
	ErrMenuItemInvalidPrice        = apperror.InvalidInput("menu item invalid price")
	ErrMenuItemInvalidCurrency     = apperror.InvalidInput("menu item invalid currency")
	ErrOptionGroupNameRequired     = apperror.InvalidInput("option group name required")
	ErrDuplicateOptionGroupName    = apperror.AlreadyExists("duplicate option group name")
	ErrOptionGroupInvalidMaxSelect = apperror.InvalidInput("option group invalid max_select")
	ErrOptionNameRequired          = apperror.InvalidInput("option name required")
	ErrDuplicateOptionName         = apperror.AlreadyExists("duplicate option name")
	ErrOptionInvalidPrice          = apperror.InvalidInput("option invalid price")
)
