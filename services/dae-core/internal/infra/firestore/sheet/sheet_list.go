package sheet

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
)

// List retrieves a paginated list of sheets (not yet implemented)
func (r *sheetRepo) List(ctx context.Context, query port.ListSheetsQuery) ([]*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetRepo.List")
	defer span.End()

	limit := query.Limit
	if limit <= 0 || limit > 1000 {
		limit = r.defaultPageSize
	}
	q := r.collection.Query
	q = q.Limit(int(limit))

	if query.Cursor != "" {
		cursorSnap, err := r.collection.Doc(query.Cursor).Get(ctx)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("get cursor document: %w", err)
		}
		q = q.StartAfter(cursorSnap)
	}

	// Add ordering (most recent first)
	q = q.OrderBy("created_at", firestore.Desc)

	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("list sheets: %w", err)
	}

	sheets := make([]*domain.Sheet, 0, len(docs))
	for _, doc := range docs {
		var sheet domain.Sheet
		if err := doc.DataTo(&sheet); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("unmarshal sheet: %w", err)
		}
		if sheet.ID == "" {
			sheet.ID = doc.Ref.ID
		}
		sheets = append(sheets, &sheet)
	}

	return sheets, nil
}
