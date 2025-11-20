package sheet

import "fmt"

// Domain errors
var (
	ErrNotFound          = fmt.Errorf("sheet not found")
	ErrInvalidArgument   = fmt.Errorf("invalid argument")
	ErrInvalidStatus     = fmt.Errorf("invalid sheet status")
	ErrAlreadyExists     = fmt.Errorf("sheet already exists")
	ErrUnauthorized      = fmt.Errorf("unauthorized")
	ErrInvalidTransition = fmt.Errorf("invalid status transition")

	// Menu validation errors
	ErrMenuItemNameRequired        = fmt.Errorf("menu item name required")
	ErrDuplicateMenuItemName       = fmt.Errorf("duplicate menu item name")
	ErrMenuItemInvalidPrice        = fmt.Errorf("menu item invalid price")
	ErrMenuItemInvalidCurrency     = fmt.Errorf("menu item invalid currency")
	ErrOptionGroupNameRequired     = fmt.Errorf("option group name required")
	ErrDuplicateOptionGroupName    = fmt.Errorf("duplicate option group name")
	ErrOptionGroupInvalidMaxSelect = fmt.Errorf("option group invalid max_select")
	ErrOptionNameRequired          = fmt.Errorf("option name required")
	ErrDuplicateOptionName         = fmt.Errorf("duplicate option name")
	ErrOptionInvalidPrice          = fmt.Errorf("option invalid price")
)
