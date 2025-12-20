package client

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/pkg/domain"
)

func (svc Service) CreateClient(ctx context.Context, params domain.CreateClientParams) (*domain.Client, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("create client validation-error: %w", err)
	}

	model := new(repository.CreateClientParams)
	model.FromDomain(params)

	newClient, err := svc.repo.CreateClient(ctx, *model)

	if err != nil {
		return nil, fmt.Errorf("create client error: %w", err)
	}

	return newClient.ToDomain(), nil
}
