package tasks

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

const workersPerTask = 1

// ErrTaskNotFound не смогли найти таску
var ErrTaskNotFound = errors.New("task not found")

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

	progress, err := q.GetTaskProgress(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, err
	}

	return progress, nil
}

func (s *Store) UpdateTask(ctx context.Context, taskID uuid.UUID, title, textTemplate *string, image [][]byte) (domain.Task, error) {
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

	if title != nil {
		task.Title = *title
	}

	if textTemplate != nil {
		task.TextTemplate = *textTemplate
	}

	if image != nil {
		task.Images = image
	}

	err = q.UpdateTask(ctx, dbmodel.UpdateTaskParams{
		TextTemplate: task.TextTemplate,
		Title:        task.Title,
		Images:       task.Images,
		ID:           taskID,
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

func (s *Store) AssignProxies(ctx context.Context, taskID uuid.UUID) (int, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return 0, store.ErrTransactionFail
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrTaskNotFound
		}

		return 0, fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if task.Status != dbmodel.DataUploadedTaskStatus {
		return 0, fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.DataUploadedTaskStatus, task.Status)
	}

	botAccounts, err := q.FindBotsForTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("failed to find bot accounts for task: %v", err)
	}

	proxies, err := q.FindProxiesForTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("failed to find proxiesIds for task: %v", err)
	}

	// after deleting botAccounts and proxies would have same length
	botAccounts, proxies, err = s.deleteUnnecessaryRows(ctx, tx, botAccounts, proxies)
	if err != nil {
		return 0, err
	}

	botIds := domain.Ids(botAccounts)
	err = q.AssignProxiesToBotsForTask(ctx, dbmodel.AssignProxiesToBotsForTaskParams{
		TaskID:  taskID,
		Proxies: domain.Strings(proxies),
		Ids:     botIds,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to set bot accounts proxy: %v", err)
	}

	err = q.AssignBotsToProxiesForTask(ctx, dbmodel.AssignBotsToProxiesForTaskParams{
		TaskID: taskID,
		Ids:    domain.Ids(proxies),
		BotIds: botIds,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to set bot accounts proxy: %v", err)
	}

	err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{Status: dbmodel.ReadyTaskStatus, ID: taskID})
	if err != nil {
		return 0, fmt.Errorf("failed to update task status: %v", err)
	}

	return len(botIds), tx.Commit(ctx)
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

func (s *Store) CreateDraftTask(ctx context.Context, userID uuid.UUID, title, textTemplate string, accounts []string, images [][]byte) (uuid.UUID, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	taskID, err := q.CreateDraftTask(ctx, dbmodel.CreateDraftTaskParams{
		ManagerID:       userID,
		TextTemplate:    textTemplate,
		LandingAccounts: accounts,
		Images:          images,
		Title:           title,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create task draft: %w", err)
	}

	return taskID, nil
}

func (s *Store) StartTask(ctx context.Context, taskID uuid.UUID) error {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}

		return fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if task.Status != dbmodel.ReadyTaskStatus {
		return fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.ReadyTaskStatus, task.Status)
	}

	imageGenerator, err := images.NewRandomGammaGenerator(task.Images)
	if err != nil {
		return err
	}

	bots, err := q.FindReadyBotsForTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to find bots for task: %v", err)
	}

	if len(bots) == 0 {
		return fmt.Errorf("no bots for task found: %v", err)
	}

	maximumTargetsNum := int32(len(bots) * postsPerBot * targetsPerPost)
	targets, err := q.FindUnprocessedTargetsForTask(ctx, dbmodel.FindUnprocessedTargetsForTaskParams{
		TaskID: taskID, Limit: maximumTargetsNum,
	})
	if err != nil {
		return fmt.Errorf("failed to find targets for task: %v", err)
	}

	neededBotsNum := len(targets) / (postsPerBot * targetsPerPost)

	logger.Infof(ctx, "going to use %d/%d bots for %d targets", neededBotsNum, len(bots), len(targets))

	if neededBotsNum < len(bots) {
		bots = bots[:neededBotsNum]
	}

	randomBot := bots[rand.Intn(len(bots)-1)]

	aliveLandings, err := s.checkAliveLandingAccounts(ctx, randomBot, task.LandingAccounts)
	if err != nil {
		return fmt.Errorf("failed to check landing accounts with bot '%s': %v", randomBot.Username, err)
	}

	logger.Infof(ctx, "got alive landing accounts: %v", aliveLandings)
	task.LandingAccounts = aliveLandings

	err = q.StartTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to start task: %v", err)
	}

	// нужно отвязаться от ctx, так как он закенселится сразу после окончания запроса
	taskCtx, taskCancel := context.WithCancel(logger.ToContext(context.Background(), logger.FromContext(ctx)))
	s.taskMu.Lock()
	s.taskCancels[task.ID] = taskCancel
	s.taskMu.Unlock()

	botsChan := make(chan *domain.BotWithTargets, 20)

	for i := 0; i < workersPerTask; i++ {
		postingWorker := &worker{
			botsQueue:      botsChan,
			dbtxf:          s.dbtxf,
			cli:            instagrapi.NewClient(),
			task:           domain.Task(task),
			generator:      imageGenerator,
			processorIndex: int64(i),
		}

		go postingWorker.run(taskCtx)
	}

	go s.asyncPushBots(taskCtx, botsChan, bots, targets)

	return nil
}

func (s *Store) checkAliveLandingAccounts(ctx context.Context, bot dbmodel.BotAccount, landingAccounts []string) ([]string, error) {
	err := s.instaClient.InitBot(ctx, domain.BotWithTargets{BotAccount: domain.BotAccount(bot)})
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %v", err)
	}

	return s.instaClient.CheckLandingAccounts(ctx, bot.Headers.AuthData.SessionID, landingAccounts)
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

func (s *Store) deleteUnnecessaryRows(ctx context.Context, tx dbmodel.Tx, accounts []dbmodel.BotAccount, proxies []dbmodel.Proxy) ([]dbmodel.BotAccount, []dbmodel.Proxy, error) {
	q := dbmodel.New(tx)

	accountsLen, proxiesLen := len(accounts), len(proxies)
	logger.Infof(ctx, "got %d accounts and %d proxies", accountsLen, proxiesLen)

	var remainRows = min(accountsLen, proxiesLen)

	switch {
	case accountsLen < proxiesLen:
		// надо удалить лишние прокси из задачи

		rowsToDelete := accountsLen - proxiesLen
		deletedRowsCount, err := q.DeleteProxiesForTask(ctx, proxiesLastIds(proxies, rowsToDelete))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to delete proxies: %v", err)
		}

		if int(deletedRowsCount) != rowsToDelete {
			return nil, nil, fmt.Errorf("wanted to delete %d proxies, deleted %d", rowsToDelete, deletedRowsCount)
		}

	case accountsLen == proxiesLen:
		return accounts, proxies, nil

	case accountsLen > proxiesLen:
		// надо удалить лишних ботов из задачи

		rowsToDelete := accountsLen - proxiesLen
		deletedRowsCount, err := q.DeleteBotAccountsForTask(ctx, accountsLastIds(accounts, rowsToDelete))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to delete bot accounts: %v", err)
		}

		if int(deletedRowsCount) != rowsToDelete {
			return nil, nil, fmt.Errorf("wanted to delete %d bot accounts, deleted %d", rowsToDelete, deletedRowsCount)
		}
	}

	return accounts[:remainRows], proxies[:remainRows], nil
}

func (s *Store) asyncPushBots(ctx context.Context, botsChan chan *domain.BotWithTargets, bots []dbmodel.BotAccount, targets []dbmodel.TargetUser) {
	startedAt := time.Now()

	var batchEnd int
	var allTargetsProcessed bool

	for i, bot := range bots {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "stopping bot pushes, pushed %d/%d bots: context done", i, len(bots))
			return
		default:
		}

		batchEnd = (i + 1) * postsPerBot * targetsPerPost
		if batchEnd > len(targets) {
			batchEnd = len(targets) - 1
			allTargetsProcessed = true
		}

		botWithTargets := &domain.BotWithTargets{
			BotAccount: domain.BotAccount(bot),
			Targets:    targets[i*postsPerBot*targetsPerPost : batchEnd],
		}
		botsChan <- botWithTargets

		if allTargetsProcessed && i != len(bots)-1 {
			logger.Warnf(ctx, "processed %d targets with %d batches, breaking", len(targets), i+1)
			break
		}
	}

	logger.Infof(ctx, "pushed all bots in %s: closing chan", time.Since(startedAt))
	close(botsChan)
}

// accountsLastIds возвращает список из rowsToDelete последних айдишников
func accountsLastIds(arr []dbmodel.BotAccount, rowsToDelete int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, rowsToDelete)
	for _, account := range arr[len(arr)-rowsToDelete:] {
		ids = append(ids, account.ID)
	}

	return ids
}

// proxiesLastIds возвращает список из rowsToDelete последних айдишников
func proxiesLastIds(arr []dbmodel.Proxy, rowsToDelete int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, rowsToDelete)
	for _, account := range arr[len(arr)-rowsToDelete:] {
		ids = append(ids, account.ID)
	}

	return ids
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
