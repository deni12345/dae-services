package user

import (
	"context"

	"github.com/deni12345/dae-core/internal/domain"
)

// AdminSetUserRoles updates user roles (admin only)
func (u *usecase) AdminSetUserRoles(ctx context.Context, req *AdminSetUserRolesReq) (*domain.User, error) {
	if req.UserID == "" {
		return nil, ErrInvalidArgument
	}

	// Validate roles
	for _, r := range req.Roles {
		if !domain.ValidRoles[r] {
			return nil, ErrInvalidRole
		}
	}

	user, err := u.userRepo.SetRoles(ctx, req.UserID, req.Roles)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AdminSetUserDisabled enables/disables a user account (admin only)
func (u *usecase) AdminSetUserDisabled(ctx context.Context, req *AdminSetUserDisabledReq) (*domain.User, error) {
	if req.UserID == "" {
		return nil, ErrInvalidArgument
	}

	user, err := u.userRepo.SetDisabled(ctx, req.UserID, req.IsDisabled)
	if err != nil {
		return nil, err
	}

	return user, nil
}
