package bank

import (
	"fmt"
	"net/http"

	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/google/uuid"
)

func HandleCreateBank(svc domain.BankService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateBankRequest

		if err := httputils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		newBank, err := svc.CreateBank(r.Context(), req.ToParams())
		if err != nil {
			return fmt.Errorf("svc.CreateBank: %w", err)
		}

		resp := new(CreateBankResponse)
		resp.FromDomain(newBank)

		w.WriteHeader(http.StatusCreated)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

type CreateBankRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (req CreateBankRequest) ToParams() domain.CreateBankParams {
	return domain.CreateBankParams{
		Name: req.Name,
		Type: domain.BankType(req.Type),
	}
}

type CreateBankResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

func (resp *CreateBankResponse) FromDomain(dm *domain.Bank) {
	resp.ID = dm.ID
	resp.Name = dm.Name
	resp.Type = string(dm.Type)
}

//go:generate easyjson -all -snake_case $GOFILE
