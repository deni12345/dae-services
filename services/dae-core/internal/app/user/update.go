package user

import (
	"context"
	"strings"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/libs/apperror"
)

// UpdateUser updates user profile information
func (u *usecase) UpdateUser(ctx context.Context, req *UpdateUserReq) (*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserUC.UpdateUser")
	defer span.End()

	if req.ID == "" {
		err := apperror.InvalidInput("user_id is required")
		span.RecordError(err)
		return nil, err
	}

	user, err := u.userRepo.Update(ctx, req.ID, func(user *domain.User) error {
		// Apply changes
		if req.UserName != nil && strings.TrimSpace(*req.UserName) != user.UserName {
			user.UserName = strings.TrimSpace(*req.UserName)
		}

		if req.AvatarURL != nil && *req.AvatarURL != user.AvatarURL {
			user.AvatarURL = *req.AvatarURL
		}

		if req.IsDisabled != nil && *req.IsDisabled != user.IsDisabled {
			user.IsDisabled = *req.IsDisabled
		}

		// Validate after applying changes
		return validateUser(user)
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}
