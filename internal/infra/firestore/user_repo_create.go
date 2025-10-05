package firestore

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

func (r *userRepo) Create(ctx context.Context, user *domain.User) (string, error) {
	doc, _, err := r.collection.Add(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}
	return doc.ID, nil
}
