package sheet

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Repository errors - these should be mapped to domain errors in use cases
var (
	ErrSheetNotFound      = errors.New("sheet not found")
	ErrSheetAlreadyExists = errors.New("sheet already exists")
	ErrConcurrentUpdate   = errors.New("concurrent update detected")
	ErrInvalidCursor      = errors.New("invalid cursor")
)

// mapFirestoreError maps Firestore gRPC status codes to repository errors
func mapFirestoreError(err error, operation string) error {
	if err == nil {
		return nil
	}

	code := status.Code(err)
	switch code {
	case codes.NotFound:
		return fmt.Errorf("%s: %w", operation, ErrSheetNotFound)
	case codes.AlreadyExists:
		return fmt.Errorf("%s: %w", operation, ErrSheetAlreadyExists)
	case codes.FailedPrecondition:
		return fmt.Errorf("%s: %w", operation, ErrConcurrentUpdate)
	case codes.InvalidArgument:
		return fmt.Errorf("%s: invalid argument: %w", operation, err)
	case codes.DeadlineExceeded:
		return fmt.Errorf("%s: deadline exceeded: %w", operation, err)
	case codes.ResourceExhausted:
		return fmt.Errorf("%s: resource exhausted: %w", operation, err)
	default:
		return fmt.Errorf("%s: %w", operation, err)
	}
}
