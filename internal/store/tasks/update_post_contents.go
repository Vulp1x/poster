package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/workers"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
	"github.com/jackc/pgx/v4"
)

func (s *Store) StartUpdatePostContents(ctx context.Context, taskID uuid.UUID) ([]string, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	if task.Status != dbmodel.DoneTaskStatus {
		return nil, fmt.Errorf("got status %d, expected 6(done)", task.Status)
	}

	if err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{
		Status: dbmodel.UpdatingPostContentsTaskStatus,
		ID:     taskID,
	}); err != nil {
		return nil, fmt.Errorf("failed to update task status: %v", err)
	}

	botsCount, err := q.GetTaskBotsCount(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to count bots for task: %v", err)
	}

	aliveBots, err := q.SetBotsEditingPostsStatus(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to update bots statuses to 'editing posts': %v", err)
	}

	logger.Infof(ctx, "have %d alive bots from %d bots in task", len(aliveBots), botsCount)

	tasks := make([]pgqueue.Task, len(aliveBots))
	for i, bot := range aliveBots {
		tasks[i] = pgqueue.Task{
			Kind:        workers.EditMediaTaskKind,
			Payload:     workers.EmptyPayload,
			ExternalKey: fmt.Sprintf("%s::%s", task.ID.String(), bot.Username),
		}
	}

	if err = s.queue.PushTasksTx(ctx, tx, tasks); err != nil {
		return nil, fmt.Errorf("failed to push %d tasks: %v", len(tasks), err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	logger.Infof(ctx, "pushed %d tasks for editing posts", len(tasks))

	return task.LandingAccounts, nil
}
