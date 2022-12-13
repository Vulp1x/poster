package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
)

type ParseMediasHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
	queue *pgqueue.Queue
}

// HandleTask обрабатывает задачу на поиск медиа у блогера
// ожидается external key по формату: "<dataset-id>::<blogger username>"
func (p *ParseMediasHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) == 0 {
		return fmt.Errorf("%w: expected '::' in external key after dataset id in '%s'", pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	datasetID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	dataset, err := dbmodel.New(p.dbTxF(ctx)).GetDatasetByID(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find dataset with id '%s': %v", datasetID, err)
	}

	resp, err := p.cli.GetBloggerMedias(ctx, &instaproxy.GetBloggerMediasRequest{
		Username:    externalKeyParts[1],
		MediasCount: dataset.PostsPerBlogger,
	})
	if err != nil {
		return fmt.Errorf("failed to get blogger's media: %v", err)
	}

	tx, err := p.dbTxF(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin new transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	domainMedias := domain.MediasFromProto(resp.Medias, datasetID)

	mediaTasks := make([]pgqueue.Task, 0, len(domainMedias))

	batch := q.SaveMedias(ctx, domainMedias.ToSaveMediasParams())
	batch.QueryRow(func(i int, m dbmodel.Media, err error) {
		if i >= len(domainMedias) {
			err = batch.Close()
			if err != nil {
				logger.Errorf(ctx, "failed to close saving medias batch: ")
			}
		}

		if err != nil {
			logger.Errorf(ctx, "failed to save media %d '%s': %v", i+1, domainMedias[i].ID, err)
			return
		}

		mediaTasks = append(mediaTasks, pgqueue.Task{
			Kind:        ParseUsersFromMediaTaskKind,
			Payload:     task.Payload,
			ExternalKey: fmt.Sprintf("%s::%s", datasetID, m.ID),
		})
	})

	logger.Infof(ctx, "saved %d/%d medias", len(mediaTasks), len(domainMedias))

	err = p.queue.PushTasksTx(ctx, tx, mediaTasks)
	if err != nil {
		return fmt.Errorf("failed to push tasks: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
