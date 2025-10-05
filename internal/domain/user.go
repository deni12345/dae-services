package domain

import "time"

type Role string

const (
	RoleUnspecified Role = ""
	RoleUser        Role = "user"
	RoleAdmin       Role = "admin"
)

type User struct {
	ID         string
	Email      string
	Name       string
	Role       []Role
	IsDisabled bool
	AvatarURL  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ExternalIdentity struct {
	ID       string // provider:subject
	UserID   string
	Provider string // map từ enum
	Subject  string
	Email    string
	LinkedAt time.Time
}
