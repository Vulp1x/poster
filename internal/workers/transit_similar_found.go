package workers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
)

type TransitToSimilarFoundHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
}

func (h *TransitToSimilarFoundHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	datasetID, err := uuid.Parse(task.ExternalKey)
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, task.ExternalKey, err)
	}

	q := dbmodel.New(h.dbTxF(ctx))

	notReadyBloggers, err := q.FindNotReadyBloggers(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find not ready bloggers: %v", err)
	}

	if len(notReadyBloggers) != 0 {
		return fmt.Errorf("dataset %s still has %d not ready bloggers: %v", datasetID,
			len(notReadyBloggers), domain.DatasetWithBloggers{Bloggers: notReadyBloggers}.Usernames())
	}

	err = q.UpdateDatasetStatus(ctx, dbmodel.UpdateDatasetStatusParams{Status: dbmodel.ReadyForParsingDatasetStatus, ID: datasetID})
	if err != nil {
		return fmt.Errorf("failed to update dataset status: %v", err)
	}

	return nil
}
