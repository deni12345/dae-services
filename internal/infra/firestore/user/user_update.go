package user

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *userRepo) Update(ctx context.Context, id string, fn func(u *domain.User) error) (*domain.User, error) {
	doc := r.collection.Doc(id)
	var out *domain.User

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(doc)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return ErrNotFound
			}
			return fmt.Errorf("get user: %w", err)
		}

		var cur domain.User
		if err := snap.DataTo(&cur); err != nil {
			return fmt.Errorf("unmarshal user: %w", err)
		}
		if cur.ID == "" {
			cur.ID = snap.Ref.ID
		}

		before := cur

		// Apply patch function (business rules from usecase)
		if err := fn(&cur); err != nil {
			return err
		}

		// Build diff
		updates := buildUserDiff(before, cur)
		if len(updates) == 0 {
			out = &cur
			return nil // no-op
		}

		// Add updated_at
		now := time.Now().UTC()
		cur.UpdatedAt = now
		updates = append(updates, firestore.Update{Path: "updated_at", Value: now})

		// Optimistic concurrency
		pre := firestore.LastUpdateTime(snap.UpdateTime)
		if err := tx.Update(doc, updates, pre); err != nil {
			if status.Code(err) == codes.FailedPrecondition {
				return ErrConcurrentUpdate
			}
			return fmt.Errorf("update user: %w", err)
		}

		out = &cur
		return nil
	})

	if err != nil {
		return nil, err
	}
	return out, nil
}

// buildUserDiff computes differences between two user instances
func buildUserDiff(before, after domain.User) []firestore.Update {
	var updates []firestore.Update

	if before.UserName != after.UserName {
		updates = append(updates, firestore.Update{Path: "username", Value: after.UserName})
	}
	if before.AvatarURL != after.AvatarURL {
		updates = append(updates, firestore.Update{Path: "avatar_url", Value: after.AvatarURL})
	}
	if before.IsDisabled != after.IsDisabled {
		updates = append(updates, firestore.Update{Path: "is_disabled", Value: after.IsDisabled})
	}

	// Compare roles
	if len(before.Roles) != len(after.Roles) {
		updates = append(updates, firestore.Update{Path: "roles", Value: after.Roles})
	} else {
		// Check if roles changed
		rolesChanged := false
		for i := range before.Roles {
			if before.Roles[i] != after.Roles[i] {
				rolesChanged = true
				break
			}
		}
		if rolesChanged {
			updates = append(updates, firestore.Update{Path: "roles", Value: after.Roles})
		}
	}

	return updates
}

// SetRoles updates user roles (admin operation)
func (r *userRepo) SetRoles(ctx context.Context, userID string, roles []domain.Role) (*domain.User, error) {
	return r.Update(ctx, userID, func(u *domain.User) error {
		u.Roles = roles
		return nil
	})
}

// SetDisabled updates user disabled status (admin operation)
func (r *userRepo) SetDisabled(ctx context.Context, userID string, isDisabled bool) (*domain.User, error) {
	return r.Update(ctx, userID, func(u *domain.User) error {
		u.IsDisabled = isDisabled
		return nil
	})
}
