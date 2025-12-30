package banks

import (
	"context"

	"github.com/edgarSucre/crm/internal/domain/bank"
)

type Repository interface {
	CreateBank(context.Context, bank.Bank) (bank.Bank, error)
}
