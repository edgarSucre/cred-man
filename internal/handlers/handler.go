package handlers

import (
	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/infrastructure/events"
	"github.com/edgarSucre/crm/pkg/terror"
)

var ErrNoCreditService = terror.Internal.New("event-handler-bad-config", "credit service is missing")

func GetCreditHandlers(
	approveCredit credits.ApproveCreditService,
	processCredit credits.ProcessCreditService,
	rejectCredit credits.RejectCreditService,
) []events.EventHandler {

	creditCreated := &CreditCreatedHandler{processCredit}
	creditApproved := &CreditApprovedHandler{approveCredit}
	creditRejected := &CreditRejectedHandler{rejectCredit}

	return []events.EventHandler{
		creditCreated,
		creditApproved,
		creditRejected,
	}
}
