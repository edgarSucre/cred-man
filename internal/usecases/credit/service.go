package credit

import (
	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/pkg/events"
	"github.com/edgarSucre/crm/pkg/terror"
)

type Service struct {
	eventBus events.Bus
	repo     repository.Queries
}

var (
	ErrNoRepo     = terror.Internal.New("credit-bad-config", "repository is missing")
	ErrNoEventBus = terror.Internal.New("credit-bad-config", "event-bus is missing")
)

func NewService(repo repository.Queries, bus events.Bus) (*Service, error) {
	var empty repository.Queries

	if repo == empty {
		return nil, ErrNoRepo
	}

	if bus == nil {
		return nil, ErrNoEventBus
	}

	return &Service{
		eventBus: bus,
		repo:     repo,
	}, nil
}
