package tasks

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

const videosFolder = "videos"

func (s *Store) SaveVideo(ctx context.Context, taskID uuid.UUID, video []byte, filename string) (task domain.Task, err error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to start transaction: %v", err)
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	var dbTask dbmodel.Task

	dbTask, err = q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("failed to find task: %v", err)
	}

	task = domain.Task(dbTask)

	if task.Status == dbmodel.StartedTaskStatus {
		return domain.Task{}, fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.DataUploadedTaskStatus, task.Status)
	}

	if task.Type != dbmodel.ReelsTaskType {
		return domain.Task{}, fmt.Errorf("%w: got type %d, expected %d (reels)", ErrUnexpectedTaskType, task.Type, dbmodel.ReelsTaskType)
	}

	videoFilename := filepath.Join(videosFolder, fmt.Sprintf("%s_%s_%s", task.ID, time.Now().Format(time.RFC3339), filename))
	f, err := os.Create(videoFilename)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to open file at '%s': %v ", err)
	}

	defer func() {
		err2 := f.Close()
		if err2 != nil {
			logger.Errorf(ctx, "failed to close file: %v", err2)
		}
	}()

	_, err = f.Write(video)
	if err != nil {
		return domain.Task{}, fmt.Errorf("faile to write video: %v", err)
	}

	err = q.SetTaskVideoFilename(ctx, dbmodel.SetTaskVideoFilenameParams{VideoFilename: &videoFilename, ID: task.ID})
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to set task vide filename: %v", err)
	}

	logger.Infof(ctx, "saved task video to file '%s'", videoFilename)

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return task, nil
}
