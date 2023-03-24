package tasks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/pkg/logger"
)

type instagrapiClient interface {
	MakePost(ctx context.Context, task domain.Task, landingAccount string, sessionID string, cheapProxy string, targets []dbmodel.TargetUser, postImage []byte) error
	MakeReelsPost(ctx context.Context, task domain.Task, landingAccount string, sessionID string, cheapProxy string, targets []dbmodel.TargetUser, reelsVideo []byte) error
	InitBot(ctx context.Context, bot domain.BotWithTargets) error
	CheckLandingAccounts(ctx context.Context, sessionID string, landingAccountUsernames []string) ([]string, error)
	FollowTargets(ctx context.Context, bot domain.BotWithTargets) error
	EditProfile(ctx context.Context, fullName, sessionID string, image []byte) error
}

type worker struct {
	botsQueue      chan *domain.BotWithTargets
	dbtxf          dbmodel.DBTXFunc
	cli            instagrapiClient
	task           domain.Task
	generator      images.Generator
	processorIndex int64
	wg             *sync.WaitGroup
	videoBytes     []byte
}

type postCreator func(ctx context.Context, task domain.Task, landingAccount string, sessionID string, cheapProxy string, targets []dbmodel.TargetUser, postBytes []byte) error

func (w *worker) getPostCreator() postCreator {
	if w.task.Type == 2 {
		return w.cli.MakeReelsPost
	}

	return w.cli.MakePost
}

func (w worker) preparePostBytess(ctx context.Context) []byte {
	if w.task.Type == dbmodel.PhotoTaskType {
		return w.generator.Next(ctx)
	}

	return w.videoBytes
}
func (w *worker) run(ctx context.Context) {
	defer w.wg.Done()
	ctx = logger.WithKV(ctx, "processor_index", w.processorIndex)
	q := dbmodel.New(w.dbtxf(ctx))
	var err error

	if w.task.PerPostSleepSeconds == 0 {
		w.task.PerPostSleepSeconds = 10
	}

	sleepDuration := time.Duration(w.task.PerPostSleepSeconds) * time.Second

	for botWithTargets := range w.botsQueue {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "exiting from worker by context done")
			return
		default:
		}

		if botWithTargets == nil {
			logger.ErrorKV(ctx, "got nil bot with targets, skpping it")
			continue
		}

		startTime := time.Now()
		taskCtx := logger.WithKV(ctx, "bot_account", botWithTargets.Username)

		logger.DebugKV(taskCtx, "got account for processing")

		err = nil

		targetsLen := len(botWithTargets.Targets)
		if targetsLen != w.task.PostsPerBot*w.task.TargetsPerPost {
			logger.Warnf(taskCtx, "got %d targets, expected %d", targetsLen, w.task.PostsPerBot*w.task.TargetsPerPost)
		}

		err = w.cli.InitBot(taskCtx, *botWithTargets)
		if err != nil {
			logger.Errorf(taskCtx, "failed to init bot: %v", err)

			err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.FailBotStatus, ID: botWithTargets.ID})
			if err != nil {
				logger.Errorf(taskCtx, "failed to set bot status to 'failed': %v", err)
			}

			continue
		}

		if len(w.task.AccountProfileImages) > 0 {
			err = w.cli.EditProfile(
				taskCtx,
				"",
				botWithTargets.Headers.AuthData.SessionID,
				domain.RandomFromSlice(w.task.AccountProfileImages),
			)
			if err != nil {
				logger.Errorf(taskCtx, "failed to edit profile: %v", err)
			}
		}

		if w.task.FollowTargets {
			err = w.cli.FollowTargets(taskCtx, *botWithTargets)
			if err != nil {
				logger.Errorf(taskCtx, "failed to follow targets: %v", err)
			}
		}

		var cheapProxy string
		if botWithTargets.WorkProxy == nil {
			logger.Warnf(taskCtx, "bot has empty cheap proxy, so using residential for post upload")
			cheapProxy = botWithTargets.ResProxy.PythonString()
		} else {
			cheapProxy = botWithTargets.WorkProxy.PythonString()
		}

		err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.StartedBotStatus, ID: botWithTargets.ID})
		if err != nil {
			logger.Errorf(taskCtx, "failed to set bot status to 'started': %v", err)
			continue
		}

		landingAccount, err2 := w.chooseAliveLandingAccount(taskCtx, botWithTargets.BotAccount)
		if err2 != nil {
			logger.Errorf(taskCtx, "failed to select alive landing account: %v", err)

			break
		}

		var (
			i, postsDone int
			shouldBreak  = false
			targetIds    []uuid.UUID
		)

		if botWithTargets.PostsCount > 0 {
			logger.Infof(ctx, "bot already has %d posts, adding new", botWithTargets.PostsCount)
			postsDone += botWithTargets.PostsCount
		}

		for i = 0; i < w.task.PostsPerBot; i++ {
			rightBorderOfTargets := (i + 1) * w.task.TargetsPerPost
			if rightBorderOfTargets >= targetsLen {
				rightBorderOfTargets = targetsLen - 1
				shouldBreak = true
			}

			targetsBatch := botWithTargets.Targets[i*w.task.TargetsPerPost : rightBorderOfTargets]
			targetIds = domain.Ids(targetsBatch)

			err = w.getPostCreator()(taskCtx, w.task, landingAccount, botWithTargets.Headers.AuthData.SessionID, cheapProxy, targetsBatch, w.preparePostBytess(ctx))
			if err != nil {
				logger.Errorf(taskCtx, "failed to create post [%d]: %v", i, err)
				err = q.SetTargetsStatus(taskCtx, dbmodel.SetTargetsStatusParams{Status: dbmodel.TargetsStatusFailed, Ids: targetIds})
				if err != nil {
					logger.Errorf(taskCtx, "failed to set targets statuses to 'failed' for targets '%v': %v", targetIds, err)
				}

				err = q.SetTargetsStatus(taskCtx, dbmodel.SetTargetsStatusParams{Status: dbmodel.TargetsStatusFailed, Ids: targetIds})
				if err != nil {
					logger.Errorf(taskCtx, "failed to set targets statuses to 'notified' for targets '%v': %v", targetIds, err)
				}

				if errors.Is(err, instagrapi.ErrBotIsBlocked) {
					break
				}

				continue
			}

			postsDone++

			err = q.SetBotPostsCount(taskCtx, dbmodel.SetBotPostsCountParams{PostsCount: postsDone, ID: botWithTargets.ID})
			if err != nil {
				logger.Errorf(taskCtx, "failed to update posts count: %v", err)
			}

			err = q.SetTargetsStatus(taskCtx, dbmodel.SetTargetsStatusParams{Status: dbmodel.TargetsStatusNotified, Ids: targetIds})
			if err != nil {
				logger.Errorf(taskCtx, "failed to set targets statuses to 'notified' for targets '%v': %v", targetIds, err)
			}

			// тегнули уже всех пользователей, больше постов не нужно
			if shouldBreak {
				break
			}

			select {
			case <-ctx.Done():
				logger.Warnf(taskCtx, "exiting from worker by context done, created %d posts", i)
				return
			case <-time.After(sleepDuration):
				logger.Infof(ctx, "slept for %s", sleepDuration)
			}
		}

		logger.Infof(taskCtx, "made %d posts, saving results time elapsed: %s", postsDone, time.Since(startTime))

		err = q.SetBotDoneStatus(taskCtx, botWithTargets.ID)
		if err != nil {
			logger.Errorf(taskCtx, "failed to mark bot account as completed: %v", err)
		}
	}

	logger.Infof(ctx, "bots queue closed, stopping worker")
}

type APIResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	ErrorType     string `json:"error_type"`
	ExceptionName string `json:"exception_name"`
}

func (w *worker) chooseAliveLandingAccount(ctx context.Context, bot domain.BotAccount) (string, error) {
	if len(w.task.LandingAccounts) == 0 {
		return "", fmt.Errorf("empty list of landing accounts")
	}

	aliveLandingAccounts, err := w.cli.CheckLandingAccounts(ctx, bot.Headers.AuthData.SessionID, w.task.LandingAccounts)
	if err != nil {
		return "", err
	}

	if len(aliveLandingAccounts) == 0 {
		return "", fmt.Errorf("all landing accounts are dead")
	}

	w.task.LandingAccounts = aliveLandingAccounts

	return domain.RandomFromSlice(aliveLandingAccounts), nil
}

// '{"upload_id":"1664888837874","xsharing_nonces":{},"status":"ok"}' webp
// got '{"upload_id":"1664889100793","xsharing_nonces":{},"status":"ok"}' from body jpeg
