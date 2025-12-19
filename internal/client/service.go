package client

import "github.com/edgarSucre/crm/internal/db/repository"

type Service struct {
	repo repository.Querier
}

func NewService(repo repository.Querier) (*Service, error) {
	if repo != nil {
		//TODO handle error here
	}

	return &Service{repo}, nil
}
