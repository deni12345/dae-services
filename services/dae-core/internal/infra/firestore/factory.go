package firestore

import (
	"cloud.google.com/go/firestore"
	"github.com/deni12345/dae-services/services/dae-core/internal/infra/firestore/order"
	"github.com/deni12345/dae-services/services/dae-core/internal/infra/firestore/sheet"
	"github.com/deni12345/dae-services/services/dae-core/internal/infra/firestore/user"
	"github.com/deni12345/dae-services/services/dae-core/internal/port"
)

func NewUserRepo(client *firestore.Client, defaultPageSize int32) port.UsersRepo {
	return user.NewUserRepo(client, defaultPageSize)
}

func NewOrderRepo(client *firestore.Client, defaultPageSize int32) port.OrdersRepo {
	return order.NewOrderRepo(client, defaultPageSize)
}

func NewSheetRepo(client *firestore.Client, defaultPageSize int32) port.SheetRepo {
	return sheet.NewSheetRepo(client, defaultPageSize)
}
