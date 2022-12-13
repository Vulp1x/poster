package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
)

type ParseUsersFromMediaHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   instaproxy.InstaProxyClient
}

func (h *ParseUsersFromMediaHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) != 2 {
		return fmt.Errorf(`%w: expected external key in format "<dataset id>::<media id>", got '%s'`, pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	datasetID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	q := dbmodel.New(h.dbTxF(ctx))
	dataset, err := q.GetDatasetByID(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("failed to find datatset with id '%s': %v", datasetID, err)
	}
	mediaID := externalKeyParts[1]

	media, err := q.FindMediaByID(ctx, mediaID)
	if err != nil {
		return fmt.Errorf("failed to find media with id '%s': %v", mediaID, err)
	}

	targetsResp, err := h.cli.ParseMedia(ctx, &instaproxy.ParseMediaRequest{
		MediaId:       mediaID,
		CommentsCount: dataset.CommentedPerPost,
		LikesCount:    dataset.LikedPerPost,
	})
	if err != nil {
		return fmt.Errorf("failed to parse targets from media '%s': %v", mediaID, err)
	}

	targets := domain.ShortUsersFromProto(targetsResp.GetTargets())

	count, err := q.SaveTargetUsers(ctx, targets.ToSaveTargetsParams(media.Pk, datasetID))
	if err != nil {
		return fmt.Errorf("failed to save targets: %v", err)
	}

	logger.Infof(ctx, "saved %d/%d targets", count, len(targets))

	return nil
}
