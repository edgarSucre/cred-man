package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) transactionManager {
	return transactionManager{pool}
}

type txKey struct{}

func (tm transactionManager) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) (err error) {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if cErr := tx.Rollback(ctx); cErr != nil && err == nil {
			err = cErr
		}
	}()

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
