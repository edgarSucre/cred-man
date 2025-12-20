package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edgarSucre/crm/internal/http/httputils"
	"github.com/edgarSucre/crm/pkg"
	"github.com/google/uuid"
)

func HandleCreateClient(svc pkg.ClientService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateClientRequest

		if err := httputils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		newClient, err := svc.CreateClient(r.Context(), req.ToParams())
		if err != nil {
			return fmt.Errorf("svc.CreateClient: %w", err)
		}

		resp := new(CreateClientResponse)
		resp.FromDomain(newClient)

		w.WriteHeader(http.StatusCreated)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

type CreateClientRequest struct {
	BirthDate *httputils.Date `json:"birthdate"`
	Country   *string         `json:"country"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
}

func (req CreateClientRequest) FullName() string {
	return fmt.Sprintf("%s %s", req.FirstName, req.LastName)
}

func (req CreateClientRequest) ToParams() pkg.CreateClientParams {
	var date *time.Time
	if req.BirthDate != nil {
		date = &req.BirthDate.Time
	}

	return pkg.CreateClientParams{
		Birthdate: date,
		Country:   req.Country,
		Email:     req.Email,
		FullName:  req.FullName(),
	}
}

type CreateClientResponse struct {
	BirthDate time.Time `json:"birthdate"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	ID        uuid.UUID `json:"id"`
}

func (resp *CreateClientResponse) FromDomain(dto *pkg.Client) {
	resp.BirthDate = dto.Birthdate
	resp.Country = dto.Country
	resp.CreatedAt = dto.CreatedAt
	resp.Email = dto.Email
	resp.FullName = dto.FullName
	resp.ID = dto.ID
}

//go:generate easyjson -all -snake_case $GOFILE
