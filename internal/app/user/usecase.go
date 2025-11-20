package user

import (
	"context"

	"github.com/deni12345/dae-core/internal/domain"
	"github.com/deni12345/dae-core/internal/port"
)

// Usecase defines all user operations
type Usecase interface {
	// Queries
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, req *ListUsersReq) (*ListUsersResp, error)

	// Commands
	UpdateUser(ctx context.Context, req *UpdateUserReq) (*domain.User, error)

	// Admin operations
	AdminSetUserRoles(ctx context.Context, req *AdminSetUserRolesReq) (*domain.User, error)
	AdminSetUserDisabled(ctx context.Context, req *AdminSetUserDisabledReq) (*domain.User, error)
}

type usecase struct {
	userRepo port.UsersRepo
}

// NewUsecase creates a new user usecase
func NewUsecase(userRepo port.UsersRepo) Usecase {
	return &usecase{
		userRepo: userRepo,
	}
}
