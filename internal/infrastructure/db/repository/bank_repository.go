package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BankRepository struct {
	q *Queries
}

func NewBankRepository(pool *pgxpool.Pool) BankRepository {
	return BankRepository{q: &Queries{db: pool}}
}

func (repo BankRepository) CreateBank(ctx context.Context, bank bank.Bank) (bank.Bank, error) {
	model, err := repo.q.CreateBank(ctx, CreateBankParams{
		Name: bank.Name(),
		Type: BankType(bank.Type().String()),
	})

	if err != nil {
		return bank, err
	}

	dBank, err := model.ToDomain()
	if err != nil {
		return bank, err
	}

	return dBank, nil
}

func (m Bank) ToDomain() (bank.Bank, error) {
	id, err := bank.NewID(m.ID.String())
	if err != nil {
		err = terror.ToInternal(err)
		return bank.Bank{}, fmt.Errorf("bank.NewID > %w", err)
	}

	t, err := bank.TypeFromString(string(m.Type))
	if err != nil {
		err = terror.ToInternal(err)
		return bank.Bank{}, fmt.Errorf("bank.TypeFromString > %w", err)
	}

	return bank.Rehydrate(id, m.Name, t), nil

}

var ErrBankNotFound = terror.NotFound.New("bank_not_found", "couldn't find bank")

func (repo BankRepository) GetBank(
	ctx context.Context,
	id bank.ID,
) (bank.Bank, error) {
	mBank, err := repo.q.GetBank(ctx, id.UUID())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return bank.Bank{}, ErrBankNotFound
		}

		return bank.Bank{}, fmt.Errorf("sqlc.getBank > %w", err)
	}

	dBank, err := mBank.ToDomain()
	if err != nil {
		return bank.Bank{}, fmt.Errorf("sqlc.bank.ToDomain > %w", err)
	}

	return dBank, nil
}
