package clients

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/pkg/pointer"
	"github.com/edgarSucre/mye"
)

type createClient struct {
	repo Repository
}

func NewCreateClientService(repo Repository) (CreateClientService, error) {
	if repo == nil {
		return nil, mye.New(
			mye.CodeInvalid,
			"createClient_service_config_error",
			"client repository is missing",
		)
	}
	return createClient{repo: repo}, nil
}

func (svc createClient) CreateClient(
	ctx context.Context,
	cmd CreateClientCommand,
) (ClientResult, error) {
	if err := cmd.Validate(); err != nil {
		return ClientResult{}, err
	}

	birthdate, err := client.NewBirthdate(cmd.Birthdate)
	if err != nil {
		return ClientResult{}, err
	}

	email, err := client.NewEmail(cmd.Email)
	if err != nil {
		return ClientResult{}, err
	}

	clientParam, err := client.New(birthdate, cmd.Country, email, cmd.FullName)
	if err != nil {
		return ClientResult{}, err
	}

	newClient, err := svc.repo.CreateClient(ctx, clientParam)
	if err != nil {
		return ClientResult{}, err
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
	err := mye.New(mye.CodeInvalid, "client_creation_failed", "validation failed").
		WithUserMsg("client creation failed due to validation")

	if cmd.Birthdate != nil && cmd.Birthdate.IsZero() {
		err.WithField("birthdate", "birthdate is not a valid date")
	}

	if _, mailErr := mail.ParseAddress(cmd.Email); mailErr != nil {
		err.WithField("email", "email is not a valid email address")
	}

	if len(strings.TrimSpace(cmd.FullName)) == 0 {
		return err.WithField("full_name", "full_name can't be empty")
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
