package clients

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/client"
)

type getClient struct {
	repo Repository
}

func NewGetClientService(repo Repository) GetClientService {
	return getClient{repo: repo}
}

func (svc getClient) GetClient(ctx context.Context, cmd GetClientCommand) (ClientResult, error) {
	if err := cmd.Validate(); err != nil {
		return ClientResult{}, fmt.Errorf("getClient.CommandValidation: %w", err)
	}

	clientID, err := client.NewID(cmd.ID)
	if err != nil {
		return ClientResult{}, fmt.Errorf("client.NewID > %w", err)
	}

	cl, err := svc.repo.GetClient(ctx, clientID)
	if err != nil {
		return ClientResult{}, fmt.Errorf("getClient.repo.GetClient: %w", err)
	}

	return mapClientResult(cl), nil
}

type GetClientCommand struct {
	ID string
}

func (cmd GetClientCommand) Validate() error {
	if len(cmd.ID) == 0 {
		// TODO: replace error with app error
		return client.ErrInvalidClientID
	}

	return nil
}
