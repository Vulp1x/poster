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

func (s *Store) StartBots(ctx context.Context, taskID uuid.UUID, usernames []string) ([]string, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("задача не найдена")
		}

		return nil, fmt.Errorf("failed to find task by id: %v", err)
	}

	bots, err := q.FindTaskBotsByUsername(ctx, dbmodel.FindTaskBotsByUsernameParams{TaskID: taskID, Usernames: usernames})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("не найдено ни одного бота")
		}

		return nil, fmt.Errorf("failed to find bots: %v", err)
	}

	postingTasks := preparePostingTasks(ctx, task, bots)
	err = s.queue.PushTasksTx(ctx, tx, postingTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to push tasks to queue: %v", err)
	}

	logger.Infof(ctx, "pushed %d tasks from %d input usernames", len(postingTasks), len(usernames))

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	var succeededUsernames = make([]string, len(bots))
	for i, bot := range bots {
		succeededUsernames[i] = bot.Username
	}

	return succeededUsernames, nil
}

func preparePostingTasks(ctx context.Context, task dbmodel.Task, bots []dbmodel.BotAccount) []pgqueue.Task {
	postingTasks := make([]pgqueue.Task, len(bots))
	for i, bot := range bots {
		postingTasks[i] = pgqueue.Task{
			Kind:        workers.MakePhotoPostsTaskKind,
			Payload:     workers.EmptyPayload,
			ExternalKey: fmt.Sprintf("%s::%s::0", task.ID.String(), bot.Username),
		}
	}

	return postingTasks
}
