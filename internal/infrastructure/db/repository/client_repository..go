package repository

import (
	"context"

	"github.com/edgarSucre/crm/internal/domain/client"
	"github.com/edgarSucre/mye"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	q Querier
}

func NewClientRepository(pool *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{New(pool)}
}

func (repo *ClientRepository) CreateClient(
	ctx context.Context,
	cl client.Client,
) (client.Client, error) {
	model, err := repo.q.CreateClient(ctx, CreateClientParams{
		FullName:  cl.FullName(),
		Email:     cl.Email().String(),
		Birthdate: cl.Birthdate().Time(),
		Country:   cl.Country(),
	})

	if err != nil {
		code, slug := CodeAndSlug(err)

		if code == mye.CodeConflict {
			return client.Client{}, mye.Wrap(err, code, slug, "client name is not unique").
				WithUserMsg("The name of the client is already taken. Choose a different name and try again")
		}

		if code == mye.CodeTimeout {
			return client.Client{}, mye.Wrap(err, code, slug, "createClient time out").
				WithUserMsg("The client creation is taking a bit too log due to high traffic. Please try again in a few seconds.")
		}

		return client.Client{}, mye.Wrap(err, code, slug, "create client failure")
	}

	dClient, err := model.ToDomain()
	if err != nil {
		return cl, err
	}

	return dClient, nil
}

func (repo *ClientRepository) GetClient(ctx context.Context, id client.ID) (client.Client, error) {
	mClient, err := repo.q.GetClient(ctx, id.UUID())
	if err != nil {
		code, slug := CodeAndSlug(err)

		if code == mye.CodeNotFound {
			return client.Client{}, mye.Wrap(err, code, slug, "client not found").
				WithAttribute("id", id.String()).
				WithUserMsg("we couldn't find that client")
		}

		if code == mye.CodeTimeout {
			return client.Client{}, mye.Wrap(err, code, slug, "getClient time out").
				WithUserMsg("The client search is taking a bit too long due to high traffic. Please try again in a few seconds")
		}

		return client.Client{}, mye.Wrap(err, code, slug, "failed to retrieve client")
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
		return client.Client{}, mye.Wrap(err, mye.CodeInternal, ErrDataIntegrity, "corrupted birthdate in the database")
	}

	email, err := client.NewEmail(m.Email)
	if err != nil {
		return client.Client{}, mye.Wrap(err, mye.CodeInternal, ErrDataIntegrity, "corrupted email in the database")
	}

	id, err := client.NewID(m.ID.String())
	if err != nil {
		return client.Client{}, mye.Wrap(err, mye.CodeInternal, ErrDataIntegrity, "corrupted ID in the database")
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
