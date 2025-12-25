package user

import (
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
	"go.opentelemetry.io/otel"
)

// Repository errors
var (
	ErrNotFound         = errors.New("user not found")
	ErrAlreadyExists    = errors.New("user already exists")
	ErrConcurrentUpdate = errors.New("concurrent update detected")
)

type userRepo struct {
	client          *firestore.Client
	collection      *firestore.CollectionRef
	defaultPageSize int32
}

var tracer = otel.Tracer("firestore/user")

func NewUserRepo(client *firestore.Client, defaultPageSize int32) port.UsersRepo {
	return &userRepo{
		client:          client,
		collection:      client.Collection("users"),
		defaultPageSize: defaultPageSize,
	}
}
