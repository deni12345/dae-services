package user

import (
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// validateUser validates user domain rules
func validateUser(u *domain.User) error {
	if len(u.UserName) <= 0 || len(u.UserName) > 64 {
		return fmt.Errorf("username must be between 1 and 64 characters")
	}
	return nil
}
