package order

import (
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/port"
)

// Repository errors
var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderExists      = errors.New("order already exists")
	ErrConcurrentUpdate = errors.New("concurrent update detected")
)

type orderRepo struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
}

func NewOrderRepo(client *firestore.Client) port.OrdersRepo {
	return &orderRepo{
		client:     client,
		collection: client.Collection("orders"),
	}
}
