package workers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
)

type PostPhotoHandler struct {
	dbTxF dbmodel.DBTXFunc
	cli   api.InstaProxyClient
}

func (s *PostPhotoHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) == 0 {
		return fmt.Errorf("%w: expected '::' in external key after dataset id in '%s'", pgqueue.ErrMustCancelTask, task.ExternalKey)
	}

	taskID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	fmt.Println(taskID)

	return nil
}
