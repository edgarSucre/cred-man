package clients

import (
	"context"

	"github.com/edgarSucre/crm/internal/domain/client"
)

type Repository interface {
	CreateClient(context.Context, client.Client) (client.Client, error)
	GetClient(context.Context, client.ID) (client.Client, error)
}
