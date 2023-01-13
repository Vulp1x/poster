package tasks

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/instagrapi"
	api "github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/internal/workers"
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

func (s *Store) StartUpdatePostContents(ctx context.Context, taskID uuid.UUID) ([]string, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, err
	}

	if task.Status != dbmodel.DoneTaskStatus {
		return nil, fmt.Errorf("got status %d, expected 6(done)", task.Status)
	}

	// TODO добавить проверку landing accounts перед началом
	// s.cli.CheckLandingAccounts()

	if err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{
		Status: dbmodel.UpdatingPostContentsTaskStatus,
		ID:     taskID,
	}); err != nil {
		return nil, fmt.Errorf("failed to update task status: %v", err)
	}

	bots, err := q.FindBotsForTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bots for task: %v", err)
	}

	tasks := make([]pgqueue.Task, len(bots))
	for i, bot := range bots {
		tasks[i] = pgqueue.Task{
			Kind:        workers.UpdatePostsContentsTaskKind,
			Payload:     workers.EmptyPayload,
			ExternalKey: fmt.Sprintf("%s::%s", task.ID.String(), bot.Username),
		}
	}

	if err = s.queue.PushTasksTx(ctx, tx, tasks); err != nil {
		return nil, fmt.Errorf("failed to push %d tasks: %v", len(tasks), err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return task.LandingAccounts, nil
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

func (s *Store) insertMoreProxies(
	ctx context.Context,
	taskID uuid.UUID,
	tx dbmodel.Tx,
	initialProxies []dbmodel.Proxy,
	proxiesToInsert int,
	isCheap bool,
) ([]dbmodel.Proxy, error) {
	var newProxies = make([]dbmodel.Proxy, 0, proxiesToInsert)
	var proxiesFromOneInitialProxy = math.Ceil(float64(proxiesToInsert) / float64(len(initialProxies)))

	for _, proxy := range initialProxies {
		for i := 0; i < int(proxiesFromOneInitialProxy); i++ {
			newProxies = append(newProxies, proxy)
		}
	}

	logger.Infof(ctx, "got %d new proxies, wanted at least %d,from %d initial",
		len(newProxies), proxiesToInsert, len(initialProxies),
	)

	q := dbmodel.New(tx)
	savedProxiesCount, err := q.SaveProxies(ctx, domain.ProxiesFromDB(newProxies).ToSaveParams(taskID, isCheap))
	if err != nil {
		return nil, fmt.Errorf("failed to save new proxies: %v", err)
	}

	logger.Infof(ctx, "saved %d new proxies, is cheap: %t", savedProxiesCount, isCheap)

	return newProxies, nil
}
