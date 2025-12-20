package credit

import (
	"net/http"

	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
)

func HandleCreateCredit(svc domain.CreditService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateCreditRequest

		if err := httputils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		params, err := req.ToParams()
		if err != nil {
			return err
		}

		if err := svc.CreateCredit(r.Context(), params); err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		return nil
	}

	return httputils.ErrorHandlerFunc(fn)
}

type CreateCreditRequest struct {
	BankID     string `json:"bank_id"`
	ClientID   string `json:"client_id"`
	CreditType string `json:"credit_type"`
}

var (
	ErrInvalidBankID   = terror.Validation.New("invalid-bank-id", "bank_id is not a valid uuid v4")
	ErrInvalidClientID = terror.Validation.New("invalid-client-id", "client_id is not a valid uuid v4")
)

func (req CreateCreditRequest) ToParams() (domain.CreateCreditParams, error) {
	var params domain.CreateCreditParams

	bankID, err := uuid.Parse(req.BankID)
	if err != nil {
		return params, ErrInvalidBankID
	}

	clientID, err := uuid.Parse(req.ClientID)
	if err != nil {
		return params, ErrInvalidClientID
	}

	params.BankID = bankID
	params.ClientID = clientID
	params.CreditType = domain.CreditType(req.CreditType)

	return params, nil
}

//go:generate easyjson -all -snake_case $GOFILE
