package sheet

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"google.golang.org/api/iterator"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// ListForUser returns sheets where user is host OR member using CollectionGroup query
func (r *sheetRepo) ListForUser(ctx context.Context, query port.ListSheetsForUserQuery) (*port.ListSheetsForUserResp, error) {
	ctx, span := tracer.Start(ctx, "SheetRepo.ListForUser")
	defer span.End()

	var resp port.ListSheetsForUserResp

	// Validate & set limit
	limit := int(query.Limit)
	if limit <= 0 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}

	// Step 1: Query members subcollection with CollectionGroup
	// Optimized: server-side filtering by user_id field
	q := r.client.CollectionGroup("members").
		Where("user_id", "==", query.UserID).
		OrderBy("joined_at", firestore.Desc).
		Limit(limit + 1) // fetch one extra to determine if there's more

	if query.Cursor != "" {
		// For proper pagination with cursor
		q = q.StartAfter(query.Cursor)
	}

	iter := q.Documents(ctx)
	defer iter.Stop()

	sheetIDs := make([]string, 0, limit)
	var lastCursor string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("iterate members: %w", err)
		}

		// Extract sheet_id from path: sheets/{sheetID}/members/{userID}
		sheetRef := doc.Ref.Parent.Parent
		sheetID := sheetRef.ID

		sheetIDs = append(sheetIDs, sheetID)

		// Store cursor (joined_at value for next page)
		if data := doc.Data(); data != nil {
			if joinedAt, ok := data["joined_at"].(time.Time); ok {
				lastCursor = joinedAt.Format(time.RFC3339Nano)
			}
		}

		if len(sheetIDs) > limit {
			break
		}
	}

	// Check if there are more results
	hasMore := len(sheetIDs) > limit
	if hasMore {
		sheetIDs = sheetIDs[:limit]
		resp.NextCursor = lastCursor
	}

	if len(sheetIDs) == 0 {
		resp.Sheets = []*domain.Sheet{}
		return &resp, nil
	}

	// Step 2: Batch get sheets
	// Firestore supports up to 500 docs in GetAll, we're limiting to MaxPageSize
	sheets := make([]*domain.Sheet, 0, len(sheetIDs))
	for _, id := range sheetIDs {
		snap, err := r.collection.Doc(id).Get(ctx)
		if err != nil {
			// Sheet might be deleted - skip
			continue
		}

		var sheet domain.Sheet
		if err := snap.DataTo(&sheet); err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("unmarshal sheet %s: %w", id, err)
		}
		if sheet.ID == "" {
			sheet.ID = snap.Ref.ID
		}

		// Apply status filter if specified
		if query.StatusFilter != nil && sheet.Status != *query.StatusFilter {
			continue
		}

		sheets = append(sheets, &sheet)
	}

	resp.Sheets = sheets
	return &resp, nil
}
