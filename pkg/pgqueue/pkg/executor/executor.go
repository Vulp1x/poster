// Легкое представление исполнителя запросов к БД.
// Необходимо для обеспечения совместимости с различными библиотеками (e.g. database-go, database-pg).
package executor

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Tx - легкое представление транзакции.
type Tx interface {
	Executor
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Row - легкое представление строки.
type Row interface {
	Scan(...interface{}) error
}

// Rows - легкое представление строк.
type Rows interface {
	Row
	Next() bool
	Err() error
	Close() error
}

type Executor interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}
