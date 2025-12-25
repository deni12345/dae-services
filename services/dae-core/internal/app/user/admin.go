package user

import (
	"context"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/libs/apperror"
)

// AdminSetUserRoles updates user roles (admin only)
func (u *usecase) AdminSetUserRoles(ctx context.Context, req *AdminSetUserRolesReq) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUC.AdminSetUserRoles")
	defer span.End()

	if req.UserID == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return nil, err
	}

	// Validate roles
	for _, r := range req.Roles {
		if !domain.ValidRoles[r] {
			span.RecordError(ErrInvalidRole)
			return nil, ErrInvalidRole
		}
	}

	user, err := u.userRepo.SetRoles(ctx, req.UserID, req.Roles)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

// AdminSetUserDisabled enables/disables a user account (admin only)
func (u *usecase) AdminSetUserDisabled(ctx context.Context, req *AdminSetUserDisabledReq) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUC.AdminSetUserDisabled")
	defer span.End()

	if req.UserID == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return nil, err
	}

	user, err := u.userRepo.SetDisabled(ctx, req.UserID, req.IsDisabled)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}
