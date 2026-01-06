package handlers

import (
	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/infrastructure/events"
)

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
