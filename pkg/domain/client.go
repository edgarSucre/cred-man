package domain

import (
	"context"
	"strings"
	"time"

	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID
	Birthdate time.Time
	Country   string
	Email     string
	FullName  string
	CreatedAt time.Time
}

type CreateClientParams struct {
	Birthdate *time.Time
	Country   *string
	Email     string
	FullName  string
}

var (
	ErrFullNameRequired = terror.Validation.New("invalid-name", "client full_name is required")
	ErrEmailRequired    = terror.Validation.New("invalid-email", "client email is required")
)

func (params CreateClientParams) Validate() error {
	if len(strings.TrimSpace(params.FullName)) == 0 {
		return ErrFullNameRequired
	}

	if len(strings.TrimSpace(params.Email)) == 0 {
		return ErrEmailRequired
	}

	return nil
}

type ClientService interface {
	CreateClient(context.Context, CreateClientParams) (*Client, error)
	GetClient(context.Context, uuid.UUID) (*Client, error)
}
