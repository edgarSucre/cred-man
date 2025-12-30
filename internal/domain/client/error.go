package client

import "github.com/edgarSucre/crm/pkg/terror"

var (
	ErrInvalidClientID       = terror.Validation.New("invalid-client-id", "id should be a valid uuid v4")
	ErrInvalidClientFullName = terror.Validation.New("invalid-client-full-name", "client full name can't be empty")
	ErrInvalidBirthdate      = terror.Validation.New("invalid-client-birthday", "client birthdate most be a valid date")
	ErrInvalidEmail          = terror.Validation.New("invalid-client-email", "client email should be a valid email address")
)
