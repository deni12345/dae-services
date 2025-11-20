package domain

import (
	"strings"
	"time"
)

type Role string

const (
	RoleUnspecified Role = ""
	RoleUser        Role = "user"
	RoleOwner       Role = "owner"
	RoleAdmin       Role = "admin"
	RoleSuperAdmin  Role = "superadmin"
)

var ValidRoles = map[Role]bool{
	RoleUser:  true,
	RoleAdmin: true,
	RoleOwner: true,
}

type User struct {
	ID         string    `firestore:"-" json:"id"`
	Email      string    `firestore:"email" json:"email"`
	UserName   string    `firestore:"user_name" json:"user_name"`
	Roles      []Role    `firestore:"roles" json:"roles"`
	IsDisabled bool      `firestore:"is_disabled" json:"is_disabled"`
	AvatarURL  string    `firestore:"avatar_url" json:"avatar_url"`
	CreatedAt  time.Time `firestore:"created_at" json:"created_at"`
	UpdatedAt  time.Time `firestore:"updated_at" json:"updated_at"`
}

func NormalizeEmail(str string) string {
	str = strings.ToLower(str)
	return strings.TrimSpace(str)
}
