package credit

import (
	"context"
	"fmt"
	"sync"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/shopspring/decimal"
)

func (svc *Service) CreateCredit(
	ctx context.Context,
	params domain.CreateCreditParams,
) error {
	if err := params.Validate(); err != nil {
		return fmt.Errorf("create credit validation-error: %w", err)
	}

	model := new(repository.CreateCreditParams)
	model.FromDomain(params)
	setModel(model)

	return svc.repo.WithTransaction(ctx, svc.createTx(ctx, *model))
}

func (svc *Service) createTx(
	ctx context.Context,
	model repository.CreateCreditParams,
) func(*repository.Queries) error {
	fn := func(q *repository.Queries) error {
		failureCh := make(chan error, 3)

		var wg sync.WaitGroup
		wg.Add(2)

		// check if bank exists
		go func() {
			defer wg.Done()

			_, err := q.GetBank(ctx, model.BankID)
			if err != nil {
				failureCh <- err
			}
		}()

		// check if client exists
		go func() {
			defer wg.Done()

			_, err := q.GetClient(ctx, model.ClientID)
			if err != nil {
				failureCh <- err
			}
		}()

		wg.Wait()
		if len(failureCh) > 0 {
			err := <-failureCh

			return fmt.Errorf("crate credit transaction error: %w", err)
		}

		credit, err := q.CreateCredit(ctx, model)
		if err != nil {
			return fmt.Errorf("repo.CreateCredit error: %w", err)
		}

		return svc.eventBus.Publish(ctx, domain.CreditCreated{
			BankID:   credit.BankID,
			ClientID: credit.ClientID,
			CreditID: credit.ID,
		})
	}

	return fn
}

func setModel(model *repository.CreateCreditParams) {
	model.Status = repository.CreditStatusPending

	// made up values
	switch model.CreditType {
	case repository.CreditTypeAuto:
		model.MinPayment = decimal.NewFromInt(1000)
		model.MaxPayment = decimal.NewFromInt(5000)
		model.TermMonths = 60
	case repository.CreditTypeMortgage:
		model.MinPayment = decimal.NewFromInt(2000)
		model.MaxPayment = decimal.NewFromInt(10000)
		model.TermMonths = 360
	default:
		model.MinPayment = decimal.NewFromInt(10000)
		model.MaxPayment = decimal.NewFromInt(50000)
		model.TermMonths = 120
	}
}
