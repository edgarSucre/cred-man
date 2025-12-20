package client

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/db/repository"
	"github.com/edgarSucre/crm/pkg"
)

func (svc Service) CreateClient(ctx context.Context, params pkg.CreateClientParams) (*pkg.Client, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("create client validation-error: %w", err)
	}

	model := new(repository.CreateClientParams)
	model.FromDomain(params)

	newClient, err := svc.repo.CreateClient(ctx, *model)

	if err != nil {
		return nil, fmt.Errorf("create client insert-error: %w", err)
	}

	return newClient.ToDomain(), nil
}
