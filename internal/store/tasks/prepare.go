package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

func (s *Store) PrepareTask(
	ctx context.Context,
	taskID uuid.UUID,
	botAccounts domain.BotAccounts,
	residentialProxies, cheapProxies domain.Proxies,
	targets domain.TargetUsers,
	filenames *tasksservice.TaskFileNames,
) error {
	tx, err := s.txf(ctx)
	if err != nil {
		return store.ErrTransactionFail
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if task.Status != dbmodel.DraftTaskStatus {
		return fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.DraftTaskStatus, task.Status)
	}

	savedCount, err := q.SaveBotAccounts(ctx, botAccounts.ToSaveParams(taskID))
	logger.Infof(ctx, "saved %d bots", savedCount)
	if err != nil {
		return fmt.Errorf("failed to save bots: %v", err)
	}

	savedCount, err = q.SaveProxies(ctx, residentialProxies.ToSaveParams(taskID, false))
	logger.Infof(ctx, "saved %d residential proxies", savedCount)
	if err != nil {
		return fmt.Errorf("failed to save residential proxies: %v", err)
	}

	savedCount, err = q.SaveProxies(ctx, cheapProxies.ToSaveParams(taskID, true))
	logger.Infof(ctx, "saved %d cheap proxies", savedCount)
	if err != nil {
		return fmt.Errorf("failed to save cheap proxies: %v", err)
	}

	savedCount, err = q.SaveTargetUsers(ctx, targets.ToSaveParams(taskID))
	logger.Infof(ctx, "saved %d target users", savedCount)
	if err != nil {
		return fmt.Errorf("failed to save targets: %v", err)
	}

	err = q.SaveUploadedDataToTask(ctx, dbmodel.SaveUploadedDataToTaskParams{
		ID:                   taskID,
		BotsFilename:         &filenames.BotsFilename,
		ResProxiesFilename:   &filenames.ResidentialProxiesFilename,
		CheapProxiesFilename: &filenames.CheapProxiesFilename,
		TargetsFilename:      &filenames.TargetsFilename,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
