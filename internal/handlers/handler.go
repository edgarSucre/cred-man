package handlers

import (
	"github.com/edgarSucre/crm/internal/infrastructure/events"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/terror"
)

var ErrNoCreditService = terror.Internal.New("event-handler-bad-config", "credit service is missing")

func GetCreditHandlers(svc domain.CreditService) ([]events.EventHandler, error) {
	if svc == nil {
		return nil, ErrNoCreditService
	}

	creditCreated := &CreditCreatedHandler{creditSvc: svc}
	creditApproved := &CreditApprovedHandler{creditSvc: svc}
	creditRejected := &CreditRejectedHandler{creditSvc: svc}

	return []events.EventHandler{
		creditCreated,
		creditApproved,
		creditRejected,
	}, nil

}
