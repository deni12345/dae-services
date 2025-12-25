package user

import "github.com/deni12345/dae-services/libs/apperror"

// Domain errors using apperror for better gRPC mapping
var (
	ErrInvalidArgument = apperror.InvalidInput("invalid argument")
	ErrNotFound        = apperror.NotFound("user not found")
	ErrInvalidRole     = apperror.InvalidInput("invalid role")

	// Create user errors
	ErrEmailRequired         = apperror.InvalidInput("email is required")
	ErrNameRequired          = apperror.InvalidInput("name is required")
	ErrProviderRequired      = apperror.InvalidInput("provider is required")
	ErrSubjectRequired       = apperror.InvalidInput("subject is required")
	ErrPasswordRequired      = apperror.InvalidInput("password is required for local provider")
	ErrEmailAlreadyExists    = apperror.AlreadyExists("email already exists")
	ErrIdentityAlreadyLinked = apperror.AlreadyExists("identity already linked to another user")
	ErrInvalidInput          = apperror.InvalidInput("invalid input")
)
