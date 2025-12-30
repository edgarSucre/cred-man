package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/domain/event"
)

type CreditApprovedHandler struct {
	svc credits.ApproveCreditService
}

func (h *CreditApprovedHandler) EventName() string {
	return "credit.approved"
}

func (h *CreditApprovedHandler) Handle(
	ctx context.Context,
	payload json.RawMessage,
) error {
	var approveEvent event.CreditApproved
	if err := json.Unmarshal(payload, &approveEvent); err != nil {
		return err
	}

	h.svc.Execute(ctx, approveEvent.CreditID)

	return nil
}
