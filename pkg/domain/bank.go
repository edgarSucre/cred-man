package domain

import (
	"context"
	"strings"

	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
)

type BankType string

const (
	BankTypePrivate    BankType = "private"
	BankTypeGovernment BankType = "government"
)

func (bt BankType) isValid() bool {
	return bt == BankTypeGovernment || bt == BankTypePrivate
}

type Bank struct {
	ID   uuid.UUID
	Name string
	Type BankType
}

type CreateBankParams struct {
	Name string
	Type BankType
}

var (
	ErrBankNameRequired = terror.Validation.New("invalid-name", "bank name is required")
	ErrBankTypeInvalid  = terror.Validation.New("invalid-bank-type", "bank type must be 'private' or 'government'")
)

func (params CreateBankParams) Validate() error {
	if len(strings.TrimSpace(params.Name)) == 0 {
		return ErrBankNameRequired
	}

	if !params.Type.isValid() {
		return ErrBankTypeInvalid
	}

	return nil
}

type BankService interface {
	CreateBank(context.Context, CreateBankParams) (*Bank, error)
}
