package port

import (
	"context"

	"github.com/deni12345/dae-core/internal/domain"
)

type ListUsersReq struct {
	Limit           int
	Cursor          string
	ExactEmail      string
	Query           string
	IncludeDisabled bool
}

type ListUsersResp struct {
	Users      []*domain.User
	NextCursor string
}

type PatchUser struct {
	ID         string
	Email      *string
	Name       *string
	AvatarURL  *string
	IsDisabled *bool
}

type UserRepo interface {
	Create(ctx context.Context, user *domain.User) (string, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, patch PatchUser) error
	List(ctx context.Context, req ListUsersReq) (ListUsersResp, error)
}

type ExternalIdentityRepo interface {
	Link(ctx context.Context, ident *domain.ExternalIdentity) error
}
