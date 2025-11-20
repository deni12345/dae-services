package user

import (
	"context"
	"strings"

	"github.com/deni12345/dae-core/internal/domain"
)

// UpdateUser updates user profile information
func (u *usecase) UpdateUser(ctx context.Context, req *UpdateUserReq) (*domain.User, error) {
	if req.ID == "" {
		return nil, ErrInvalidArgument
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
		return nil, err
	}

	return user, nil
}
