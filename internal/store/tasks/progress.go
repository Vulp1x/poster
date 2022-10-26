package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/jackc/pgx/v4"
)

func (s *Store) TaskProgress(ctx context.Context, taskID uuid.UUID) (domain.TaskProgress, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TaskProgress{}, ErrTaskNotFound
		}
	}

	if task.Status < dbmodel.StartedTaskStatus {
		return domain.TaskProgress{}, fmt.Errorf("%w: expected at least %d status, got %d",
			ErrTaskInvalidStatus, dbmodel.StartedTaskStatus, task.Status)
	}

	progress, err := q.GetBotsProgress(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, err
	}

	targetCounters, err := q.GetTaskTargetsCount(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, fmt.Errorf("failed to get target counters: %v", err)
	}

	return domain.TaskProgress{
		BotsProgress:   progress,
		TargetCounters: targetCounters,
		Done:           task.Status > dbmodel.StartedTaskStatus,
	}, nil
}
