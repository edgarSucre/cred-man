package http

import (
	"net/http"
	"time"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/mye"
)

type CreditHandler struct {
	createCredit credits.CreateCreditService
	getCredit    credits.GetCreditService
}

func NewCreditHandler(
	createCredit credits.CreateCreditService,
	getCredit credits.GetCreditService,
) (CreditHandler, error) {
	err := mye.New(
		mye.CodeInternal,
		"credit_handler_config_error",
		"credit handler validation error",
	)

	if createCredit == nil {
		err.WithField("createCredit", "createCredit service is missing")
	}

	if getCredit == nil {
		err.WithField("getCredit", "getCredit service is missing")
	}

	if err.HasFields() {
		return CreditHandler{}, nil
	}
	return CreditHandler{createCredit, getCredit}, nil
}

/* ============================================================================================== */
/*                                       HandleCreateCredit                                       */
/* ============================================================================================== */

func HandleCreateCredit(svc credits.CreateCreditService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		var req CreateCreditRequest

		if err := httputils.Unmarshal(r.Body, &req); err != nil {
			return err
		}

		creditResult, err := svc.Execute(r.Context(), req.ToCommand())
		if err != nil {
			return err
		}

		resp := new(CreditResponse)
		resp.FromResult(creditResult)

		w.WriteHeader(http.StatusCreated)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

//easyjson:json
type CreateCreditRequest struct {
	BankID     string `json:"bank_id"`
	ClientID   string `json:"client_id"`
	CreditType string `json:"credit_type"`
}

func (req CreateCreditRequest) ToCommand() credits.CreateCreditCommand {
	return credits.CreateCreditCommand{
		BankID:     req.BankID,
		ClientID:   req.ClientID,
		CreditType: req.CreditType,
	}
}

//easyjson:json
type CreditResponse struct {
	BankID     string `json:"bank_id"`
	ClientID   string `json:"client_id"`
	CreditType string `json:"credit_type"`
	CreatedAt  string `json:"created_at"`
	ID         string `json:"id"`
	MaxPayment string `json:"max_payment"`
	MinPayment string `json:"min_payment"`
	Status     string `json:"status"`
	TermMonths int    `json:"term_months"`
}

func (resp *CreditResponse) FromResult(creditResult credits.CreditResult) {
	resp.BankID = creditResult.BankID
	resp.ClientID = creditResult.ClientID
	resp.CreatedAt = creditResult.CreatedAt.Format(time.DateOnly)
	resp.ID = creditResult.ID
	resp.MaxPayment = creditResult.MaxPayment
	resp.MinPayment = creditResult.MinPayment
	resp.Status = creditResult.Status
	resp.TermMonths = creditResult.TermMonths
}

/* ============================================================================================== */
/*                                         HandleGetCredit                                        */
/* ============================================================================================== */

func HandleGetCredit(svc credits.GetCreditService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		creditID := r.PathValue("id")

		creditResult, err := svc.Execute(r.Context(), credits.GetCreditCommand{ID: creditID})
		if err != nil {
			return err
		}

		resp := new(CreditResponse)
		resp.FromResult(creditResult)

		w.WriteHeader(http.StatusOK)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

//go:generate easyjson -snake_case $GOFILE
