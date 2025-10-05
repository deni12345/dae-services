package firestore

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/deni12345/dae-core/internal/configs"
)

func NewFirestoreInstance(ctx context.Context) (*firestore.Client, error) {
	fbConfig := &firebase.Config{ProjectID: configs.Values.Firestore.ProjectID}
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
