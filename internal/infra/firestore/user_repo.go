package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/domain"
	"github.com/deni12345/dae-core/internal/port"
	"google.golang.org/api/iterator"
)

const DEFAULT_PAGE_SIZE = 20

type userRepo struct {
	collection *firestore.CollectionRef
}

func NewUserRepo(client *firestore.Client) port.UserRepo {
	return &userRepo{
		collection: client.Collection("users"),
	}
}

func (r *userRepo) List(ctx context.Context, req port.ListUsersReq) (port.ListUsersResp, error) {
	var resp port.ListUsersResp

	query := r.collection.Query

	if req.ExactEmail != "" {
		query = query.Where("Email", "==", req.ExactEmail)
	}
	if !req.IncludeDisabled {
		query = query.Where("IsDisabled", "==", false)
	}

	query = query.OrderBy(firestore.DocumentID, firestore.Desc)

	if req.Cursor != "" {
		query = query.StartAfter(req.Cursor)
	}

	limit := req.Limit
	if limit <= 0 {
		limit = DEFAULT_PAGE_SIZE
	}
	query = query.Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()
	var lastID string

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return resp, fmt.Errorf("list users: %w", err)
		}
		var user domain.User
		if err := doc.DataTo(&user); err != nil {
			return resp, fmt.Errorf("data to user: %w", err)
		}
		if user.ID == "" {
			user.ID = doc.Ref.ID
		}
		resp.Users = append(resp.Users, &user)
		lastID = doc.Ref.ID
	}

	if len(resp.Users) == limit && lastID != "" {
		resp.NextCursor = lastID
	}
	return resp, nil
}
