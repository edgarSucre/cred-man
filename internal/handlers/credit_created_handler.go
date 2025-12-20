package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/pkg/domain"
)

type CreditCreatedHandler struct {
	creditSvc domain.CreditService
}

func (h *CreditCreatedHandler) EventName() string {
	return "credit.created"
}

func (h *CreditCreatedHandler) Handle(
	ctx context.Context,
	payload json.RawMessage,
) error {
	var event domain.CreditCreated
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return h.creditSvc.ProcessCredit(ctx, event)
}
