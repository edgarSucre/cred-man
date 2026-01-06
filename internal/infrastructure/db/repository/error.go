package repository

import (
	"context"
	"errors"

	"github.com/edgarSucre/mye"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

/* ========================================================================================== */
/*                                         error slugs                                        */
/* ========================================================================================== */

const (
	ErrNotFound      = "record_not_found"
	ErrConn          = "database_connection_err"
	ErrDuplicate     = "unique_constraint_violation"
	ErrForeignKey    = "foreign_key_violation"
	ErrSchema        = "schema_mismatch"
	ErrDefault       = "database_error"
	ErrTimeOut       = "query_time_out"
	ErrCancelled     = "query_cancelled"
	ErrDataIntegrity = "data_integrity_error"
)

func CodeAndSlug(err error) (mye.Code, string) {
	if err == nil {
		return 0, ""
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return mye.CodeTimeout, ErrTimeOut
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return mye.CodeNotFound, ErrNotFound
	}

	var connErr *pgconn.ConnectError
	if errors.As(err, &connErr) {
		return mye.CodeUnavailable, ErrConn
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return mye.CodeConflict, ErrDuplicate
		case "23503": // foreign_key_violation
			return mye.CodeInvalid, ErrForeignKey
		case "42P01", "42703": // undefined_table, undefined_column
			return mye.CodeInternal, ErrSchema
		case "08000", "08003", "08006": // connection errors
			return mye.CodeUnavailable, ErrConn
		}
	}

	return mye.CodeInternal, ErrDefault
}
