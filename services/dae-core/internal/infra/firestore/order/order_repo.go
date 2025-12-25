package order

import (
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"go.opentelemetry.io/otel"
)

// Repository errors
var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderExists      = errors.New("order already exists")
	ErrConcurrentUpdate = errors.New("concurrent update detected")
	tracer              = otel.Tracer("firestore/order")
)

type orderRepo struct {
	client          *firestore.Client
	collection      *firestore.CollectionRef
	defaultPageSize int32
}

func NewOrderRepo(client *firestore.Client, defaultPageSize int32) port.OrdersRepo {
	return &orderRepo{
		client:          client,
		collection:      client.Collection("orders"),
		defaultPageSize: defaultPageSize,
	}
}
