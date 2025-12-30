package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/domain/event"
)

type CreditCreatedHandler struct {
	svc credits.ProcessCreditService
}

func (h *CreditCreatedHandler) EventName() string {
	return "credit.created"
}

func (h *CreditCreatedHandler) Handle(
	ctx context.Context,
	payload json.RawMessage,
) error {
	var event event.CreditCreated
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return h.svc.Execute(ctx, credits.ProcessCreditCommand{
		ClientID: event.ClientID,
		CreditID: event.CreditID,
	})
}
