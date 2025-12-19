package client

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/pkg"
)

func (svc Service) CreateClient(ctx context.Context, params pkg.CreateClientParams) (*pkg.Client, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("create client validation-error: %w", err)
	}

	newClient, err := svc.repo.CreateClient(ctx, params.ToModel())

	if err != nil {
		return nil, fmt.Errorf("create client insert-error: %w", err)
	}

	resp := new(pkg.Client)
	resp.FromModel(newClient)

	return resp, nil
}
