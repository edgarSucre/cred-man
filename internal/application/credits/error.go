package credits

import "github.com/edgarSucre/crm/pkg/terror"

var (
	ErrNoCreditID        = terror.Validation.New("no_credit_id", "credit_id is required")
	ErrNoBankID          = terror.Validation.New("no_bank_id", "bank_id is required")
	ErrNoClientID        = terror.Validation.New("no_client_id", "client_id is required")
	ErrInvalidMaxPayment = terror.Validation.New("invalid_max_payment", "max_payment can't be zero")
	ErrInvalidMinPayment = terror.Validation.New("invalid_min_payment", "min_payment can't be zero")
	ErrNoCreditType      = terror.Validation.New("no_credit_type", "credit_type is required")
	ErrNoCreditStatus    = terror.Validation.New("no_credit_type", "credit_status is required")
	ErrInvalidTerm       = terror.Validation.New("invalid_term_months", "term_months most be greater than zero")
)
