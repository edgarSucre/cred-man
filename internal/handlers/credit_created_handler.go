package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/domain/event"
	"github.com/edgarSucre/mye"
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
		return mye.Wrap(
			err,
			mye.CodeInternal, // there is no point in using a code for retry
			"handle_credit_created_failed",
			"failed to unmarshal credit.created event",
		)
	}

	return h.svc.Execute(ctx, credits.ProcessCreditCommand{
		ClientID: event.ClientID,
		CreditID: event.CreditID,
	})
}
