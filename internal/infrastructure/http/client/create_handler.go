package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/google/uuid"
)

func HandleCreateClient(svc domain.ClientService) http.Handler {
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

func (req CreateClientRequest) ToParams() domain.CreateClientParams {
	var date *time.Time
	if req.BirthDate != nil {
		date = &req.BirthDate.Time
	}

	return domain.CreateClientParams{
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

func (resp *CreateClientResponse) FromDomain(dm *domain.Client) {
	resp.BirthDate = dm.Birthdate
	resp.Country = dm.Country
	resp.CreatedAt = dm.CreatedAt
	resp.Email = dm.Email
	resp.FullName = dm.FullName
	resp.ID = dm.ID
}

//go:generate easyjson -all -snake_case $GOFILE
