package pkg

import (
	"time"

	"github.com/edgarSucre/crm/internal/db/repository"
	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID
	BirthDate time.Time
	Country   string
	Email     string
	FullName  string
	CreatedAt time.Time
}

func (cl *Client) FromModel(rClient repository.Client) {
	cl.BirthDate = rClient.Birthdate
	cl.Country = rClient.Country
	cl.CreatedAt = rClient.CreatedAt
	cl.Email = rClient.Email
	cl.FullName = rClient.FullName
	cl.ID = rClient.ID
}

type CreateClientParams struct {
	BirthDate time.Time
	Country   string
	Email     string
	FullName  string
}

func (params CreateClientParams) ToModel() repository.CreateClientParams {
	return repository.CreateClientParams{
		Birthdate: params.BirthDate,
		Country:   params.Country,
		Email:     params.Email,
		FullName:  params.FullName,
	}
}

func (params CreateClientParams) Validate() error {
	//TODO handle error here
	return nil
}

type ClientService interface {
}
