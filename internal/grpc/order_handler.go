package grpc

import (
	"context"

	"github.com/deni12345/dae-core/internal/app/order"
	"github.com/deni12345/dae-core/internal/grpc/converter"
	"github.com/deni12345/dae-core/internal/grpc/errors"
	corev1 "github.com/deni12345/dae-core/proto/gen"
)

type OrderHandler struct {
	corev1.UnimplementedOrdersServiceServer
	uc order.Usecase
}

func NewOrderHandler(uc order.Usecase) *OrderHandler {
	return &OrderHandler{
		uc: uc,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *corev1.CreateOrderReq) (*corev1.CreateOrderResp, error) {
	o, err := h.uc.CreateOrder(ctx, converter.CreateOrderReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}
	return &corev1.CreateOrderResp{
		Order: converter.OrderToProto(o),
	}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *corev1.GetOrderReq) (*corev1.GetOrderResp, error) {
	o, err := h.uc.GetOrderByID(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}
	return &corev1.GetOrderResp{
		Order: converter.OrderToProto(o),
	}, nil
}

func (h *OrderHandler) ListOrders(ctx context.Context, req *corev1.ListOrdersReq) (*corev1.ListOrdersResp, error) {
	resp, err := h.uc.ListOrders(ctx, converter.ListOrdersReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return converter.ListOrdersRespToProto(resp), nil
}
