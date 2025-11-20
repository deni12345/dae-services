package sheet

import (
	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/port"
)

type sheetRepo struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
}

// NewSheetRepo creates a new Firestore-backed sheet repository
func NewSheetRepo(client *firestore.Client) port.SheetRepo {
	return &sheetRepo{
		client:     client,
		collection: client.Collection("sheets"),
	}
}
