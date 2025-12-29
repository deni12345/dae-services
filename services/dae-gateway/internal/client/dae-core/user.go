package daecore

import (
	"context"

	pb "github.com/deni12345/dae-services/proto/gen"
)

func (c *Client) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.User.CreateUser(ctx, req)
}

func (c *Client) ListUser(ctx context.Context, req *pb.ListUsersReq) (*pb.ListUsersResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.User.ListUsers(ctx, req)
}

func (c *Client) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.User.GetUser(ctx, req)
}

func (c *Client) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.User.UpdateUser(ctx, req)
}

func (c *Client) AdminSetUserRoles(ctx context.Context, req *pb.AdminSetUserRolesReq) (*pb.AdminSetUserRolesResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.User.AdminSetUserRoles(ctx, req)
}

func (c *Client) AdminSetUserDisabled(ctx context.Context, req *pb.AdminSetUserDisabledReq) (*pb.AdminSetUserDisabledResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.User.AdminSetUserDisabled(ctx, req)
}
