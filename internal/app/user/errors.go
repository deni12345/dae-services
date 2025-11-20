package user

import "fmt"

// Domain errors
var (
	ErrInvalidArgument = fmt.Errorf("invalid argument")
	ErrNotFound        = fmt.Errorf("user not found")
	ErrInvalidRole     = fmt.Errorf("invalid role")
)
