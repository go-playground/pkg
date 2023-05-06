package sqlext

import (
	"context"
	"database/sql"
	resultext "github.com/go-playground/pkg/v5/values/result"
)

// DoTransaction is a helper function that abstracts some complexities of dealing with a transaction and rolling it back.
func DoTransaction[T any](ctx context.Context, opts *sql.TxOptions, conn *sql.DB, fn func(*sql.Tx) resultext.Result[T, error]) resultext.Result[T, error] {
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return resultext.Err[T, error](err)
	}
	result := fn(tx)
	if result.IsErr() {
		_ = tx.Rollback()
		return result
	}
	err = tx.Commit()
	if err != nil {
		return resultext.Err[T, error](err)
	}
	return result
}
