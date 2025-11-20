package user

import (
	"github.com/deni12345/dae-core/internal/domain"
)

// Command DTOs - for write operations

type UpdateUserReq struct {
	ID             string
	IdempotencyKey string // Required for write operations
	UserName       *string
	AvatarURL      *string
	IsDisabled     *bool
}

type AdminSetUserRolesReq struct {
	UserID string
	Roles  []domain.Role
}

type AdminSetUserDisabledReq struct {
	UserID     string
	IsDisabled bool
}

// Query DTOs - for read operations

type ListUsersReq struct {
	PageSize        int32
	Cursor          string
	IncludeDisabled bool
	Query           string
	EmailExact      string
}

type ListUsersResp struct {
	Users      []*domain.User
	NextCursor string
}
