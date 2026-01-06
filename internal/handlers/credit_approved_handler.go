package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/domain/event"
	"github.com/edgarSucre/mye"
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
		return mye.Wrap(
			err,
			mye.CodeInternal, // there is no point in using a code for retry
			"handle_credit_approve_failed",
			"failed to unmarshal credit.approved event",
		)
	}

	h.svc.Execute(ctx, approveEvent.CreditID)

	return nil
}
