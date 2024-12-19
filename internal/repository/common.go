package repository

import (
	"context"
	"database/sql"
)

// DB is common interface for both sql.DB and sql.Tx.
// It is used to abstract the database operations from the repository layer.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// getExecutor returns the executor based on the context.
// If the context has a transaction, it returns the transaction.
// Otherwise, it returns the database.
func getExecutor(ctx context.Context, db *sql.DB) DB {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}
