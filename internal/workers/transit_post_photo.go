package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
)

type TransitPostPhotoHandler struct {
	dbTxF         dbmodel.DBTXFunc
	cli           api.InstaProxyClient
	instagrapiCLi instagrapi.Client
	queue         *pgqueue.Queue
}

func (s *TransitPostPhotoHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)
	taskID, err := uuid.Parse(task.ExternalKey)
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, task.ExternalKey, err)
	}

	q := dbmodel.New(s.dbTxF(ctx))

	notFinishedBots, err := q.FindNotFinishedPostingTaskBots(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to find bots in progress: %v", err)
	}

	if len(notFinishedBots) != 0 {
		logger.Warnf(ctx, "still has %d bots in not finished statuses: %v", len(notFinishedBots), notFinishedBots)

		err = s.queue.RetryTasks(ctx, MakePhotoPostsTaskKind, 10, 1000)
		if err != nil {
			return fmt.Errorf("failed to retry posting tasks: %v", err)
		}

		return fmt.Errorf("still has %d bots in not finished statuses: %v", len(notFinishedBots), notFinishedBots)
	}

	err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{Status: dbmodel.DoneTaskStatus, ID: taskID})
	if err != nil {
		return fmt.Errorf("failed to update task status: %v", err)
	}

	return nil
}
