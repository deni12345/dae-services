package port

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
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

type CreateUserRequest struct {
	Email       string
	Name        string
	DisplayName string
	PhotoURL    string
	Phone       string
	Provider    domain.IdentityProvider
	Subject     string
	Password    string // Only for local provider
}

type UsersRepo interface {
	// User CRUD
	Create(ctx context.Context, req CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id string, fn func(u *domain.User) error) (*domain.User, error)
	List(ctx context.Context, query ListUserQuery) ([]*domain.User, error)

	// Identity management
	CreateIdentity(ctx context.Context, userID string, identity *domain.UserIdentity) error
	GetIdentityByProvider(ctx context.Context, userID string, provider domain.IdentityProvider) (*domain.UserIdentity, error)

	// Uniqueness checks (for transactions)
	CheckEmailUnique(ctx context.Context, email string) (bool, error)
	CheckIdentityUnique(ctx context.Context, provider domain.IdentityProvider, subject string) (bool, error)

	// Admin
	SetRoles(ctx context.Context, id string, roles []domain.Role) (*domain.User, error)
	SetDisabled(ctx context.Context, id string, isDisabled bool) (*domain.User, error)
}
