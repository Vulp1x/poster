package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
)

type TransitEditPhotoHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   api.InstaProxyClient
	queue *pgqueue.Queue
}

func (s *TransitEditPhotoHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)
	taskID, err := uuid.Parse(task.ExternalKey)
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, task.ExternalKey, err)
	}

	q := dbmodel.New(s.dbTxF(ctx))

	notFinishedBots, err := q.FindNotFinishedEditingTaskBots(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to find bots in progress: %v", err)
	}

	if len(notFinishedBots) != 0 {
		logger.Warnf(ctx, "still has %d bots in not all done status: %v", len(notFinishedBots), notFinishedBots)

		err = s.queue.RetryTasks(ctx, EditMediaTaskKind, 10, 1000)
		if err != nil {
			return fmt.Errorf("failed to retry posting tasks: %v", err)
		}

		return nil
	}

	err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{Status: dbmodel.AllDoneTaskStatus, ID: taskID})
	if err != nil {
		return fmt.Errorf("failed to update task status: %v", err)
	}

	return nil
}
