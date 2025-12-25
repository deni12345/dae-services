package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func NewFirestoreClient(ctx context.Context, projectID string) (*firestore.Client, error) {
	fbConfig := &firebase.Config{
		ProjectID: projectID,
	}

	app, err := firebase.NewApp(ctx, fbConfig)
	if err != nil {
		return nil, fmt.Errorf("init firebase app on err %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firestore client on err %w", err)
	}
	return client, nil
}
