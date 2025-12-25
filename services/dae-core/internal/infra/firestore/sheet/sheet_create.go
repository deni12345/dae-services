package sheet

import (
	"context"
	"fmt"

	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
)

// Create creates a new sheet in Firestore
func (r *sheetRepo) Create(ctx context.Context, sheet *domain.Sheet) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetRepo.Create")
	defer span.End()

	if sheet.ID == "" {
		err := fmt.Errorf("sheet ID is required")
		span.RecordError(err)
		return nil, err
	}

	docRef := r.collection.Doc(sheet.ID)
	_, err := docRef.Create(ctx, sheet)
	if err != nil {
		span.RecordError(err)
		return nil, mapFirestoreError(err, "create sheet")
	}

	return sheet, nil
}
