package credits

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/crm/internal/domain/event"
	"github.com/edgarSucre/mye"
)

type createCredit struct {
	bankRepo   bankRepository
	clientRepo clientRepository
	creditRepo creditRepository
	eventBus   eventBus
	txManager  transactionManager
}

func NewCreateCreditService(
	bankRepo bankRepository,
	clientRepo clientRepository,
	creditRepo creditRepository,
	eventBus eventBus,
	txManager transactionManager,
) (CreateCreditService, error) {
	if err := validateCreateService(bankRepo, clientRepo, creditRepo, eventBus, txManager); err != nil {
		return nil, err
	}

	return createCredit{bankRepo, clientRepo, creditRepo, eventBus, txManager}, nil
}

func (svc createCredit) Execute(
	ctx context.Context,
	cmd CreateCreditCommand,
) (CreditResult, error) {
	if err := cmd.validate(); err != nil {
		return CreditResult{}, fmt.Errorf("CreateCredit > %w", err)
	}

	newCreditOpts, err := newCreditOpts(cmd)
	if err := cmd.validate(); err != nil {
		return CreditResult{}, fmt.Errorf("CreateCredit > %w", err)
	}

	_, _, err = svc.getBankAndClient(ctx, newCreditOpts.ClientID, newCreditOpts.BankID)
	if err != nil {
		return CreditResult{}, fmt.Errorf("CreateCredit > %w", err)
	}

	newCredit, err := credit.New(newCreditOpts)
	if err != nil {
		return CreditResult{}, fmt.Errorf("createCredit > %w", err)
	}

	var createdCredit credit.Credit

	err = svc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		createdCredit, err = svc.createCreditAndPublishEvent(txCtx, newCredit)
		if err != nil {
			return fmt.Errorf("createCredit.createCredit > %w", err)
		}

		return nil
	})

	if err != nil {
		return CreditResult{}, fmt.Errorf("txManager.WithTransaction > %w", err)
	}

	return CreditResult{
		BankID:     createdCredit.BankID().String(),
		ClientID:   createdCredit.ClientID().String(),
		CreatedAt:  createdCredit.CreatedAt(),
		CreditType: createdCredit.CreditType().String(),
		ID:         createdCredit.ID().String(),
		MaxPayment: createdCredit.MaxPayment().String(),
		MinPayment: createdCredit.MinPayment().String(),
		TermMonths: createdCredit.TermMonths(),
		Status:     createdCredit.Status().String(),
	}, nil
}

type CreateCreditCommand struct {
	BankID     string
	ClientID   string
	CreditType string
}

func (cmd CreateCreditCommand) validate() error {
	err := mye.New(mye.CodeInvalid, "credit_creation_failed", "validation error").
		WithUserMsg("credit creation failed due to input validation")

	if len(cmd.BankID) == 0 {
		return err.WithField("bank_id", "bank_id can't be empty")
	}

	if len(cmd.ClientID) == 0 {
		return err.WithField("client_id", "client_id can't be empty")
	}

	if len(cmd.CreditType) == 0 {
		return err.WithField("credit_type", "credit_type can't be empty")
	}

	if err.HasFields() {
		return err
	}

	return nil
}

func newCreditOpts(cmd CreateCreditCommand) (credit.NewCreditOpts, error) {
	bankId, err := bank.NewID(cmd.BankID)
	if err != nil {
		return credit.NewCreditOpts{}, fmt.Errorf("CreateCredit > %w", err)
	}

	clientID, err := client.NewID(cmd.ClientID)
	if err != nil {
		return credit.NewCreditOpts{}, fmt.Errorf("CreateCredit > %w", err)
	}

	creditType, err := credit.CreditTypeFromString(cmd.CreditType)
	if err != nil {
		return credit.NewCreditOpts{}, fmt.Errorf("CreateCredit > %w", err)
	}

	return credit.NewCreditOpts{
		BankID:     bankId,
		ClientID:   clientID,
		CreditType: creditType,
	}, nil
}

func (svc createCredit) getBankAndClient(
	ctx context.Context,
	clientID client.ID,
	bankID bank.ID,
) (bank.Bank, client.Client, error) {
	clientCh := make(chan client.Client, 2)
	bankCh := make(chan bank.Bank, 2)
	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		c, err := svc.clientRepo.GetClient(ctx, clientID)
		if err != nil {
			errCh <- err
			return
		}

		clientCh <- c
	}()

	go func() {
		defer wg.Done()

		b, err := svc.bankRepo.GetBank(ctx, bankID)
		if err != nil {
			errCh <- err
			return
		}

		bankCh <- b
	}()

	wg.Wait()

	if len(errCh) > 0 {
		return bank.Bank{}, client.Client{}, fmt.Errorf("verifyEntities > %w", <-errCh)
	}

	return <-bankCh, <-clientCh, nil
}

func (svc createCredit) createCreditAndPublishEvent(
	ctx context.Context,
	newCredit credit.Credit,
) (credit.Credit, error) {
	c, err := svc.creditRepo.CreateCredit(ctx, newCredit)
	if err != nil {
		return credit.Credit{}, fmt.Errorf("creditRepo.CreateCredit > %w", err)
	}

	creditCreated := event.CreditCreated{
		BankID:   c.BankID().String(),
		ClientID: c.ClientID().String(),
		CreditID: c.ID().String(),
	}

	if err := svc.eventBus.Publish(ctx, creditCreated); err != nil {
		return credit.Credit{}, fmt.Errorf("eventBus.Publish > %w", err)
	}

	return c, nil
}

type CreditResult struct {
	BankID     string
	ClientID   string
	CreatedAt  time.Time
	CreditType string
	ID         string
	MaxPayment string
	MinPayment string
	Status     string
	TermMonths int
}

//nolint:errcheck
func validateCreateService(
	bankRepo bankRepository,
	clientRepo clientRepository,
	creditRepo creditRepository,
	eventBus eventBus,
	txManager transactionManager,
) error {
	err := mye.New(mye.CodeInternal, "createCredit_service_config_error", "createService parameter validation failed")

	if bankRepo == nil {
		err.WithField("bankRepo", "bank repository is missing")

	}

	if clientRepo == nil {
		err.WithField("clientRepo", "client repository is missing")
	}

	if creditRepo == nil {
		err.WithField("creditRepo", "credit repository is missing")
	}

	if eventBus == nil {
		err.WithField("eventBus", "event bus is missing")
	}

	if txManager == nil {
		err.WithField("txManager", "transaction manager is missing")
	}

	if err.HasFields() {
		return err
	}

	return nil
}
