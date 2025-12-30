package credits

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/shopspring/decimal"
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
) CreateCreditService {
	return createCredit{bankRepo, clientRepo, creditRepo, eventBus, txManager}
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
	MaxPayment float64
	MinPayment float64
	TermMonths int
}

func (cmd CreateCreditCommand) validate() error {
	if len(cmd.BankID) == 0 {
		return fmt.Errorf("CreateCreditCommand.validate > %w", ErrNoBankID)
	}

	if len(cmd.ClientID) == 0 {
		return fmt.Errorf("CreateCreditCommand.validate > %w", ErrNoClientID)
	}

	if len(cmd.CreditType) == 0 {
		return fmt.Errorf("CreateCreditCommand.validate > %w", ErrNoCreditType)
	}

	if cmd.MaxPayment <= 0 {
		return fmt.Errorf("CreateCreditCommand.validate > %w", ErrInvalidMaxPayment)
	}

	if cmd.MinPayment <= 0 {
		return fmt.Errorf("CreateCreditCommand.validate > %w", ErrInvalidMinPayment)
	}

	if cmd.TermMonths <= 0 {
		return fmt.Errorf("CreateCreditCommand.validate > %w", ErrInvalidTerm)
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

	maxPayment := decimal.NewFromFloat(cmd.MaxPayment)
	minPayment := decimal.NewFromFloat(cmd.MinPayment)

	return credit.NewCreditOpts{
		BankID:     bankId,
		ClientID:   clientID,
		CreditType: creditType,
		MaxPayment: maxPayment,
		MinPayment: minPayment,
		TermMonths: cmd.TermMonths,
	}, nil
}

func (svc createCredit) getBankAndClient(
	ctx context.Context,
	clientID client.ID,
	bankID bank.ID,
) (bank.Bank, client.Client, error) {
	clientCh := make(chan client.Client)
	bankCh := make(chan bank.Bank)
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
		return
	}()

	go func() {
		defer wg.Done()

		b, err := svc.bankRepo.GetBank(ctx, bankID)
		if err != nil {
			errCh <- err
			return
		}

		bankCh <- b
		return
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

	if err := svc.eventBus.Publish(ctx, credit.CreditCreated{ID: c.ID().String()}); err != nil {
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
