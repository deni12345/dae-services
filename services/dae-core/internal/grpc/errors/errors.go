package errors

import (
	"github.com/deni12345/dae-services/libs/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// appErrorToGRPC maps AppError codes to gRPC status codes
func ToGRPCStatus(err error) error {
	var code codes.Code

	switch apperror.GetCode(err) {
	case apperror.CodeNotFound:
		code = codes.NotFound
	case apperror.CodeAlreadyExists:
		code = codes.AlreadyExists
	case apperror.CodeInvalidInput:
		code = codes.InvalidArgument
	case apperror.CodeUnauthorized:
		code = codes.Unauthenticated
	case apperror.CodeForbidden:
		code = codes.PermissionDenied
	case apperror.CodeConflict:
		code = codes.Aborted
	case apperror.CodeInternal:
		code = codes.Internal
	default:
		code = codes.Internal
	}

	return status.Error(code, err.Error())
}
