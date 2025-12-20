package client

import (
	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/pkg/terror"
)

type Service struct {
	repo repository.Querier
}

var ErrNoRepo = terror.Internal.New("clientService-bad-config", "repository is missing")

func NewService(repo repository.Querier) (*Service, error) {
	if repo == nil {
		return nil, ErrNoRepo
	}

	return &Service{repo}, nil
}
