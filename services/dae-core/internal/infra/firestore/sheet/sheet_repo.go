package sheet

import (
	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"go.opentelemetry.io/otel"
)

type sheetRepo struct {
	client          *firestore.Client
	collection      *firestore.CollectionRef
	defaultPageSize int32
}

var tracer = otel.Tracer("firestore/sheet")

// NewSheetRepo creates a new Firestore-backed sheet repository
func NewSheetRepo(client *firestore.Client, defaultPageSize int32) port.SheetRepo {
	return &sheetRepo{
		client:          client,
		collection:      client.Collection("sheets"),
		defaultPageSize: defaultPageSize,
	}
}
