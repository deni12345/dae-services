package interceptor

import (
	"context"
	"errors"
	"log/slog"

	"github.com/deni12345/dae-services/services/dae-core/internal/infra/observability"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Validator interface{ ValidateAll() error }
type ProtoValidator interface{ Validate() error }

type multiError interface{ AllErrors() []error }
type fieldError interface {
	Field() string
	Reason() string
	Cause() error
	Key() bool
	ErrorName() string
}

func ValidateRequestInterceptor(metrics *observability.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if req == nil {
			return nil, status.Error(codes.InvalidArgument, "request must not be nil")
		}

		var err error

		// try ValidateAll() or generated Validate()
		if v, ok := req.(Validator); ok {
			err = v.ValidateAll()
		} else if p, ok := req.(ProtoValidator); ok {
			err = p.Validate()
		}

		if err != nil {
			if metrics != nil {
				metrics.ValidationFailures.Add(ctx, 1)
			}

			st := status.New(codes.InvalidArgument, err.Error())
			br := &errdetails.BadRequest{}
			var mErr multiError
			if errors.As(err, &mErr) {
				inner := mErr.AllErrors()
				br.FieldViolations = make([]*errdetails.BadRequest_FieldViolation, 0, len(inner))
				for _, ie := range inner {
					var fe fieldError
					if errors.As(ie, &fe) {
						br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{Field: fe.Field(), Description: fe.Reason()})
					} else {
						br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{Field: "", Description: ie.Error()})
					}
				}
			} else {
				// Single validation error (generated validators expose Field() and Reason())
				var fe fieldError
				if errors.As(err, &fe) {
					br.FieldViolations = []*errdetails.BadRequest_FieldViolation{{Field: fe.Field(), Description: fe.Reason()}}
				} else {
					br.FieldViolations = []*errdetails.BadRequest_FieldViolation{{Field: "", Description: err.Error()}}
				}
			}
			stWithDetails, e := st.WithDetails(br)
			if e != nil {
				slog.Warn("failed to attach validation details to status", "error", e)
				return nil, st.Err()
			}
			return nil, stWithDetails.Err()
		}
		return handler(ctx, req)
	}
}
