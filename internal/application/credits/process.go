package credits

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/crm/internal/domain/event"
)

type processCredit struct {
	bus                eventBus
	creditRepo         creditRepository
	transactionManager transactionManager
}

func NewProcessCreditService(
	bus eventBus,
	creditRepo creditRepository,
	transactionManager transactionManager,
) ProcessCreditService {
	return processCredit{bus, creditRepo, transactionManager}
}

func (svc processCredit) Execute(
	ctx context.Context,
	cmd ProcessCreditCommand,
) error {
	if err := cmd.validate(); err != nil {
		return fmt.Errorf("processCredit > %w", err)
	}

	creditID, err := credit.NewIDFromString(cmd.ClientID)
	if err != nil {
		return fmt.Errorf("processCredit > %w", err)
	}

	clientID, err := client.NewID(cmd.ClientID)
	if err != nil {
		return fmt.Errorf("processCredit > %w", err)
	}

	creditAggregate, err := svc.creditRepo.GetAggregate(ctx, creditID, clientID)
	if err != nil {
		return fmt.Errorf("creditRepo.GetAggregate > %w", err)
	}

	creditAggregate.Process()

	if err := svc.creditRepo.ProcessCredit(ctx, *creditAggregate); err != nil {
		return fmt.Errorf("creditRepo.ProcessCredit > %w", err)
	}

	svc.transactionManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := svc.processCreditAndPublishEvent(txCtx, *creditAggregate); err != nil {
			return fmt.Errorf("processCreditAndPublishEvent > %w", err)
		}

		return nil
	})

	return nil
}

type ProcessCreditCommand struct {
	ClientID string
	CreditID string
}

func (cmd ProcessCreditCommand) validate() error {
	if len(cmd.ClientID) == 0 {
		return fmt.Errorf("ProcessCreditCommand.validate > %w", ErrNoClientID)
	}

	if len(cmd.CreditID) == 0 {
		return fmt.Errorf("ProcessCreditCommand.validate > %w", ErrNoCreditID)
	}

	return nil
}

func (svc processCredit) processCreditAndPublishEvent(
	ctx context.Context,
	aggregate credit.CreditAggregate,
) error {
	if err := svc.creditRepo.ProcessCredit(ctx, aggregate); err != nil {
		return fmt.Errorf("creditRepo.ProcessCredit > %w", err)
	}

	var creditEvent event.Event

	if aggregate.Status() == credit.CreditStatusApproved {
		creditEvent = event.CreditApproved{
			CreditID: aggregate.ID().String(),
		}
	}

	if aggregate.Status() == credit.CreditStatusRejected {
		creditEvent = event.CreditRejected{
			CreditID: aggregate.ID().String(),
		}
	}

	if err := svc.bus.Publish(ctx, creditEvent); err != nil {
		return fmt.Errorf("eventBux.Publish > %w", err)
	}

	return nil
}
