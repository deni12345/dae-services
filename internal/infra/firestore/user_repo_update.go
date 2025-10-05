package firestore

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/port"
)

type applicator func([]firestore.Update, port.PatchUser) []firestore.Update

var updateStrategies = []applicator{
	func(u []firestore.Update, p port.PatchUser) []firestore.Update {
		if p.Email != nil {
			u = append(u, firestore.Update{Path: "Email", Value: *p.Email})
		}
		return u
	},
	func(u []firestore.Update, p port.PatchUser) []firestore.Update {
		if p.Name != nil {
			u = append(u, firestore.Update{Path: "Name", Value: *p.Name})
		}
		return u
	},
	func(u []firestore.Update, p port.PatchUser) []firestore.Update {
		if p.AvatarURL != nil {
			u = append(u, firestore.Update{Path: "AvatarURL", Value: *p.AvatarURL})
		}
		return u
	},
	func(u []firestore.Update, p port.PatchUser) []firestore.Update {
		if p.IsDisabled != nil {
			u = append(u, firestore.Update{Path: "IsDisabled", Value: *p.IsDisabled})
		}
		return u
	},
}

func (r *userRepo) Update(ctx context.Context, patch port.PatchUser) error {
	updatesFields := make([]firestore.Update, 0, 4)
	for _, apply := range updateStrategies {
		updatesFields = apply(updatesFields, patch)
	}

	if len(updatesFields) == 0 {
		return nil
	}

	updatesFields = append(updatesFields, firestore.Update{Path: "UpdatedAt", Value: time.Now().UTC()})
	_, err := r.collection.Doc(patch.ID).Update(ctx, updatesFields)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
