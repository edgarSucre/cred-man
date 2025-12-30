package credits

import (
	"context"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/crm/internal/domain/credit"
	"github.com/edgarSucre/crm/internal/domain/event"
)

type creditRepository interface {
	CreateCredit(context.Context, credit.Credit) (credit.Credit, error)
	GetAggregate(context.Context, credit.ID, client.ID) (*credit.CreditAggregate, error)
	GetCredit(context.Context, credit.ID) (credit.Credit, error)
	ProcessCredit(context.Context, credit.CreditAggregate) error
}

type transactionManager interface {
	WithTransaction(context.Context, func(context.Context) error) error
}

type clientRepository interface {
	GetClient(context.Context, client.ID) (client.Client, error)
}

type bankRepository interface {
	GetBank(context.Context, bank.ID) (bank.Bank, error)
}

type eventBus interface {
	Publish(context.Context, event.Event) error
}
