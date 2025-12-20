package bank

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/pkg/domain"
)

func (svc *Service) CreateBank(
	ctx context.Context,
	params domain.CreateBankParams,
) (*domain.Bank, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("create bank validation-error: %w", err)
	}

	model := new(repository.CreateBankParams)
	model.FromDomain(params)

	newBank, err := svc.repo.CreateBank(ctx, *model)

	if err != nil {
		return nil, fmt.Errorf("create bank error: %w", err)
	}

	return newBank.ToDomain(), nil
}
