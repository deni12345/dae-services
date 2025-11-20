package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// GetByID retrieves a sheet by ID
func (r *sheetRepo) GetByID(ctx context.Context, id string) (*domain.Sheet, error) {
	snap, err := r.collection.Doc(id).Get(ctx)
	if err != nil {
		return nil, mapFirestoreError(err, "get sheet")
	}

	var sheet domain.Sheet
	if err := snap.DataTo(&sheet); err != nil {
		return nil, fmt.Errorf("unmarshal sheet: %w", err)
	}

	if sheet.ID == "" {
		sheet.ID = snap.Ref.ID
	}

	return &sheet, nil
}
