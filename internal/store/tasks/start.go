package tasks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

const landingAccountPlaceholder = "@account"

// ErrTaskInvalidTextTemplate у задачи нет фотографий для постов
var ErrTaskInvalidTextTemplate = errors.New("task doesn't have landing account placeholder in text tamplate")

// ErrTaskWithEmptyPostImages у задачи нет фотографий для постов
var ErrTaskWithEmptyPostImages = errors.New("task doesn't have post images")

// ErrTaskWithEmptyTargetsPerPost у задачи нет фотографий для постов
var ErrTaskWithEmptyTargetsPerPost = errors.New("task have 0 targets per post")

// ErrTaskWithEmptyPostsPerBot у задачи нет фотографий для постов
var ErrTaskWithEmptyPostsPerBot = errors.New("task have 0 posts per bot")

// ErrTaskWithEmptyLandingAccounts у задачи нет фотографий для постов
var ErrTaskWithEmptyLandingAccounts = errors.New("task have 0 landing accounts")

// StartTask начинает выполнение задачи
func (s *Store) StartTask(ctx context.Context, taskID uuid.UUID) ([]string, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		return nil, fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if err = validateTaskBeforeStart(task); err != nil {
		return nil, err
	}

	imageGenerator, err := images.NewRandomGammaGenerator(task.Images)
	if err != nil {
		return nil, err
	}

	bots, err := q.FindReadyBotsForTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find bots for task: %v", err)
	}

	if len(bots) == 0 {
		err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{Status: dbmodel.DoneTaskStatus, ID: task.ID})
		if err != nil {
			return nil, fmt.Errorf("failed to set task status to done: %v", err)
		}

		return nil, fmt.Errorf(" в задаче нет ботов, готовых к работе")
	}

	maximumTargetsNum := int32(len(bots)) * task.PostsPerBot * task.TargetsPerPost
	targets, err := q.FindUnprocessedTargetsForTask(ctx, dbmodel.FindUnprocessedTargetsForTaskParams{
		TaskID: taskID, Limit: maximumTargetsNum,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find targets for task: %v", err)
	}

	if len(targets) == 0 {
		err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{Status: dbmodel.DoneTaskStatus, ID: task.ID})
		if err != nil {
			return nil, fmt.Errorf("failed to set task status to done: %v", err)
		}

		return nil, fmt.Errorf("в задаче нет неоповещенных целей")
	}

	neededBotsNum := int32(len(targets)) / (task.PostsPerBot * task.TargetsPerPost)

	logger.Infof(ctx, "going to use %d/%d bots for %d targets", neededBotsNum, len(bots), len(targets))

	if neededBotsNum < int32(len(bots)) {
		bots = bots[:neededBotsNum]
	}

	randomBot := domain.RandomFromSlice(bots)

	aliveLandings, err := s.checkAliveLandingAccounts(ctx, randomBot, task.LandingAccounts)
	if err != nil {
		if errors.Is(err, instagrapi.ErrBotIsBlocked) {
			err2 := q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.BlockedBotStatus, ID: randomBot.ID})
			if err2 != nil {
				logger.Errorf(ctx, "failed to set bot status to blocked: %v", err)
			}

			return nil, fmt.Errorf("аккаунт '%s' для проверки лендингов был заблокирован, попробуйте ещё раз", randomBot.Username)
		}

		return nil, fmt.Errorf("failed to check landing accounts with bot '%s': %v", randomBot.Username, err)
	}

	logger.Infof(ctx, "got alive landing accounts: %v", aliveLandings)
	task.LandingAccounts = aliveLandings

	var videoBytes []byte

	if task.Type == dbmodel.ReelsTaskType {
		if !strings.HasSuffix(strings.ToLower(*task.VideoFilename), "mp4") {
			return nil, fmt.Errorf("only mp4 files are supported, got: %s", *task.VideoFilename)
		}

		f, err := os.Open(*task.VideoFilename)
		if err != nil {
			return nil, fmt.Errorf("failed to open video file at '%s': %v", *task.VideoFilename, err)
		}

		videoBytes, err = io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read video: %v", err)
		}

		err = f.Close()
		if err != nil {
			logger.Errorf(ctx, "failed to close file at '%s': %v", *task.VideoFilename, err)
		}

	}

	err = q.StartTaskByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to start task: %v", err)
	}

	// нужно отвязаться от ctx, так как он закенселится сразу после окончания запроса
	taskCtx, taskCancel := context.WithCancel(logger.ToContext(context.Background(), logger.FromContext(ctx)))
	s.taskMu.Lock()
	s.taskCancels[task.ID] = taskCancel
	s.taskMu.Unlock()

	botsChan := make(chan *domain.BotWithTargets, 20)

	wg := &sync.WaitGroup{}

	for i := 0; i < workersPerTask; i++ {
		postingWorker := &worker{
			botsQueue:      botsChan,
			dbtxf:          s.dbtxf,
			cli:            s.instaClient,
			task:           domain.Task(task),
			generator:      imageGenerator,
			processorIndex: int64(i),
			wg:             wg,
			videoBytes:     videoBytes,
		}

		go postingWorker.run(taskCtx)
	}

	wg.Add(workersPerTask)

	go s.asyncPushBots(taskCtx, q, task, botsChan, bots, targets, wg)

	return aliveLandings, nil
}

func validateTaskBeforeStart(task dbmodel.Task) error {
	if task.Status != dbmodel.ReadyTaskStatus {
		return fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.ReadyTaskStatus, task.Status)
	}

	if len(task.Images) == 0 && task.Type == dbmodel.PhotoTaskType {
		return ErrTaskWithEmptyPostImages
	}

	if task.VideoFilename == nil && task.Type == dbmodel.ReelsTaskType {
		return ErrTaskWithEmptyPostImages
	}

	if !strings.Contains(task.TextTemplate, landingAccountPlaceholder) {
		return ErrTaskInvalidTextTemplate
	}

	if task.TargetsPerPost == 0 {
		return ErrTaskWithEmptyTargetsPerPost
	}

	if task.PostsPerBot == 0 {
		return ErrTaskWithEmptyPostsPerBot
	}

	if len(task.LandingAccounts) == 0 {
		return ErrTaskWithEmptyLandingAccounts
	}

	return nil
}

func (s *Store) checkAliveLandingAccounts(ctx context.Context, bot dbmodel.BotAccount, landingAccounts []string) ([]string, error) {
	err := s.instaClient.InitBot(ctx, domain.BotWithTargets{BotAccount: domain.BotAccount(bot)})
	if err != nil {
		return nil, fmt.Errorf("failed to init bot when checking alive accounts: %w", err)
	}

	return s.instaClient.CheckLandingAccounts(ctx, bot.Headers.AuthData.SessionID, landingAccounts)
}

func (s *Store) asyncPushBots(ctx context.Context, q *dbmodel.Queries, task dbmodel.Task, botsChan chan *domain.BotWithTargets, bots []dbmodel.BotAccount, targets []dbmodel.TargetUser, wg *sync.WaitGroup) {
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

		targetsPerPost := int(task.PostsPerBot * task.TargetsPerPost)

		batchEnd = (i + 1) * targetsPerPost
		if batchEnd > len(targets) {
			batchEnd = len(targets) - 1
			allTargetsProcessed = true
		}

		botWithTargets := &domain.BotWithTargets{
			BotAccount: domain.BotAccount(bot),
			Targets:    targets[i*targetsPerPost : batchEnd],
		}
		botsChan <- botWithTargets

		if allTargetsProcessed && i != len(bots)-1 {
			logger.Warnf(ctx, "processed %d targets with %d batches, breaking", len(targets), i+1)
			break
		}
	}

	logger.Infof(ctx, "pushed all bots in %s: closing chan", time.Since(startedAt))
	close(botsChan)

	wg.Wait()

	logger.Infof(ctx, "all workers done in %s: setting task status to done ", time.Since(startedAt))
	err := q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{
		Status: dbmodel.DoneTaskStatus,
		ID:     task.ID,
	})
	if err != nil {
		logger.Errorf(ctx, "failed to set task status to done: %v", err)
	}
}
