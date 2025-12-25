package grpc

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/app/user"
	"github.com/deni12345/dae-services/services/dae-core/internal/grpc/converter"
	"github.com/deni12345/dae-services/services/dae-core/internal/grpc/errors"
	corev1 "github.com/deni12345/dae-services/proto/gen"
)

type UserHandler struct {
	corev1.UnimplementedUsersServiceServer
	uc user.Usecase
}

func NewUserHandler(uc user.Usecase) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *corev1.CreateUserReq) (*corev1.CreateUserResp, error) {
	u, err := h.uc.CreateUser(ctx, converter.CreateUserReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.CreateUserResp{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *corev1.GetUserReq) (*corev1.GetUserResp, error) {
	u, err := h.uc.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.GetUserResp{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *corev1.UpdateUserReq) (*corev1.UpdateUserResp, error) {
	u, err := h.uc.UpdateUser(ctx, converter.UpdateUserReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.UpdateUserResp{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *corev1.ListUsersReq) (*corev1.ListUsersResp, error) {
	resp, err := h.uc.ListUsers(ctx, converter.ListUsersReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return converter.ListUsersRespToProto(resp), nil
}

func (h *UserHandler) AdminSetUserRoles(ctx context.Context, req *corev1.AdminSetUserRolesReq) (*corev1.AdminSetUserRolesResp, error) {
	u, err := h.uc.AdminSetUserRoles(ctx, converter.AdminSetUserRolesReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.AdminSetUserRolesResp{
		User: converter.UserToProto(u),
	}, nil
}

func (h *UserHandler) AdminSetUserDisabled(ctx context.Context, req *corev1.AdminSetUserDisabledReq) (*corev1.AdminSetUserDisabledResp, error) {
	u, err := h.uc.AdminSetUserDisabled(ctx, converter.AdminSetUserDisabledReqFromProto(req))
	if err != nil {
		return nil, errors.ToGRPCStatus(err)
	}

	return &corev1.AdminSetUserDisabledResp{
		User: converter.UserToProto(u),
	}, nil
}
