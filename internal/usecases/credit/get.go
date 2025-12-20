package credit

import (
	"context"
	"fmt"
	"strings"

	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
)

var ErrNotFound = terror.NotFound.New("not-found", "credit not found")

func (svc *Service) GetCredit(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Credit, error) {
	credit, err := svc.repo.GetCredit(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("svc.GetClient: %w", err)
	}

	return credit.ToDomain(), nil
}
