package repository

import (
	"context"

	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/jackc/pgx/v5"
)

func NewRepository(ctx context.Context, connStr string) (*Queries, error) {
	conn, err := pgx.Connect(context.Background(), connStr)

	return New(conn), errDBConn(err)
}

func errDBConn(err error) error {
	if err != nil {
		return terror.Internal.New(
			"db-connection-error",
			err.Error(),
		)
	}

	return nil
}
