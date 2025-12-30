package http

import (
	"fmt"
	"net/http"

	"github.com/edgarSucre/crm/internal/application/banks"
	"github.com/edgarSucre/crm/internal/infrastructure/thttp/httputils"
)

type BankHandler struct {
	createBank banks.CreateBankService
}

func NewBankHandler(createBank banks.CreateBankService) BankHandler {
	return BankHandler{createBank}
}

func HandleCreateBank(svc banks.CreateBankService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateBankRequest

		if err := httputils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		newBank, err := svc.Execute(r.Context(), req.ToParams())
		if err != nil {
			return fmt.Errorf("svc.CreateBank: %w", err)
		}

		resp := new(CreateBankResponse)
		resp.FromResult(newBank)

		w.WriteHeader(http.StatusCreated)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

type CreateBankRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (req CreateBankRequest) ToParams() banks.CreateBankCmd {
	return banks.CreateBankCmd{
		Name: req.Name,
		Type: req.Type,
	}
}

type CreateBankResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (resp *CreateBankResponse) FromResult(bankResult banks.CreateBankResult) {
	resp.ID = bankResult.ID
	resp.Name = bankResult.Name
	resp.Type = bankResult.Type
}

//go:generate easyjson -all -snake_case $GOFILE
