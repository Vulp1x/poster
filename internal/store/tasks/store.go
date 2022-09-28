package tasks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/internal/transport"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

type Store struct {
	tasksChan   chan domain.TaskWithCtx
	taskCancels map[uuid.UUID]func()
	taskMu      *sync.Mutex
	pushTimeout time.Duration
	dbtxf       dbmodel.DBTXFunc
	txf         dbmodel.TxFunc
}

func (s *Store) ForceDelete(ctx context.Context, taskID uuid.UUID) error {
	tx, err := s.txf(ctx)
	if err != nil {
		return store.ErrTransactionFail
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	_, err = q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}
	var deleteCount int64

	deleteCount, err = q.ForceDeleteBotAccountsForTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete bot accounts: %v", err)
	}

	logger.Infof(ctx, "deleted %d bot accounts", deleteCount)

	deleteCount, err = q.ForceDeleteProxiesForTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete proxies: %v", err)
	}

	logger.Infof(ctx, "deleted %d proxies", deleteCount)

	deleteCount, err = q.ForceDeleteTargetUsersForTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete target users: %v", err)
	}

	logger.Infof(ctx, "deleted %d target users", deleteCount)

	err = q.ForceDeleteTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return tx.Commit(ctx)
}

// ErrTaskNotFound не смогли найти таску
var ErrTaskNotFound = errors.New("task not found")

var ErrTaskInvalidStatus = errors.New("invalid task status")

func (s *Store) PrepareTask(
	ctx context.Context,
	taskID uuid.UUID,
	botAccounts domain.BotAccounts,
	proxies domain.Proxies,
	targets domain.TargetUsers,
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
		return err
	}

	savedCount, err = q.SaveProxies(ctx, proxies.ToSaveParams(taskID))
	logger.Infof(ctx, "saved %d proxies", savedCount)
	if err != nil {
		return err
	}

	savedCount, err = q.SaveTargetUsers(ctx, targets.ToSaveParams(taskID))
	logger.Infof(ctx, "saved %d target users", savedCount)
	if err != nil {
		return err
	}

	err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{
		Status: dbmodel.DataUploadedTaskStatus,
		ID:     taskID,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func NewStore(timeout time.Duration, dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc) *Store {
	return &Store{
		tasksChan:   make(chan domain.TaskWithCtx, 10),
		taskCancels: make(map[uuid.UUID]func()),
		pushTimeout: timeout,
		dbtxf:       dbtxFunc,
		txf:         txFunc,
	}
}

const workersPerTask = 10

func (s *Store) CreateDraftTask(ctx context.Context, userID uuid.UUID, title, textTemplate string, image []byte) (uuid.UUID, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	taskID, err := q.CreateDraftTask(ctx, dbmodel.CreateDraftTaskParams{
		ManagerID:    userID,
		TextTemplate: textTemplate,
		Image:        image,
		Title:        title,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create task draft: %w", err)
	}

	return taskID, nil
}

func (s *Store) StartTask(ctx context.Context, taskID uuid.UUID) error {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.StartTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	taskCtx, taskCancel := context.WithCancel(ctx)

	taskWithCtx := domain.TaskWithCtx{
		Task: task,
		Ctx:  taskCtx,
	}

	botsChan := make(chan *domain.TaskPerBot, 20)

	var workers []*worker
	for i := 0; i < workersPerTask; i++ {
		workers = append(workers, &worker{
			tasksQueue:     botsChan,
			dbtxf:          s.dbtxf,
			cli:            transport.InitHTTPClient(),
			processorIndex: int64(i),
		})
	}

	select {
	case s.tasksChan <- taskWithCtx:
		s.taskCancels[task.ID] = taskCancel
		return nil

	case <-time.After(s.pushTimeout):
		logger.Debugf(ctx, "waited for %s, failed to push task to queue")
		break
	}

	taskCancel()

	return fmt.Errorf("failed to push task to queue")
}

func (s *Store) StopTask(ctx context.Context, taskID uuid.UUID) error {
	logger.Infof(ctx, "stopping task '%s'", taskID)
	cancel, ok := s.taskCancels[taskID]
	if !ok {
		return fmt.Errorf("failed to find task '%s' in tasks: %#v", taskID, s.taskCancels)
	}

	cancel()

	s.taskMu.Lock()
	defer s.taskMu.Unlock()

	delete(s.taskCancels, taskID)

	return nil
}
