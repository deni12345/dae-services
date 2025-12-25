package apperror

import "errors"

type Code string

const (
	CodeNotFound      Code = "NOT_FOUND"
	CodeAlreadyExists Code = "ALREADY_EXISTS"
	CodeInvalidInput  Code = "INVALID_INPUT"
	CodeUnauthorized  Code = "UNAUTHORIZED"
	CodeForbidden     Code = "FORBIDDEN"
	CodeInternal      Code = "INTERNAL_ERROR"
	CodeConflict      Code = "CONFLICT"
)

type AppError struct {
	message string
	code    Code
	cause   error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) Code() Code {
	return e.code
}

// New creates a new AppError with the given message and code
func New(message string, code Code) *AppError {
	return &AppError{
		message: message,
		code:    code,
	}
}

// Wrap wraps an error with an AppError, preserving the error chain
func Wrap(err error, message string, code Code) *AppError {
	return &AppError{
		message: message,
		code:    code,
		cause:   err,
	}
}

// Helper constructors for common error types
func NotFound(message string) *AppError {
	return New(message, CodeNotFound)
}

func AlreadyExists(message string) *AppError {
	return New(message, CodeAlreadyExists)
}

func InvalidInput(message string) *AppError {
	return New(message, CodeInvalidInput)
}

func Unauthorized(message string) *AppError {
	return New(message, CodeUnauthorized)
}

func Forbidden(message string) *AppError {
	return New(message, CodeForbidden)
}

func Internal(message string) *AppError {
	return New(message, CodeInternal)
}

func Conflict(message string) *AppError {
	return New(message, CodeConflict)
}

// GetCode extracts the Code from an error if it's an AppError, otherwise returns CodeInternal
func GetCode(err error) Code {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code()
	}
	return CodeInternal
}
