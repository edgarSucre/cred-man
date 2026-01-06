package credits

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/crm/internal/domain/event"
	"github.com/edgarSucre/mye"
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
) (ProcessCreditService, error) {
	if err := validateProcessService(bus, creditRepo, transactionManager); err != nil {
		return nil, err
	}

	return processCredit{bus, creditRepo, transactionManager}, nil
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
	err := mye.New(mye.CodeInvalid, "credit process failed", "validation error")

	if len(cmd.ClientID) == 0 {
		err.WithField("client_id", "client_id can't be empty")
	}

	if len(cmd.CreditID) == 0 {
		err.WithField("credit_id", "credit_id")
	}

	if err.HasFields() {
		return err
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

func validateProcessService(
	bus eventBus,
	creditRepo creditRepository,
	transactionManager transactionManager,
) error {
	err := mye.New(
		mye.CodeInternal,
		"processCredit_service_config_error",
		"processService parameter validation failed",
	)

	if bus == nil {
		err.WithField("bus", "event bus is missing")
	}

	if creditRepo == nil {
		err.WithField("creditRepo", "credit repository is missing")
	}

	if transactionManager == nil {
		err.WithField("transactionManager", "transaction manager is missing")
	}

	if err.HasFields() {
		return err
	}

	return nil
}
