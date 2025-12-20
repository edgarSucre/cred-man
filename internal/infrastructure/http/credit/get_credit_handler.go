package credit

import (
	"net/http"
	"time"

	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/google/uuid"
)

var ErrInvalidCreditID = terror.Validation.New("invalid-credit-id", "id is not a valid uuid v4")

func HandleGetCredit(svc domain.CreditService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		creditID := r.PathValue("id")

		id, err := uuid.Parse(creditID)
		if err != nil {
			return ErrInvalidCreditID
		}

		credit, err := svc.GetCredit(r.Context(), id)
		if err != nil {
			return err
		}

		resp := new(GetCreditResponse)
		resp.FromDomain(credit)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}

type GetCreditResponse struct {
	BankID     uuid.UUID
	ClientID   uuid.UUID
	CreatedAt  time.Time
	CreditType string
	ID         uuid.UUID
	MaxPayment string
	MinPayment string
	Status     string
	TermMonths int
}

func (resp *GetCreditResponse) FromDomain(dm *domain.Credit) {
	resp.BankID = dm.BankID
	resp.ClientID = dm.ClientID
	resp.CreatedAt = dm.CreatedAt
	resp.CreditType = string(dm.CreditType)
	resp.ID = dm.ID
	resp.MaxPayment = dm.MaxPayment.String()
	resp.MinPayment = dm.MinPayment.String()
	resp.Status = string(dm.Status)
	resp.TermMonths = dm.TermMonths
}

//go:generate easyjson -all -snake_case $GOFILE
