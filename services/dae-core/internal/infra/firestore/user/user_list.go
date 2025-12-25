package user

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/domain"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"google.golang.org/api/iterator"
)

func (r *userRepo) List(ctx context.Context, query port.ListUserQuery) ([]*domain.User, error) {
	ctx, span := tracer.Start(ctx, "UserRepo.List")
	defer span.End()

	var resp []*domain.User

	q := r.collection.Query
	q = q.Where("IsDisabled", "==", query.IncludeDisabled)

	q = q.OrderBy(firestore.DocumentID, firestore.Desc)
	if query.Cursor != "" {
		q = q.StartAfter(query.Cursor)
	}

	limit := query.Limit
	if limit <= 0 || limit > 1000 {
		limit = r.defaultPageSize
	}
	q = q.Limit(int(limit))

	iter := q.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			span.RecordError(err)
			return resp, fmt.Errorf("list users: %w", err)
		}

		var user domain.User
		if err := doc.DataTo(&user); err != nil {
			span.RecordError(err)
			return resp, fmt.Errorf("data to user: %w", err)
		}
		if user.ID == "" {
			user.ID = doc.Ref.ID
		}
		resp = append(resp, &user)
	}

	return resp, nil
}
