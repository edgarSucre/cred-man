package credits

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/mye"
)

type getCredit struct {
	creditRepo creditRepository
}

func NewGetCreditService(creditRepo creditRepository) GetCreditService {
	return getCredit{creditRepo}
}

func (svc getCredit) Execute(ctx context.Context, cmd GetCreditCommand) (CreditResult, error) {
	if err := cmd.validate(); err != nil {
		return CreditResult{}, fmt.Errorf("GetCreditCommand.validate > %w", err)
	}

	uID, err := credit.NewIDFromString(cmd.ID)
	if err != nil {
		return CreditResult{}, fmt.Errorf("getCredit > %w", err)
	}

	dCredit, err := svc.creditRepo.GetCredit(ctx, uID)
	if err != nil {
		return CreditResult{}, fmt.Errorf("getCredit > %w", err)
	}

	return CreditResult{
		BankID:     dCredit.BankID().String(),
		ClientID:   dCredit.ClientID().String(),
		CreatedAt:  dCredit.CreatedAt(),
		CreditType: dCredit.CreditType().String(),
		ID:         dCredit.ID().String(),
		MaxPayment: dCredit.MaxPayment().String(),
		MinPayment: dCredit.MinPayment().String(),
		Status:     dCredit.Status().String(),
		TermMonths: dCredit.TermMonths(),
	}, nil
}

type GetCreditCommand struct {
	ID string
}

func (cmd GetCreditCommand) validate() error {
	if len(cmd.ID) == 0 {
		return mye.New(mye.CodeInternal, "get_credit_failed", "validation error").
			WithUserMsg("can't get credit, due to input validation").
			WithField("id", "id can't be empty")
	}

	return nil
}
