package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v4"
)

// primary key
const cancelTask = `-- pgqueue: CancelTask :exec
UPDATE pgqueue SET
	status = 'cancelled', messages = array_append(messages, $1::text), updated_at = NOW()
WHERE id = $2 AND status IN ('new', 'must_retry')
`

type CancelTaskParams struct {
	Reason string
	ID     int64
}

// CancelTask переводит задачу в статус 'cancelled', если она была открытой.
func (q *Queries) CancelTask(ctx context.Context, arg CancelTaskParams) error {
	_, err := q.executor.Exec(ctx, cancelTask, arg.Reason, arg.ID)

	return err
}

// pgqueue_idempotency_idx
const cancelTaskByKey = `-- pgqueue: CancelTaskByKey :exec
UPDATE pgqueue SET
	status = 'cancelled', messages = array_append(messages, $1::text), updated_at = NOW()
WHERE kind = $2 AND external_key = $3 AND status IN ('new', 'must_retry')
`

type CancelTaskByKeyParams struct {
	Reason      string
	Kind        int16
	ExternalKey sql.NullString
}

// CancelTaskByKey переводит задачу в статус 'cancelled', если она была открытой.
func (q *Queries) CancelTaskByKey(ctx context.Context, arg CancelTaskByKeyParams) error {
	_, err := q.executor.Exec(ctx, cancelTaskByKey, arg.Reason, arg.Kind, arg.ExternalKey)

	return err
}

// pgqueue_terminal_tasks_idx
const cleanupTasks = `-- pgqueue: CleanupTasks :exec
DELETE FROM pgqueue
WHERE kind = $1 AND updated_at <= NOW() - $2::interval AND status IN ('cancelled', 'succeeded')
`

type CleanupTasksParams struct {
	Kind    int16
	Timeout string
	Period  string
}

// CleanupTasks удаляет задачи, которые находятся в терминальном статусе больше заданного времени.
func (q *Queries) CleanupTasks(ctx context.Context, arg CleanupTasksParams) error {
	if err := q.TryPassRegistry(ctx, jobCleanupTasks(arg.Kind), arg.Period); err != nil {
		return err
	}
	_, err := q.executor.Exec(ctx, cleanupTasks, arg.Kind, arg.Timeout)

	return err
}

// primary key
const completeTask = `-- pgqueue: CompleteTask :exec
UPDATE pgqueue SET status = 'succeeded', updated_at = NOW() WHERE id = $1
`

// CompleteTask переводит задачу в статус 'succeeded'.
func (q *Queries) CompleteTask(ctx context.Context, id int64) error {
	_, err := q.executor.Exec(ctx, completeTask, id)

	return err
}

// primary key
const deleteTask = `-- pgqueue: DeleteTask :exec
DELETE FROM pgqueue WHERE id = $1;
`

// DeleteTask удаляет задачу по id.
func (q *Queries) DeleteTask(ctx context.Context, id int64) error {
	_, err := q.executor.Exec(ctx, deleteTask, id)

	return err
}

// pgqueue_open_tasks_idx
const getOpenTasks = `-- pgqueue: GetOpenTasks :many
UPDATE pgqueue SET
	status = (CASE WHEN attempts_left > 1 THEN 'must_retry' ELSE 'no_attempts_left' END)::pgqueue_status,
	attempts_left = attempts_left - 1,
	attempts_elapsed = attempts_elapsed + 1,
	delayed_till = NOW() + $3::interval,
	updated_at = NOW()
WHERE id IN (
	SELECT id FROM pgqueue sub
	WHERE
		sub.kind = $1 AND
		sub.delayed_till <= NOW() AND
		sub.status IN ('new', 'must_retry')
	ORDER BY sub.delayed_till ASC
	LIMIT $2
	FOR NO KEY UPDATE SKIP LOCKED
)
RETURNING id, kind, payload, external_key, attempts_left, attempts_elapsed
`

type GetOpenTasksParams struct {
	Kind          int16
	Limit         int32
	UntilDeadline string
}

type GetOpenTasksRow struct {
	ID              int64
	Kind            int16
	Payload         []byte
	ExternalKey     sql.NullString
	AttemptsLeft    int16
	AttemptsElapsed int16
}

// GetOpenTasks возвращает задачи, ждущие выполнения. Строки сортируются по 'delayed_till',
// поэтому сперва возвращаются задачи, которые еще ни разу не выполнялись.
func (q *Queries) GetOpenTasks(ctx context.Context, arg GetOpenTasksParams) ([]GetOpenTasksRow, error) {
	rows, err := q.executor.Query(ctx, getOpenTasks, arg.Kind, arg.Limit, arg.UntilDeadline)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOpenTasksRow
	for rows.Next() {
		var i GetOpenTasksRow
		if err := rows.Scan(
			&i.ID,
			&i.Kind,
			&i.Payload,
			&i.ExternalKey,
			&i.AttemptsLeft,
			&i.AttemptsElapsed,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// primary key
const refuseTask = `-- pgqueue: RefuseTask :exec
UPDATE pgqueue SET
	status = (CASE WHEN attempts_left = 0 THEN 'no_attempts_left' ELSE 'must_retry' END)::pgqueue_status,
	messages =
		CASE WHEN $1::smallint > 0 AND cardinality(messages) >= $1::smallint THEN
			(array_append(messages, $2::text))[2:]
		ELSE
			array_append(messages, $2::text)
		END,
	delayed_till = NOW() + $3::interval,
	updated_at = NOW()
WHERE id = $4 AND status <> 'cancelled'
`

type RefuseTaskParams struct {
	MessagesLimit int16
	Reason        string
	Delay         string
	ID            int64
}

// RefuseTask в зависимости от числа оставшихся попыток
// переводит задачу в статус 'no_attempts_left' или 'must_retry'.
func (q *Queries) RefuseTask(ctx context.Context, arg RefuseTaskParams) error {
	_, err := q.executor.Exec(ctx, refuseTask, arg.MessagesLimit, arg.Reason, arg.Delay, arg.ID)

	return err
}

// pgqueue_broken_tasks_idx
const retryTasks = `-- pgqueue: RetryTasks :exec
UPDATE pgqueue SET status = 'must_retry', attempts_left = $2, attempts_elapsed = 0, updated_at = NOW()
WHERE id IN (
	SELECT id FROM pgqueue sub
	WHERE sub.kind = $1 AND sub.status IN ('no_attempts_left')
	ORDER BY created_at ASC
	LIMIT $3
	FOR NO KEY UPDATE SKIP LOCKED
)
`

type RetryTasksParams struct {
	Kind         int16
	AttemptsLeft int16
	Limit        int32
}

// RetryTasks обновляет количество попыток у задач в статусе 'no_attempts_left',
// переводя их в статус `must_retry` в порядке добавления в очередь.
func (q *Queries) RetryTasks(ctx context.Context, arg RetryTasksParams) error {
	_, err := q.executor.Exec(ctx, retryTasks, arg.Kind, arg.AttemptsLeft, arg.Limit)

	return err
}

// pgqueue_idempotency_idx
const pushTasks = `-- pgqueue: PushTasks :exec
INSERT INTO pgqueue (kind, payload, external_key, attempts_left, delayed_till)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING
`

type PushTasksBatchResults struct {
	br  pgx.BatchResults
	ind int
}

type PushTasksParams struct {
	Kind         int16          `json:"kind"`
	Payload      []byte         `json:"payload"`
	ExternalKey  sql.NullString `json:"external_key"`
	AttemptsLeft int16          `json:"attempts_left"`
	DelayedTill  time.Time      `json:"delayed_till"`
}

func (q *Queries) PushTasks(ctx context.Context, arg []PushTasksParams) *PushTasksBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.Kind,
			a.Payload,
			a.ExternalKey,
			a.AttemptsLeft,
			a.DelayedTill,
		}
		batch.Queue(pushTasks, vals...)
	}
	br := q.executor.SendBatch(ctx, batch)
	return &PushTasksBatchResults{br, 0}
}

func (b *PushTasksBatchResults) Exec(f func(int, error)) {
	for {
		_, err := b.br.Exec()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, err)
		}
		b.ind++
	}
}

func (b *PushTasksBatchResults) Close() error {
	return b.br.Close()
}
