package grpc

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/app/health"
	"github.com/deni12345/dae-services/services/dae-core/internal/grpc/errors"
	corev1 "github.com/deni12345/dae-services/proto/gen"
	"github.com/deni12345/dae-services/libs/apperror"
)

// HealthHandler is a thin adapter that delegates health checks to the usecase.
type HealthHandler struct {
	corev1.UnimplementedHealthServiceServer
	h health.Usecase
}

func NewHealthHandler(hc health.Usecase) *HealthHandler {
	return &HealthHandler{h: hc}
}

func (h *HealthHandler) CheckHealth(ctx context.Context, req *corev1.HealthCheckReq) (*corev1.HealthCheckResp, error) {
	if err := h.h.Check(ctx); err != nil {
		return &corev1.HealthCheckResp{Message: "NOT_SERVING: " + err.Error()}, errors.ToGRPCStatus(apperror.Internal(err.Error()))
	}
	return &corev1.HealthCheckResp{Message: "SERVING"}, nil
}
