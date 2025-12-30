package clients

import "github.com/edgarSucre/crm/pkg/terror"

var (
	ErrInvalidBirthDate = terror.Validation.New("invalid_birthdate", "birthdate is not a valid date")
	ErrInvalidEmail     = terror.Validation.New("invalid_email", "email is not a valid email-address")
	ErrNoFullName       = terror.Validation.New("invalid_name", "name can't be empty")
)
