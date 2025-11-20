package port

import (
	"context"

	"github.com/deni12345/dae-core/internal/domain"
)

type ListUserQuery struct {
	Limit           int32
	Cursor          string
	IncludeDisabled bool
}

type UpdateUserRequest struct {
	ID         string
	UserName   *string
	AvatarURL  *string
	IsDisabled *bool
}

type UsersRepo interface {
	// User
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id string, fn func(u *domain.User) error) (*domain.User, error)
	List(ctx context.Context, query ListUserQuery) ([]*domain.User, error)

	// Admin
	SetRoles(ctx context.Context, id string, roles []domain.Role) (*domain.User, error)
	SetDisabled(ctx context.Context, id string, isDisabled bool) (*domain.User, error)
}
