package domain

import (
	"context"
	"time"

	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type (
	CreditType   string
	CreditStatus string
)

const (
	CreditTypeAuto       CreditType = "auto"
	CreditTypeMortgage   CreditType = "mortgage"
	CreditTypeCommercial CreditType = "commercial"

	CreditStatusPending  CreditStatus = "pending"
	CreditStatusApproved CreditStatus = "approved"
	CreditStatusRejected CreditStatus = "rejected"
)

func (ct CreditType) isValid() bool {
	return ct == CreditTypeAuto || ct == CreditTypeCommercial || ct == CreditTypeMortgage
}

type (
	Credit struct {
		BankID     uuid.UUID
		ClientID   uuid.UUID
		CreatedAt  time.Time
		CreditType CreditType
		ID         uuid.UUID
		MaxPayment decimal.Decimal
		MinPayment decimal.Decimal
		Status     CreditStatus
		TermMonths int
	}

	CreditCreated struct {
		BankID   uuid.UUID
		ClientID uuid.UUID
		CreditID uuid.UUID
	}

	CreditApproved struct {
		CreditID uuid.UUID
	}

	CreditRejected struct {
		CreditID uuid.UUID
	}
)

func (CreditCreated) EventName() string {
	return "credit.created"
}

func (CreditApproved) EventName() string {
	return "credit.approved"
}

func (CreditRejected) EventName() string {
	return "credit.rejected"
}

type CreateCreditParams struct {
	BankID     uuid.UUID
	ClientID   uuid.UUID
	CreditType CreditType
}

var (
	ErrCreditTypeInvalid = terror.Validation.New(
		"invalid-credit-type",
		"credit_type must one of 'auto', 'mortgage', 'commercial'",
	)
)

func (params CreateCreditParams) Validate() error {
	if !params.CreditType.isValid() {
		return ErrCreditTypeInvalid
	}

	return nil
}

type CreditService interface {
	CreateCredit(context.Context, CreateCreditParams) error
	GetCredit(context.Context, uuid.UUID) (*Credit, error)
	ProcessCredit(context.Context, CreditCreated) error

	// these are just for logging, not implemented
	ProcessApproval(context.Context, CreditApproved) error
	ProcessRejection(context.Context, CreditRejected) error
}
