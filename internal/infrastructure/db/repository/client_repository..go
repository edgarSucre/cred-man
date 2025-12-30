package repository

import (
	"context"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	pool *pgxpool.Pool
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{pool}
}

func (repo *ClientRepository) CreateClient(
	ctx context.Context,
	cl client.Client,
) (client.Client, error) {
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return cl, err
	}

	defer conn.Conn().Close(ctx)

	q := &Queries{db: conn}

	model, err := q.CreateClient(ctx, CreateClientParams{
		FullName:  cl.FullName(),
		Email:     cl.Email().String(),
		Birthdate: cl.Birthdate().Time(),
		Country:   cl.Country(),
	})

	if err != nil {
		return cl, err
	}

	dClient, err := model.ToDomain()
	if err != nil {
		return cl, err
	}

	return dClient, nil
}

func (repo *ClientRepository) GetClient(ctx context.Context, id client.ID) (client.Client, error) {
	conn, err := repo.pool.Acquire(ctx)
	if err != nil {
		return client.Client{}, err
	}

	defer conn.Conn().Close(ctx)

	q := &Queries{db: conn}

	mClient, err := q.GetClient(ctx, id.UUID())
	if err != nil {
		return client.Client{}, err
	}

	dClient, err := mClient.ToDomain()
	if err != nil {
		return client.Client{}, err
	}

	return dClient, nil
}

func (m Client) ToDomain() (client.Client, error) {
	birthdate, err := client.NewBirthdate(m.Birthdate)
	if err != nil {
		return client.Client{}, err
	}

	email, err := client.NewEmail(m.Email)
	if err != nil {
		return client.Client{}, err
	}

	id, err := client.NewID(m.ID.String())
	if err != nil {
		return client.Client{}, err
	}

	return client.Rehydrate(
		birthdate,
		m.Country,
		m.CreatedAt,
		email,
		m.FullName,
		id,
	), nil

}
