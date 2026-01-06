package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/domain/event"
	"github.com/edgarSucre/mye"
)

type CreditRejectedHandler struct {
	svc credits.RejectCreditService
}

func (h *CreditRejectedHandler) EventName() string {
	return "credit.rejected"
}

func (h *CreditRejectedHandler) Handle(
	ctx context.Context,
	payload json.RawMessage,
) error {
	var rejectEvent event.CreditRejected
	if err := json.Unmarshal(payload, &rejectEvent); err != nil {
		return mye.Wrap(
			err,
			mye.CodeInternal, // there is no point in using a code for retry
			"handle_credit_rejected_failed",
			"failed to unmarshal credit.rejected event",
		)
	}

	h.svc.Execute(ctx, rejectEvent.CreditID)

	return nil
}
