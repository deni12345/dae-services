package user

import (
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/libs/apperror"
)

// validateUser validates user domain rules
func validateUser(u *domain.User) error {
	if len(u.UserName) <= 0 || len(u.UserName) > 64 {
		return apperror.InvalidInput(fmt.Sprintf("username must be between 1 and 64 characters, got length %d", len(u.UserName)))
	}
	return nil
}
