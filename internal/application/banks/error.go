package banks

import "github.com/edgarSucre/crm/pkg/terror"

var (
	ErrBankNameRequired = terror.Validation.New("invalid-name", "bank name is required")
	ErrBankTypeInvalid  = terror.Validation.New("invalid-bank-type", "bank type must be 'private' or 'government'")
)
