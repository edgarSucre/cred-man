package client

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/google/uuid"
)

func (svc Service) GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	cl, err := svc.repo.GetClient(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("svc.GetClient: %w", err)
	}

	return cl.ToDomain(), nil
}
