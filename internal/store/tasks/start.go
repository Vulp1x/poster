package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

const landingAccountPlaceholder = "@account"

// ErrTaskInvalidTextTemplate у задачи нет фотографий для постов
var ErrTaskInvalidTextTemplate = errors.New("task doesn't have landing account placeholder in text tamplate")

// ErrTaskWithEmptyPostImages у задачи нет фотографий для постов
var ErrTaskWithEmptyPostImages = errors.New("task doesn't have post images")

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

	if task.Status != dbmodel.ReadyTaskStatus {
		return nil, fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.ReadyTaskStatus, task.Status)
	}

	if len(task.Images) == 0 {
		return nil, ErrTaskWithEmptyPostImages
	}

	if !strings.Contains(task.TextTemplate, landingAccountPlaceholder) {
		return nil, ErrTaskInvalidTextTemplate
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
		return nil, fmt.Errorf("no bots for task found: %v", err)
	}

	maximumTargetsNum := int32(len(bots) * postsPerBot * targetsPerPost)
	targets, err := q.FindUnprocessedTargetsForTask(ctx, dbmodel.FindUnprocessedTargetsForTaskParams{
		TaskID: taskID, Limit: maximumTargetsNum,
	})
	if err != nil || len(targets) == 0 {
		return nil, fmt.Errorf("failed to find targets for task: %v", err)
	}

	neededBotsNum := len(targets) / (postsPerBot * targetsPerPost)

	logger.Infof(ctx, "going to use %d/%d bots for %d targets", neededBotsNum, len(bots), len(targets))

	if neededBotsNum < len(bots) {
		bots = bots[:neededBotsNum]
	}

	randomBot := domain.RandomFromSlice(bots)

	aliveLandings, err := s.checkAliveLandingAccounts(ctx, randomBot, task.LandingAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to check landing accounts with bot '%s': %v", randomBot.Username, err)
	}

	logger.Infof(ctx, "got alive landing accounts: %v", aliveLandings)
	task.LandingAccounts = aliveLandings

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

	for i := 0; i < workersPerTask; i++ {
		postingWorker := &worker{
			botsQueue:      botsChan,
			dbtxf:          s.dbtxf,
			cli:            s.instaClient,
			task:           domain.Task(task),
			generator:      imageGenerator,
			processorIndex: int64(i),
		}

		go postingWorker.run(taskCtx)
	}

	go s.asyncPushBots(taskCtx, botsChan, bots, targets)

	return aliveLandings, nil
}

func (s *Store) checkAliveLandingAccounts(ctx context.Context, bot dbmodel.BotAccount, landingAccounts []string) ([]string, error) {
	err := s.instaClient.InitBot(ctx, domain.BotWithTargets{BotAccount: domain.BotAccount(bot)})
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %v", err)
	}

	return s.instaClient.CheckLandingAccounts(ctx, bot.Headers.AuthData.SessionID, landingAccounts)
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
