package order

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *orderRepo) Update(ctx context.Context, id string, fn func(o *domain.Order) error) (*domain.Order, error) {
	ctx, span := tracer.Start(ctx, "OrderRepo.Update")
	defer span.End()

	docRef := r.collection.Doc(id)
	var out *domain.Order

	err := r.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		snap, err := tx.Get(docRef)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return ErrOrderNotFound
			}
			return fmt.Errorf("get order: %w", err)
		}

		var cur domain.Order
		if err := snap.DataTo(&cur); err != nil {
			return fmt.Errorf("unmarshal order: %w", err)
		}

		if cur.ID == "" {
			cur.ID = snap.Ref.ID
		}

		before := cur

		// Apply business logic from usecase
		if err := fn(&cur); err != nil {
			return err
		}

		// Build diff to avoid unnecessary updates
		updates := buildOrderDiff(before, cur)
		if len(updates) == 0 {
			out = &cur
			return nil // no-op
		}

		// Add updated_at
		now := time.Now().UTC()
		cur.UpdatedAt = now
		updates = append(updates, firestore.Update{Path: "updated_at", Value: now})

		// Optimistic concurrency
		pre := firestore.LastUpdateTime(snap.UpdateTime)
		if err := tx.Update(docRef, updates, pre); err != nil {
			if status.Code(err) == codes.FailedPrecondition {
				return ErrConcurrentUpdate
			}
			return fmt.Errorf("update order: %w", err)
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

// buildOrderDiff computes differences between two orders
func buildOrderDiff(before, after domain.Order) []firestore.Update {
	var updates []firestore.Update

	// Compare lines (simplified - always update if fn was called)
	if len(before.Lines) != len(after.Lines) || before.Note != after.Note {
		updates = append(updates,
			firestore.Update{Path: "lines", Value: after.Lines},
			firestore.Update{Path: "subtotal", Value: after.Subtotal},
			firestore.Update{Path: "total", Value: after.Total},
			firestore.Update{Path: "note", Value: after.Note},
		)
	}

	return updates
}
