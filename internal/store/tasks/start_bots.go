package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	api "github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/internal/workers"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
	"github.com/jackc/pgx/v4"
)

func (s *Store) StartBots(ctx context.Context, taskID uuid.UUID, usernames []string) ([]string, error) {
	q := dbmodel.New(s.db)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("задача не найдена")
		}

		return nil, fmt.Errorf("failed to find task by id: %v", err)
	}

	bots, err := q.FindTaskBotsByUsername(ctx, dbmodel.FindTaskBotsByUsernameParams{TaskID: taskID, Usernames: usernames})
	if err != nil {
		return nil, fmt.Errorf("failed to find bots: %v", err)
	}

	if len(bots) == 0 {
		return nil, fmt.Errorf("не найдено ни одного бота")
	}

	err = s.checkAndUpdateTaskLandingAccounts(ctx, task, q)
	if err != nil {
		return nil, err
	}

	tx, err := s.txf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q = dbmodel.New(tx)

	err = q.SetBotsStatus(ctx, dbmodel.SetBotsStatusParams{Status: dbmodel.StartedBotStatus, Ids: domain.Ids(bots)})
	if err != nil {
		return nil, fmt.Errorf("failed to set bot statuses to 'started': %v", err)
	}

	postingTasks := preparePostingTasks(ctx, task, bots)
	err = s.queue.PushTasksTx(ctx, tx, postingTasks)
	if err != nil {
		return nil, fmt.Errorf("failed to push tasks to queue: %v", err)
	}

	if err = s.queue.PushTaskTx(ctx, tx, pgqueue.Task{
		Kind:        workers.TransitPostsMadeTaskKind,
		Payload:     workers.EmptyPayload,
		ExternalKey: taskID.String(),
	}); err != nil {
		return nil, fmt.Errorf("failed to push 'update status after all bots done' task: %v", err)
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

func (s *Store) checkAndUpdateTaskLandingAccounts(ctx context.Context, task dbmodel.Task, q *dbmodel.Queries) error {
	checkedLandings, err := s.cli.CheckLandings(ctx, &api.CheckLandingsRequest{Usernames: task.LandingAccounts})
	if err != nil {
		return fmt.Errorf("failed to check landings: %v", err)
	}

	if len(task.LandingAccounts) > len(checkedLandings.AliveLandings) {
		logger.Warnf(ctx, "before check got %d landing accounts, but only %d are alive: %v",
			len(task.LandingAccounts), len(checkedLandings.AliveLandings), checkedLandings.AliveLandings,
		)

		err = q.UpdateTaskLandingAccounts(ctx, dbmodel.UpdateTaskLandingAccountsParams{LandingAccounts: checkedLandings.AliveLandings, ID: task.ID})
		if err != nil {
			return fmt.Errorf("failed to update task landing accounts: %v", err)
		}
	}
	return nil
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
