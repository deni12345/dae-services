package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/infra/firestore/order"
	"github.com/deni12345/dae-core/internal/infra/firestore/sheet"
	"github.com/deni12345/dae-core/internal/infra/firestore/user"
	"github.com/deni12345/dae-core/internal/port"
)

// NewUserRepo creates a new Firestore-backed user repository
func NewUserRepo(client *firestore.Client) port.UsersRepo {
	return user.NewUserRepo(client)
}

// NewOrderRepo creates a new Firestore-backed order repository
func NewOrderRepo(client *firestore.Client, projectID string) port.OrdersRepo {
	return order.NewOrderRepo(client)
}

// NewSheetRepo creates a new Firestore-backed sheet repository
func NewSheetRepo(client *firestore.Client) port.SheetRepo {
	return sheet.NewSheetRepo(client)
}
