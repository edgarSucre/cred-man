package clients

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/pkg/pointer"
)

type createClient struct {
	repo Repository
}

func NewCreateClientService(repo Repository) CreateClientService {
	return createClient{repo: repo}
}

func (svc createClient) CreateClient(
	ctx context.Context,
	cmd CreateClientCommand,
) (ClientResult, error) {
	if err := cmd.Validate(); err != nil {
		return ClientResult{}, fmt.Errorf("createClient.CommandValidation: %w", err)
	}

	birthdate, err := client.NewBirthdate(cmd.Birthdate)
	if err != nil {
		return ClientResult{}, fmt.Errorf("client.NewBirthdate > %w", err)
	}

	email, err := client.NewEmail(cmd.Email)
	if err != nil {
		return ClientResult{}, fmt.Errorf("client.NewEmail > %w", err)
	}

	clientParam, err := client.New(birthdate, cmd.Country, email, cmd.FullName)
	if err != nil {
		return ClientResult{}, fmt.Errorf("createCLient.client.New: %w", err)
	}

	newClient, err := svc.repo.CreateClient(ctx, clientParam)
	if err != nil {
		return ClientResult{}, fmt.Errorf("createClient.repo.CreateClient: %w", err)
	}

	return mapClientResult(newClient), nil
}

type CreateClientCommand struct {
	Birthdate *time.Time
	Country   *string
	Email     string
	FullName  string
}

func (cmd CreateClientCommand) Validate() error {
	if cmd.Birthdate != nil && cmd.Birthdate.IsZero() {
		return ErrInvalidBirthDate
	}

	if _, err := mail.ParseAddress(cmd.Email); err != nil {
		return ErrInvalidEmail
	}

	if len(strings.TrimSpace(cmd.FullName)) == 0 {
		return ErrNoFullName
	}

	return nil
}

type ClientResult struct {
	ID        string
	Birthdate string
	Country   string
	Email     string
	FullName  string
	CreatedAt string
}

func mapClientResult(cl client.Client) ClientResult {
	return ClientResult{
		ID:        cl.ID().String(),
		Birthdate: cl.Birthdate().String(),
		Country:   pointer.ValueOrEmpty(cl.Country()),
		Email:     cl.Email().String(),
		FullName:  cl.FullName(),
		CreatedAt: cl.CreatedAt().Format(time.DateTime),
	}
}
