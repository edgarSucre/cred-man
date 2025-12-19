package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edgarSucre/crm/internal/client"
	"github.com/edgarSucre/crm/internal/http/utils"
	"github.com/edgarSucre/crm/pkg"
	"github.com/google/uuid"
)

func HandleCreateClient(svc *client.Service) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateClientRequest

		if err := utils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		newClient, err := svc.CreateClient(r.Context(), req.ToParams())
		if err != nil {
			return fmt.Errorf("svc.CreateClient: %w", err)
		}

		resp := new(CreateClientResponse)
		resp.FromDomain(newClient)

		return utils.Marshal(w, resp)
	}

	return utils.ErrorHandlerFunc(fn)
}

type CreateClientRequest struct {
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

func (req CreateClientRequest) FullName() string {
	return fmt.Sprintf("%s %s", req.FirstName, req.LastName)
}

func (req CreateClientRequest) ToParams() pkg.CreateClientParams {
	return pkg.CreateClientParams{
		BirthDate: req.BirthDate,
		Country:   req.Country,
		Email:     req.Email,
		FullName:  req.FullName(),
	}
}

type CreateClientResponse struct {
	BirthDate time.Time `json:"birth_date"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	FullName  string    `json:"ful_name"`
	ID        uuid.UUID `json:"id"`
}

func (resp *CreateClientResponse) FromDomain(dto *pkg.Client) {
	resp.BirthDate = dto.BirthDate
	resp.Country = dto.Country
	resp.CreatedAt = dto.CreatedAt
	resp.Email = dto.Email
	resp.FullName = dto.FullName
	resp.ID = dto.ID
}

//go:generate easyjson -all -snake_case $GOFILE
