package credit

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/google/uuid"
)

func (svc *Service) GetCredit(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Credit, error) {
	credit, err := svc.repo.GetCredit(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("svc.GetClient: %w", err)
	}

	return credit.ToDomain(), nil
}
