package converter

import (
	"github.com/deni12345/dae-core/common/utils"
	"github.com/deni12345/dae-core/internal/app/user"
	"github.com/deni12345/dae-core/internal/domain"
	corev1 "github.com/deni12345/dae-core/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var protoToDomainRolesMap = map[corev1.UserRole]domain.Role{
	corev1.UserRole_USER_ROLE_ADMIN: domain.RoleAdmin,
	corev1.UserRole_USER_ROLE_USER:  domain.RoleUser,
}

var domainToProtoRolesMap = map[domain.Role]corev1.UserRole{
	domain.RoleAdmin: corev1.UserRole_USER_ROLE_ADMIN,
	domain.RoleUser:  corev1.UserRole_USER_ROLE_USER,
}

// Proto to DTO conversions

func UpdateUserReqFromProto(req *corev1.UpdateUserReq) *user.UpdateUserReq {
	return &user.UpdateUserReq{
		ID:         req.Id,
		UserName:   req.DisplayName,
		AvatarURL:  req.AvatarUrl,
		IsDisabled: req.IsDisabled,
	}
}

func AdminSetUserRolesReqFromProto(req *corev1.AdminSetUserRolesReq) *user.AdminSetUserRolesReq {
	roles := make([]domain.Role, len(req.Roles))
	for i, r := range req.Roles {
		roles[i] = domain.Role(r)
	}
	return &user.AdminSetUserRolesReq{
		UserID: req.UserId,
		Roles:  roles,
	}
}

func AdminSetUserDisabledReqFromProto(req *corev1.AdminSetUserDisabledReq) *user.AdminSetUserDisabledReq {
	return &user.AdminSetUserDisabledReq{
		UserID:     req.UserId,
		IsDisabled: req.IsDisabled,
	}
}

func ListUsersReqFromProto(req *corev1.ListUsersReq) *user.ListUsersReq {
	dto := &user.ListUsersReq{
		PageSize: req.GetPageSize(),
	}

	// Parse cursor if provided (ID only)
	if cursor := req.GetCursor(); cursor != nil && cursor.GetId() != "" {
		dto.Cursor = cursor.GetId()
	}

	// Parse filter if provided
	if filter := req.GetFilter(); filter != nil {
		dto.IncludeDisabled = filter.GetIncludeDisabled()
		dto.Query = filter.GetQuery()
		dto.EmailExact = filter.GetEmailExact()
	}

	return dto
}

// Domain to Proto conversions

func UserToProto(u *domain.User) *corev1.User {
	if u == nil {
		return nil
	}

	roles := rolesToProto(u.Roles)

	return &corev1.User{
		Id:          u.ID,
		Email:       u.Email,
		DisplayName: u.UserName,
		AvatarUrl:   u.AvatarURL,
		Roles:       roles,
		IsDisabled:  u.IsDisabled,
		CreatedAt:   timestamppb.New(u.CreatedAt),
		UpdatedAt:   timestamppb.New(u.UpdatedAt),
	}
}

func rolesToProto(roles []domain.Role) []corev1.UserRole {
	roles = utils.ToSet(roles)
	if len(roles) == 0 {
		return []corev1.UserRole{}
	}

	res := make([]corev1.UserRole, 0, len(roles))
	for _, r := range roles {
		if ur, ok := domainToProtoRolesMap[r]; ok {
			res = append(res, ur)
		}
	}
	return res
}

func UsersToProto(users []*domain.User) []*corev1.User {
	result := make([]*corev1.User, len(users))
	for i, u := range users {
		result[i] = UserToProto(u)
	}
	return result
}

func DomainRolesFromProto(roles []corev1.UserRole) []domain.Role {
	roles = utils.ToSet(roles)
	if len(roles) == 0 {
		return []domain.Role{}
	}

	res := make([]domain.Role, 0, len(roles))
	for _, r := range roles {
		if dr, ok := protoToDomainRolesMap[r]; ok {
			res = append(res, dr)
		}
	}
	return res
}

// ListUsersRespToProto converts DTO response to proto
func ListUsersRespToProto(resp *user.ListUsersResp) *corev1.ListUsersResp {
	if resp == nil {
		return &corev1.ListUsersResp{}
	}

	protoResp := &corev1.ListUsersResp{
		Users: UsersToProto(resp.Users),
	}

	// Add next cursor if available (ID only)
	if resp.NextCursor != "" {
		protoResp.NextCursor = &corev1.Cursor{
			Id: resp.NextCursor,
		}
	}

	return protoResp
}
