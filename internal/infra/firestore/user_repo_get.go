package firestore

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	doc, err := r.collection.Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("data to user: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	doc, err := r.collection.Where("Email", "==", email).Limit(1).Documents(ctx).Next()
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	var user domain.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("data to user: %w", err)
	}
	return &user, nil
}
