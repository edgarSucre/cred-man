package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/pkg/domain"
)

type CreditApprovedHandler struct {
	creditSvc domain.CreditService
}

func (h *CreditApprovedHandler) EventName() string {
	return "credit.approved"
}

func (h *CreditApprovedHandler) Handle(
	ctx context.Context,
	payload json.RawMessage,
) error {
	var event domain.CreditCreated
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return h.creditSvc.ProcessCredit(ctx, event)
}
