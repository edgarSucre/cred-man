package credit

import (
	"fmt"
	"time"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/shopspring/decimal"
)

type Credit struct {
	bankID     bank.ID
	clientID   client.ID
	createdAt  time.Time
	creditType CreditType
	id         ID
	maxPayment decimal.Decimal
	minPayment decimal.Decimal
	status     CreditStatus
	termMonths int
}

type NewCreditOpts struct {
	BankID     bank.ID
	ClientID   client.ID
	CreditType CreditType
	MaxPayment decimal.Decimal
	MinPayment decimal.Decimal
	TermMonths int
}

func (opts NewCreditOpts) validate() error {
	if opts.BankID.IsEmpty() {
		return fmt.Errorf("NewCreditOpts.validate > %w", bank.ErrInvalidBankID)
	}

	if opts.ClientID.IsEmpty() {
		return fmt.Errorf("NewCreditOpts.validate > %w", client.ErrInvalidClientID)
	}

	if opts.CreditType.IsInvalid() {
		return fmt.Errorf("NewCreditOpts.validate > %w", ErrInvalidCreditType)
	}

	if opts.MaxPayment.IsZero() {
		return fmt.Errorf("NewCreditOpts.validate > %w", ErrInvalidMaxPayment)
	}

	if opts.MinPayment.IsZero() {
		return fmt.Errorf("NewCreditOpts.validate > %w", ErrInvalidMinPayment)
	}

	if opts.MinPayment.GreaterThan(opts.MaxPayment) {
		return fmt.Errorf("NewCreditOpts.validate > %w", ErrInvalidPayment)
	}

	if opts.TermMonths <= 0 {
		return ErrInvalidTerm
	}

	return nil
}

func New(opts NewCreditOpts) (Credit, error) {
	if err := opts.validate(); err != nil {
		return Credit{}, fmt.Errorf("credit.New > %w", err)
	}

	return Credit{
		bankID:     opts.BankID,
		clientID:   opts.ClientID,
		creditType: opts.CreditType,
		maxPayment: opts.MaxPayment,
		minPayment: opts.MinPayment,
		status:     CreditStatusPending,
		termMonths: opts.TermMonths,
	}, nil
}

type RehydrateOpts struct {
	BankID     bank.ID
	ClientID   client.ID
	CreditType CreditType
	CreatedAt  time.Time
	ID         ID
	MaxPayment decimal.Decimal
	MinPayment decimal.Decimal
	Status     CreditStatus
	TermMonths int
}

func Rehydrate(opts RehydrateOpts) Credit {
	return Credit{
		bankID:     opts.BankID,
		clientID:   opts.ClientID,
		createdAt:  opts.CreatedAt,
		creditType: opts.CreditType,
		id:         opts.ID,
		maxPayment: opts.MaxPayment,
		minPayment: opts.MinPayment,
		status:     opts.Status,
		termMonths: opts.TermMonths,
	}
}

func (credit Credit) BankID() bank.ID {
	return credit.bankID
}

func (credit Credit) ClientID() client.ID {
	return credit.clientID
}

func (credit Credit) CreatedAt() time.Time {
	return credit.createdAt
}

func (credit Credit) CreditType() CreditType {
	return credit.creditType
}

func (credit Credit) ID() ID {
	return credit.id
}

func (credit Credit) MaxPayment() decimal.Decimal {
	return credit.maxPayment
}

func (credit Credit) MinPayment() decimal.Decimal {
	return credit.minPayment
}

func (credit Credit) Status() CreditStatus {
	return credit.status
}

func (credit Credit) TermMonths() int {
	return credit.termMonths
}

func (credit Credit) Approve() {
	credit.status = CreditStatusApproved
}

func (credit Credit) IsEqual(c Credit) bool {
	return credit.id.IsEqual(c.id)
}
