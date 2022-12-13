package db

import (
	"context"
	"errors"
	"fmt"
)

const tryPassRegistry = `-- pgqueue: TryPassRegistry :exec
INSERT INTO pgqueue_registry (job) VALUES ($1)
ON CONFLICT (job) DO UPDATE
SET updated_at = CASE WHEN pgqueue_registry.updated_at + $2::interval <= now() THEN now() ELSE pgqueue_registry.updated_at END
RETURNING updated_at >= now()
`

// TryPassRegistry пробует зарегистрироваться на выполнение задачи [job].
// Возвращает nil, если это получилось. Критерий этого - прошло хотя бы время [period] с прошлой регистрации.
func (q *Queries) TryPassRegistry(ctx context.Context, job, period string) error {
	var passed bool
	if err := q.executor.QueryRow(ctx, tryPassRegistry, job, period).Scan(&passed); err != nil {
		return err
	} else if !passed {
		return ErrRegistryIsBusy
	}
	return nil
}

func jobCleanupTasks(kind int16) string {
	return fmt.Sprintf("cleanup_tasks_%v", kind)
}

// ErrRegistryIsBusy возвращается,
// если с прошлой регистрации на job прошло недостаточно времени.
var ErrRegistryIsBusy = errors.New("registry is busy")
