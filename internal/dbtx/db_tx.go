package dbtx

import (
	"context"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

// RollbackUnlessCommitted rollsback the transaction regardless
// it has already been committed or rolled back.
//
// Useful to defer tx.RollbackUnlessCommitted(), so you don't
// have to handle N failure cases.
// Keep in mind the only way to detect an error on the rollback
// is via the log.
func RollbackUnlessCommitted(ctx context.Context, tx dbmodel.Tx) {
	err := tx.Rollback(ctx)
	if err == pgx.ErrTxClosed {
		return
	}

	if err != nil {
		logger.Errorf(ctx, "rollback transaction: %s", err)
	}
}
