package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/edgarSucre/crm/internal/application/clients"
	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/mye"
)

type ClientHandler struct {
	createClient clients.CreateClientService
	getClient    clients.GetClientService
}

//nolint:errcheck
func NewClientHandler(
	createClient clients.CreateClientService,
	getClient clients.GetClientService,
) (ClientHandler, error) {
	err := mye.New(
		mye.CodeInternal,
		"client_handler_config_error",
		"client handler validation error",
	)

	if createClient == nil {
		err.WithField("createClient", "createClient service is missing")
	}

	if getClient == nil {
		err.WithField("getClient", "getClient service is missing")
	}

	if err.HasFields() {
		return ClientHandler{}, err
	}

	return ClientHandler{createClient, getClient}, nil
}

/* ============================================================================================== */
/*                                       HandleCreateClient                                       */
/* ============================================================================================== */

func HandleCreateClient(svc clients.CreateClientService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateClientRequest

		if err := httputils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		clientResult, err := svc.CreateClient(r.Context(), req.ToCommand())
		if err != nil {
			return fmt.Errorf("svc.CreateClient: %w", err)
		}

		resp := new(ClientResponse)
		resp.FromResult(clientResult)

		w.WriteHeader(http.StatusCreated)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

//easyjson:json
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

func (req CreateClientRequest) ToCommand() clients.CreateClientCommand {
	var date *time.Time
	if req.BirthDate != nil {
		date = &req.BirthDate.Value
	}

	return clients.CreateClientCommand{
		Birthdate: date,
		Country:   req.Country,
		Email:     req.Email,
		FullName:  req.FullName(),
	}
}

//easyjson:json
type ClientResponse struct {
	BirthDate string `json:"birthdate"`
	Country   string `json:"country"`
	CreatedAt string `json:"created_at"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	ID        string `json:"id"`
}

func (resp *ClientResponse) FromResult(clientResult clients.ClientResult) {
	resp.BirthDate = clientResult.Birthdate
	resp.Country = clientResult.Country
	resp.CreatedAt = clientResult.CreatedAt
	resp.Email = clientResult.Email
	resp.FullName = clientResult.FullName
	resp.ID = clientResult.ID
}

/* ============================================================================================== */
/*                                         HandleGetClient                                        */
/* ============================================================================================== */
func HandleGetCLient(svc clients.GetClientService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		clientId := r.PathValue("id")

		clientResult, err := svc.GetClient(r.Context(), clients.GetClientCommand{ID: clientId})
		if err != nil {
			return err
		}

		response := new(ClientResponse)
		response.FromResult(clientResult)

		return httputils.Marshal(w, response)
	}

	return httputils.ErrorHandlerFunc(fn)
}

//go:generate easyjson -snake_case $GOFILE
