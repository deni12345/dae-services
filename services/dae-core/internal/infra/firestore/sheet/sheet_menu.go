package sheet

import (
	"context"
	"errors"
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"google.golang.org/api/iterator"
)

func (r *sheetRepo) GetMenuItemByID(ctx context.Context, sheetID string, id string) (*domain.MenuItem, error) {
	doc, err := r.collection.Doc(sheetID).Collection("menu").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var menuItem domain.MenuItem
	if err := doc.DataTo(&menuItem); err != nil {
		return nil, err
	}

	return &menuItem, nil
}

func (r *sheetRepo) GetMenuItems(ctx context.Context, sheetID string) ([]*domain.MenuItem, error) {
	iter := r.collection.Doc(sheetID).Collection("menu").Documents(ctx)
	defer iter.Stop()

	var menuItems []*domain.MenuItem
	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			return nil, err
		}

		var menuItem domain.MenuItem
		if err := doc.DataTo(&menuItem); err != nil {
			return nil, err
		}
		menuItems = append(menuItems, &menuItem)
	}

	return menuItems, nil
}

// AttachMenuItems attaches menu items to a sheet's menu subcollection
// Handles batching automatically for lists larger than 500 items
func (r *sheetRepo) AttachMenuItems(ctx context.Context, sheetID string, menuItems []*domain.MenuItem) error {
	if len(menuItems) == 0 {
		return nil // no-op
	}

	menuCollection := r.collection.Doc(sheetID).Collection("menu")

	// Firestore batch limit is 500 operations
	const batchSize = 500

	for i := 0; i < len(menuItems); i += batchSize {
		end := i + batchSize
		if end > len(menuItems) {
			end = len(menuItems)
		}

		batch := r.client.Batch()
		chunk := menuItems[i:end]

		for _, item := range chunk {
			if item.ID == "" {
				return fmt.Errorf("menu item ID is required")
			}

			menuRef := menuCollection.Doc(item.ID)
			batch.Set(menuRef, item)
		}

		if _, err := batch.Commit(ctx); err != nil {
			return fmt.Errorf("attach menu items (batch %d-%d): %w", i, end, err)
		}
	}

	return nil
}
