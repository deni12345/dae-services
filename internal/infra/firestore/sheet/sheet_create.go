package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-core/internal/domain"
)

// Create creates a new sheet in Firestore
func (r *sheetRepo) Create(ctx context.Context, sheet *domain.Sheet) (*domain.Sheet, error) {
	if sheet.ID == "" {
		return nil, fmt.Errorf("sheet ID is required")
	}

	docRef := r.collection.Doc(sheet.ID)
	_, err := docRef.Create(ctx, sheet)
	if err != nil {
		return nil, mapFirestoreError(err, "create sheet")
	}

	return sheet, nil
}
