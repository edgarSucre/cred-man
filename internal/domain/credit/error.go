package credit

import "github.com/edgarSucre/crm/pkg/terror"

var (
	ErrInvalidCreditID     = terror.Validation.New("invalid-credit-id", "id most be a valid uuid v4")
	ErrInvalidMaxPayment   = terror.Validation.New("invalid-max-payment", "max_payment can't be zero")
	ErrInvalidMinPayment   = terror.Validation.New("invalid-min-payment", "min_payment can't be zero")
	ErrInvalidPayment      = terror.Validation.New("invalid-payment", "min_payment can't be greater than max_payment")
	ErrInvalidCreditType   = terror.Validation.New("invalid-credit-type", "credit_type most be 'auto', 'mortgage' or 'commercial'")
	ErrInvalidCreditStatus = terror.Validation.New("invalid-credit-type", "credit_status most be 'pending', 'approved' or 'rejected'")
	ErrInvalidTerm         = terror.Validation.New("invalid-term-months", "term_months most be greater than zero")
	ErrNoCreatedAt         = terror.Internal.New("empty-created-at", "credit was created with an empty date")
)
