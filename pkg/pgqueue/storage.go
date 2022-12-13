package pgqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue/internal/db"
	"github.com/inst-api/poster/pkg/pgqueue/pkg/executor"
)

const (
	// Таймаут на выполнение запросов к БД.
	storageTimeout = 5 * time.Second
	// Таймаут на ретрай запросов к БД.
	storageBackoffMaxElapsedTime = 3 * storageTimeout
)

type storage struct {
	*db.Queries
}

func newStorage(executor executor.Executor) *storage {
	return &storage{db.New(executor)}
}

// CancelTask переводит задачу в статус 'cancelled', если она была открытой.
func (s *storage) CancelTask(parentCtx context.Context, arg db.CancelTaskParams) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.CancelTask(ctx, arg)
	})
}

// CancelTaskByKey переводит задачу в статус 'cancelled', если она была открытой.
func (s *storage) CancelTaskByKey(parentCtx context.Context, arg db.CancelTaskByKeyParams) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.CancelTaskByKey(ctx, arg)
	})
}

// CleanupTasks удаляет задачи, которые находятся в терминальном статусе больше заданного времени.
func (s *storage) CleanupTasks(parentCtx context.Context, arg db.CleanupTasksParams) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.CleanupTasks(ctx, arg)
	})
}

// CompleteTask переводит задачу в статус 'succeeded'.
func (s *storage) CompleteTask(parentCtx context.Context, id int64) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.CompleteTask(ctx, id)
	})
}

// DeleteTask удаляет задачу по id.
func (s *storage) DeleteTask(parentCtx context.Context, id int64) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.DeleteTask(ctx, id)
	})
}

// GetOpenTasks возвращает задачи, ждущие выполнения.
// Не нуждается в ретрае, так как запускается часто и не является критическим.
func (s *storage) GetOpenTasks(parentCtx context.Context, params db.GetOpenTasksParams) ([]Task, error) {
	ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
	defer cancel()

	rows, err := s.Queries.GetOpenTasks(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.Queries.GetOpenTasks failed: %w", err)
	}

	tasks := make([]Task, len(rows))
	for i, row := range rows {
		tasks[i] = Task{
			Kind:        row.Kind,
			Payload:     row.Payload,
			ExternalKey: row.ExternalKey.String,

			id:              row.ID,
			attemptsLeft:    row.AttemptsLeft,
			attemptsElapsed: row.AttemptsElapsed,
		}
	}

	return tasks, err
}

// RefuseTask в зависимости от числа оставшихся попыток
// переводит задачу в статус 'no_attempts_left' или 'must_retry'.
func (s *storage) RefuseTask(parentCtx context.Context, arg db.RefuseTaskParams) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.RefuseTask(ctx, arg)
	})
}

// RetryTasks обновляет количество попыток у задач в статусе 'no_attempts_left',
// переводя их в статус `must_retry` в порядке добавления в очередь.
func (s *storage) RetryTasks(parentCtx context.Context, arg db.RetryTasksParams) error {
	return retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		return s.Queries.RetryTasks(ctx, arg)
	})
}

// PushTasks добавляет задачи в очередь батчом. При конфликте по ключу идемпотентности - DO NOTHING.
func (s *storage) PushTasks(parentCtx context.Context, kds map[int16]kindData, tasks []Task, opts ...PushOption) error {
	dbParams, err := getPushManyParams(kds, tasks, opts...)
	if err != nil {
		return err
	}

	if err := retry(func() error {
		ctx, cancel := context.WithTimeout(parentCtx, storageTimeout)
		defer cancel()

		batch := s.Queries.PushTasks(ctx, dbParams)

		batch.Exec(func(i int, err2 error) {
			if err2 != nil {
				err = fmt.Errorf("failed to push task: %v", err2)

				err2 = batch.Close()
				if err2 != nil {
					logger.Error(ctx, "failed to close batch after error: %v", err2)
				}
			}
		})

		return err
	}); err != nil {
		return fmt.Errorf("PushTasks failed: %w", err)
	}

	collectTasksInStatusNew(kds, tasks)

	return nil
}

func getPushManyParams(kds map[int16]kindData, tasks []Task, opts ...PushOption) ([]db.PushTasksParams, error) {
	dbParams := make([]db.PushTasksParams, len(tasks))
	for i, task := range tasks {
		kd, ok := kds[task.Kind]
		if !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedTaskKind, task.Kind)
		}

		params := db.PushTasksParams{
			Kind:         task.Kind,
			Payload:      task.Payload,
			ExternalKey:  db.NullString(task.ExternalKey),
			AttemptsLeft: kd.opts.MaxAttempts,
			DelayedTill:  time.Now(),
		}

		for _, opt := range opts {
			params = opt(params)
		}

		dbParams[i] = params
	}

	return dbParams, nil
}

type PushOption func(db.PushTasksParams) db.PushTasksParams

func WithDelay(delay time.Duration) PushOption {
	return func(params db.PushTasksParams) db.PushTasksParams {
		params.DelayedTill = time.Now().Add(delay)
		return params
	}
}

func collectTasksInStatusNew(kds map[int16]kindData, tasks []Task) {
	counts := make(map[int16]int32)
	for _, task := range tasks {
		counts[task.Kind]++
	}

	// for kind, count := range counts {
	// mc.CollectTasksInStatus(kds[kind].opts.Name, status.New, count)
	// }
}

func retry(fn func() error) error {
	return backoff.Retry(fn, getBackoff())
}

func getBackoff() backoff.BackOff {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = storageBackoffMaxElapsedTime
	return expBackoff
}
