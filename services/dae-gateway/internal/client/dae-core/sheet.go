package daecore

import (
	"context"

	pb "github.com/deni12345/dae-services/proto/gen"
)

func (c *Client) CreateSheet(ctx context.Context, req *pb.CreateSheetReq) (*pb.CreateSheetResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.CreateSheet(ctx, req)
}

func (c *Client) GetSheet(ctx context.Context, req *pb.GetSheetReq) (*pb.GetSheetResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.GetSheet(ctx, req)
}

func (c *Client) UpdateSheet(ctx context.Context, req *pb.UpdateSheetReq) (*pb.UpdateSheetResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.UpdateSheet(ctx, req)
}

func (c *Client) ListSheets(ctx context.Context, req *pb.ListSheetsReq) (*pb.ListSheetsResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.ListSheets(ctx, req)
}

func (c *Client) JoinSheet(ctx context.Context, req *pb.JoinSheetRequest) (*pb.JoinSheetResponse, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.JoinSheet(ctx, req)
}

func (c *Client) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.RemoveMember(ctx, req)
}

func (c *Client) ListMembers(ctx context.Context, req *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.ListMembers(ctx, req)
}

func (c *Client) AttachMenuWithPayload(ctx context.Context, req *pb.AttachMenuWithPayloadReq) (*pb.AttachMenuWithPayloadResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.AttachMenuWithPayload(ctx, req)
}

func (c *Client) GetMenu(ctx context.Context, req *pb.GetMenuReq) (*pb.GetMenuResp, error) {
	ctx, cancel := withTimeout(ctx, c.defaultTimeOut)
	defer cancel()

	return c.Sheet.GetMenu(ctx, req)
}
