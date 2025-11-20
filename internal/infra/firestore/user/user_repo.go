package user

import (
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-core/internal/port"
)

// Repository errors
var (
	ErrNotFound         = errors.New("user not found")
	ErrAlreadyExists    = errors.New("user already exists")
	ErrConcurrentUpdate = errors.New("concurrent update detected")
)

type userRepo struct {
	client     *firestore.Client
	collection *firestore.CollectionRef
}

func NewUserRepo(client *firestore.Client) port.UsersRepo {
	return &userRepo{
		client:     client,
		collection: client.Collection("users"),
	}
}
