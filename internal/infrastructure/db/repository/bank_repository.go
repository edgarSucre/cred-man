package repository

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/bank"

	"github.com/edgarSucre/mye"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BankRepository struct {
	q *Queries
}

func NewBankRepository(pool *pgxpool.Pool) BankRepository {
	return BankRepository{q: &Queries{db: pool}}
}

func (repo BankRepository) CreateBank(ctx context.Context, bankIn bank.Bank) (bank.Bank, error) {
	model, err := repo.q.CreateBank(ctx, CreateBankParams{
		Name: bankIn.Name(),
		Type: BankType(bankIn.Type().String()),
	})

	if err != nil {
		code, slug := CodeAndSlug(err)
		if code == mye.CodeConflict {
			return bank.Bank{}, mye.Wrap(err, code, slug, "bank fails to be unique").
				WithUserMsg("The name of the bank is already taken. Choose a different name and try again.")
		}

		if code == mye.CodeTimeout {
			return bank.Bank{}, mye.Wrap(err, code, slug, "insert into table banks timeout").
				WithUserMsg("The bank creation is taking a bit too log due to high traffic. Please try again in a few seconds.")
		}

		return bank.Bank{}, mye.Wrap(err, code, slug, "createBank failure")
	}

	dBank, err := model.ToDomain()
	if err != nil {
		return bankIn, err
	}

	return dBank, nil
}

func (m Bank) ToDomain() (bank.Bank, error) {
	id, err := bank.NewID(m.ID.String())
	if err != nil {
		err = mye.Wrap(err, mye.CodeInternal, ErrDataIntegrity, "corrupted bank ID in the database")

		return bank.Bank{}, err
	}

	t, err := bank.TypeFromString(string(m.Type))
	if err != nil {
		err = mye.Wrap(err, mye.CodeInternal, ErrDataIntegrity, "corrupted bank type in the database")

		return bank.Bank{}, fmt.Errorf("bank.TypeFromString > %w", err)
	}

	return bank.Rehydrate(id, m.Name, t), nil

}

func (repo BankRepository) GetBank(
	ctx context.Context,
	id bank.ID,
) (bank.Bank, error) {
	mBank, err := repo.q.GetBank(ctx, id.UUID())
	if err != nil {
		code, slug := CodeAndSlug(err)

		if code == mye.CodeNotFound {
			return bank.Bank{}, mye.Wrap(err, code, slug, "bank not found").
				WithAttribute("id", id.String()).
				WithUserMsg("we couldn't find that bank")

		}

		if code == mye.CodeTimeout {
			return bank.Bank{}, mye.Wrap(err, code, slug, "get bank query timeout").
				WithUserMsg("The bank search is taking a bit too log due to high traffic. Please try again in a few seconds.")
		}

		return bank.Bank{}, mye.Wrap(err, code, slug, "get bank query error")
	}

	dBank, err := mBank.ToDomain()
	if err != nil {
		return bank.Bank{}, err
	}

	return dBank, nil
}
