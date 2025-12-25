package sheet

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
)

// Update updates a sheet using the callback pattern with optimistic locking
func (r *sheetRepo) Update(ctx context.Context, id string, fn func(*domain.Sheet) error) (*domain.Sheet, error) {
	ctx, span := tracer.Start(ctx, "SheetRepo.Update")
	defer span.End()

	doc := r.collection.Doc(id)
	var out *domain.Sheet

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(doc)
		if err != nil {
			return mapFirestoreError(err, "get sheet for update")
		}

		var cur domain.Sheet
		if err := snap.DataTo(&cur); err != nil {
			return fmt.Errorf("unmarshal sheet: %w", err)
		}
		if cur.ID == "" {
			cur.ID = snap.Ref.ID
		}

		before := cur

		// Apply patch function
		if err := fn(&cur); err != nil {
			return err
		}

		// Build diff
		updates := buildDiff(before, cur)
		if len(updates) == 0 {
			out = &cur
			return nil // no-op
		}

		// Add updated_at with actual timestamp
		now := time.Now().UTC()
		cur.UpdatedAt = now
		updates = append(updates, firestore.Update{Path: "updated_at", Value: now})

		// Optimistic concurrency
		pre := firestore.LastUpdateTime(snap.UpdateTime)
		if err := tx.Update(doc, updates, pre); err != nil {
			return mapFirestoreError(err, "update sheet")
		}

		out = &cur
		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return out, nil
}

// buildDiff computes the differences between two sheet instances
func buildDiff(before, after domain.Sheet) []firestore.Update {
	var updates []firestore.Update

	if before.Status != after.Status {
		updates = append(updates, firestore.Update{Path: "status", Value: after.Status})
	}
	if before.DeliveryFee != after.DeliveryFee {
		updates = append(updates, firestore.Update{Path: "delivery_fee", Value: after.DeliveryFee})
	}
	if before.Discount != after.Discount {
		updates = append(updates, firestore.Update{Path: "discount", Value: after.Discount})
	}
	if before.Description != after.Description {
		updates = append(updates, firestore.Update{Path: "description", Value: after.Description})
	}
	// Note: MemberIDs should be updated via AddMember/RemoveMember methods
	// to keep subcollection in sync, not through Update patch function

	return updates
}
