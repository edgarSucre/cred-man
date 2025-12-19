package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewRepository(ctx context.Context, connStr string) (*Queries, error) {
	conn, err := pgx.Connect(context.Background(), connStr)

	return New(conn), err
}
