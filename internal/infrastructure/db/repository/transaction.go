package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (q *Queries) WithTransaction(ctx context.Context, fn func(tx *Queries) error) error {
	db := q.db.(*pgx.Conn)

	transaction, err := db.Begin(ctx)
	defer transaction.Rollback(ctx)

	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	if err := fn(q.WithTx(transaction)); err != nil {
		return fmt.Errorf("repository.Tx error: %w", err)
	}

	return nil
}
