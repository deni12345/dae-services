package daecore

import (
	"context"

	pb "github.com/deni12345/dae-services/proto/gen"
)

func (c *Client) CheckHealth(ctx context.Context, req *pb.HealthCheckReq) (*pb.HealthCheckResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Health.CheckHealth(ctx, req)
}
