package daecore

import (
	"context"

	pb "github.com/deni12345/dae-services/proto/gen"
)

func (c *Client) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Order.CreateOrder(ctx, req)
}

func (c *Client) UpdateOrder(ctx context.Context, req *pb.UpdateOrderReq) (*pb.UpdateOrderResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Order.UpdateOrder(ctx, req)
}

func (c *Client) GetOrder(ctx context.Context, req *pb.GetOrderReq) (*pb.GetOrderResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Order.GetOrder(ctx, req)
}

func (c *Client) ListOrders(ctx context.Context, req *pb.ListOrdersReq) (*pb.ListOrdersResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Order.ListOrders(ctx, req)
}
