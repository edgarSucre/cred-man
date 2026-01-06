package clients

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/mye"
)

type getClient struct {
	repo Repository
}

func NewGetClientService(repo Repository) (GetClientService, error) {
	if repo == nil {
		return nil, mye.New(
			mye.CodeInvalid,
			"getClient_service_config_error",
			"client repository is missing",
		)
	}
	return getClient{repo: repo}, nil
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
		return mye.New(mye.CodeInvalid, "get_client_failed", "validation failed").
			WithUserMsg("can't search client, validation failed").
			WithField("client_id", "client_id can't be empty")
	}

	return nil
}
