package domain

import (
	"time"
)

type Role string

const (
	RoleUnspecified Role = ""
	RoleUser        Role = "user"
	RoleHost        Role = "host"
	RoleAdmin       Role = "admin"
	RoleSuperAdmin  Role = "superadmin"
)

var ValidRoles = map[Role]bool{
	RoleUser:  true,
	RoleAdmin: true,
	RoleHost:  true,
}

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type IdentityProvider string

const (
	IdentityProviderGoogle IdentityProvider = "google"
	IdentityProviderLocal  IdentityProvider = "local"
)

type User struct {
	ID              string     `firestore:"-" json:"id"`
	Email           string     `firestore:"email" json:"email"`
	EmailNormalized string     `firestore:"email_normalized" json:"email_normalized"`
	EmailVerified   bool       `firestore:"email_verified" json:"email_verified"`
	Name            string     `firestore:"name" json:"name"`
	DisplayName     string     `firestore:"display_name" json:"display_name"`
	PhotoURL        string     `firestore:"photo_url" json:"photo_url"`
	Phone           string     `firestore:"phone" json:"phone"`
	Roles           []Role     `firestore:"roles" json:"roles"`
	Status          UserStatus `firestore:"status" json:"status"`
	CreatedAt       time.Time  `firestore:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `firestore:"updated_at" json:"updated_at"`
	LastLoginAt     *time.Time `firestore:"last_login_at,omitempty" json:"last_login_at,omitempty"`

	// Password-based authentication (optional)
	PasswordHash      *string    `firestore:"password_hash,omitempty" json:"-"`
	PasswordAlgo      *string    `firestore:"password_algo,omitempty" json:"-"`
	PasswordUpdatedAt *time.Time `firestore:"password_updated_at,omitempty" json:"password_updated_at,omitempty"`

	// Legacy fields for backward compatibility
	UserName   string `firestore:"user_name,omitempty" json:"user_name,omitempty"`
	AvatarURL  string `firestore:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	IsDisabled bool   `firestore:"is_disabled,omitempty" json:"is_disabled,omitempty"`
}

type UserIdentity struct {
	ID            string           `firestore:"-" json:"id"` // {uid}/identities/{provider}
	UserID        string           `firestore:"user_id" json:"user_id"`
	Provider      IdentityProvider `firestore:"provider" json:"provider"`
	Subject       string           `firestore:"subject" json:"subject"` // Provider's user ID
	EmailAtSignup string           `firestore:"email_at_signup" json:"email_at_signup"`
	LinkedAt      time.Time        `firestore:"linked_at" json:"linked_at"`
	LastLoginAt   *time.Time       `firestore:"last_login_at,omitempty" json:"last_login_at,omitempty"`
}

type UniqueEmail struct {
	Email     string    `firestore:"-" json:"email"` // doc id
	UserID    string    `firestore:"user_id" json:"user_id"`
	CreatedAt time.Time `firestore:"created_at" json:"created_at"`
}

type UniqueIdentity struct {
	ProviderSubject string    `firestore:"-" json:"provider_subject"` // doc id = {provider}:{subject}
	UserID          string    `firestore:"user_id" json:"user_id"`
	CreatedAt       time.Time `firestore:"created_at" json:"created_at"`
}
