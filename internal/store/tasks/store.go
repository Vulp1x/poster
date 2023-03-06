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
	"github.com/inst-api/poster/internal/instagrapi"
	api "github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const workersPerTask = 10
const maxDBLimit = 1000000

// ErrTaskNotFound не смогли найти таску
var ErrTaskNotFound = errors.New("task not found")

// ErrTaskInvalidStatus переход по статусам не возможен
var ErrTaskInvalidStatus = errors.New("invalid task status")

// ErrUnexpectedTaskType ожидали другой тип таски
var ErrUnexpectedTaskType = errors.New("unexpected task type")

func NewStore(timeout time.Duration, dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc, instagrapiHost string, conn *grpc.ClientConn, queue *pgqueue.Queue, db dbmodel.DBTX) *Store {
	return &Store{
		tasksChan:   make(chan domain.Task, 10),
		taskCancels: make(map[uuid.UUID]func()),
		pushTimeout: timeout,
		dbtxf:       dbtxFunc,
		txf:         txFunc,
		taskMu:      &sync.Mutex{},
		instaClient: instagrapi.NewClient(instagrapiHost),
		cli:         api.NewInstaProxyClient(conn),
		queue:       queue,
		db:          db,
	}
}

type Store struct {
	tasksChan   chan domain.Task
	taskCancels map[uuid.UUID]func()
	taskMu      *sync.Mutex
	pushTimeout time.Duration
	dbtxf       dbmodel.DBTXFunc
	txf         dbmodel.TxFunc
	instaClient instagrapiClient
	cli         api.InstaProxyClient
	queue       *pgqueue.Queue
	db          dbmodel.DBTX
}

func (s *Store) TaskBots(ctx context.Context, taskID uuid.UUID) (domain.BotAccounts, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	_, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	bots, err := q.FindBotsForTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bots: %v", err)
	}

	return domain.BotsFromDBModels(bots), nil
}

func (s *Store) TaskTargets(ctx context.Context, taskID uuid.UUID) (domain.Targets, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	_, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	targets, err := q.FindUnprocessedTargetsForTask(ctx, dbmodel.FindUnprocessedTargetsForTaskParams{
		TaskID: taskID,
		Limit:  maxDBLimit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find targets ")
	}

	return targets, nil
}

func (s *Store) ListTasks(ctx context.Context, userID uuid.UUID) (domain.TasksWithCounters, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	tasks, err := q.FindTasksByManagerID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("failed to find tasks for manager '%s': %v", userID, err)
	}

	domainTasks := make([]domain.TaskWithCounters, len(tasks))
	for i, task := range tasks {
		domainTasks[i] = domain.TaskWithCounters{Task: task}
	}

	return domainTasks, nil
}

func (s *Store) GetTask(ctx context.Context, taskID uuid.UUID) (domain.TaskWithCounters, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TaskWithCounters{}, ErrTaskNotFound
		}

		return domain.TaskWithCounters{}, fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	counters, err := q.SelectCountsForTask(ctx, taskID)
	if err != nil {
		return domain.TaskWithCounters{}, fmt.Errorf("failed to select counters: %v", err)
	}

	return domain.TaskWithCounters{Task: task, SelectCountsForTaskRow: counters}, nil
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

	deleteCount, err = q.ForceDeleteProxiesForTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete proxies: %v", err)
	}

	logger.Infof(ctx, "deleted %d proxies", deleteCount)

	deleteCount, err = q.ForceDeleteBotAccountsForTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete bot accounts: %v", err)
	}

	logger.Infof(ctx, "deleted %d bot accounts", deleteCount)

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

func (s *Store) StopTask(ctx context.Context, taskID uuid.UUID) error {
	logger.Infof(ctx, "stopping task '%s'", taskID)
	s.taskMu.Lock()
	cancel, ok := s.taskCancels[taskID]
	s.taskMu.Unlock()
	if !ok {
		return fmt.Errorf("failed to find task '%s' in tasks: %#v", taskID, s.taskCancels)
	}

	cancel()

	s.taskMu.Lock()
	delete(s.taskCancels, taskID)
	s.taskMu.Unlock()

	return nil
}
