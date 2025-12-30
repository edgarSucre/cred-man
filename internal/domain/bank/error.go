package bank

import "github.com/edgarSucre/crm/pkg/terror"

var (
	ErrInvalidBankID   = terror.Validation.New("invalid-bank-id", "id should be a valid uuid v4")
	ErrInvalidBankName = terror.Validation.New("invalid-name", "bank name can't be empty")
	ErrBankTypeInvalid = terror.Validation.New("invalid-bank-type", "bank type must be 'private' or 'government'")
)
