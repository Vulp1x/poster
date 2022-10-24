package tasks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

const workersPerTask = 1

// ErrTaskNotFound не смогли найти таску
var ErrTaskNotFound = errors.New("task not found")

// ErrTaskInvalidStatus переход по статусам не возможен
var ErrTaskInvalidStatus = errors.New("invalid task status")

func NewStore(timeout time.Duration, dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc) *Store {
	return &Store{
		tasksChan:   make(chan domain.Task, 10),
		taskCancels: make(map[uuid.UUID]func()),
		pushTimeout: timeout,
		dbtxf:       dbtxFunc,
		txf:         txFunc,
		taskMu:      &sync.Mutex{},
		instaClient: instagrapi.NewClient(),
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
}

func (s *Store) TaskProgress(ctx context.Context, taskID uuid.UUID) (domain.TaskProgress, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TaskProgress{}, ErrTaskNotFound
		}
	}

	if task.Status != dbmodel.StartedTaskStatus && task.Status != dbmodel.DoneTaskStatus {
		return domain.TaskProgress{}, fmt.Errorf("%w: expected statuses [%d, %d], got %d",
			ErrTaskInvalidStatus, dbmodel.StartedTaskStatus, dbmodel.DoneTaskStatus, task.Status)
	}

	progress, err := q.GetBotsProgress(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, err
	}

	targetCounters, err := q.GetTaskTargetsCount(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, fmt.Errorf("failed to get target counters: %v", err)
	}

	return domain.TaskProgress{BotsProgress: progress, TargetCounters: targetCounters}, nil
}

func (s *Store) UpdateTask(ctx context.Context, taskID uuid.UUID, opts ...UpdateOption) (domain.Task, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return domain.Task{}, store.ErrTransactionFail
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if task.Status == dbmodel.StartedTaskStatus {
		return domain.Task{}, fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.DataUploadedTaskStatus, task.Status)
	}

	for _, opt := range opts {
		opt(&task)
	}

	err = q.UpdateTask(ctx, dbmodel.UpdateTaskParams{
		TextTemplate:         task.TextTemplate,
		Title:                task.Title,
		Images:               task.Images,
		AccountNames:         task.AccountNames,
		AccountLastNames:     task.AccountLastNames,
		AccountUrls:          task.AccountUrls,
		AccountProfileImages: task.AccountProfileImages,
		LandingAccounts:      task.LandingAccounts,
		ID:                   taskID,
	})
	if err != nil {
		return domain.Task{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Task{}, err
	}

	return domain.Task(task), nil
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
		return err
	}

	savedCount, err = q.SaveProxies(ctx, residentialProxies.ToSaveParams(taskID, false))
	logger.Infof(ctx, "saved %d residential proxies", savedCount)
	if err != nil {
		return err
	}

	savedCount, err = q.SaveProxies(ctx, cheapProxies.ToSaveParams(taskID, true))
	logger.Infof(ctx, "saved %d cheap proxies", savedCount)
	if err != nil {
		return err
	}

	savedCount, err = q.SaveTargetUsers(ctx, targets.ToSaveParams(taskID))
	logger.Infof(ctx, "saved %d target users", savedCount)
	if err != nil {
		return err
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

// UpdateOption позволяет добавить опциональные поля для создания драфтовой задачи
type UpdateOption func(params *dbmodel.Task)

func WithBotNamesUpdateOption(names []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(names) != 0 {
			task.AccountNames = names
		}
	}
}

func WithBotLasNamesUpdateOption(lastNames []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(lastNames) != 0 {
			task.AccountLastNames = lastNames
		}
	}
}

func WithBotURLsUpdateOption(urls []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(urls) != 0 {
			task.AccountUrls = urls
		}
	}
}

func WithBotImagesUpdateOption(images [][]byte) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(images) != 0 {
			task.AccountProfileImages = images
		}
	}
}

func WithImagesUpdateOption(images [][]byte) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(images) != 0 {
			task.Images = images
		}
	}
}

func WithLandingAccountsUpdateOption(landingAccounts []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(landingAccounts) != 0 {
			task.LandingAccounts = landingAccounts
		}
	}
}

func WithTextTemplateUpdateOption(template *string) UpdateOption {
	return func(task *dbmodel.Task) {
		if template != nil {
			task.TextTemplate = *template
		}
	}
}

func WithTitleUpdateOption(title *string) UpdateOption {
	return func(task *dbmodel.Task) {
		if title != nil {
			task.Title = *title
		}
	}
}
