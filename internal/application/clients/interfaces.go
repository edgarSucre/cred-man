package clients

import "context"

type CreateClientService interface {
	CreateClient(context.Context, CreateClientCommand) (ClientResult, error)
}

type GetClientService interface {
	GetClient(context.Context, GetClientCommand) (ClientResult, error)
}
