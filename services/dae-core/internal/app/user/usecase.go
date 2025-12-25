package user

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"go.opentelemetry.io/otel"
)

// Usecase defines all user operations
type Usecase interface {
	// Commands
	CreateUser(ctx context.Context, req *CreateUserReq) (*domain.User, error)
	UpdateUser(ctx context.Context, req *UpdateUserReq) (*domain.User, error)

	// Queries
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context, req *ListUsersReq) (*ListUsersResp, error)

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

var tracer = otel.Tracer("usecase/user")
