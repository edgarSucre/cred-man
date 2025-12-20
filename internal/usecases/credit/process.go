package credit

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
)

func (svc *Service) ProcessCredit(ctx context.Context, msg domain.CreditCreated) error {
	credits, err := svc.repo.GetClientCredits(ctx, msg.ClientID)
	if err != nil {
		return fmt.Errorf("repo.GetCredits: %w", err)
	}

	newStatus := repository.CreditStatusApproved

	if shouldRejectCredit(credits, msg.CreditID) {
		newStatus = repository.CreditStatusRejected
	}

	return svc.repo.UpdateCreditStatus(ctx, newStatus)
}

// my completely arbitrary rules
func shouldRejectCredit(credits []repository.Credit, creditID uuid.UUID) bool {
	var creditToProcess repository.Credit

	mortgages := make([]repository.Credit, 0, len(credits))
	autos := make([]repository.Credit, 0, len(credits))
	commercial := make([]repository.Credit, 0, len(credits))

	for _, v := range credits {
		if v.ID == creditID {
			creditToProcess = v

			continue
		}

		if v.CreditType == repository.CreditTypeMortgage {
			mortgages = append(mortgages, v)

			continue
		}

		if v.CreditType == repository.CreditTypeAuto {
			autos = append(mortgages, v)

			continue
		}

		if v.CreditType == repository.CreditTypeCommercial {
			commercial = append(commercial, v)
		}
	}

	if creditToProcess.CreditType == repository.CreditTypeMortgage && len(mortgages) >= 4 {
		return true
	}

	if creditToProcess.CreditType == repository.CreditTypeAuto && len(autos) >= 2 {
		return true
	}

	if creditToProcess.CreditType == repository.CreditTypeCommercial && len(commercial) >= 3 {
		return true
	}

	return false
}

func errCreditNotFound(id uuid.UUID) error {
	detail := fmt.Sprintf("can't find credit with id: %s", id)

	return terror.NotFound.New("credit-not-found", detail)
}
