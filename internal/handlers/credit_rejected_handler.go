package handlers

import (
	"context"
	"encoding/json"

	"github.com/edgarSucre/crm/pkg/domain"
)

type CreditRejectedHandler struct {
	creditSvc domain.CreditService
}

func (h *CreditRejectedHandler) EventName() string {
	return "credit.rejected"
}

func (h *CreditRejectedHandler) Handle(
	ctx context.Context,
	payload json.RawMessage,
) error {
	var event domain.CreditRejected
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return h.creditSvc.ProcessRejection(ctx, event)
}
