package tasks

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/pkg/logger"
)

type instagrapiClient interface {
	MakePost(ctx context.Context, sessionID, caption string, image []byte) error
	InitBot(ctx context.Context, bot domain.BotWithTargets) error
}

type worker struct {
	botsQueue      chan *domain.BotWithTargets
	dbtxf          dbmodel.DBTXFunc
	cli            instagrapiClient
	task           domain.Task
	generator      images.Generator
	processorIndex int64
	captionFormat  string
}

const (
	postsPerBot    = 7
	targetsPerPost = 15
)

func (w *worker) run(ctx context.Context) {
	ctx = logger.WithKV(ctx, "processor_index", w.processorIndex)
	q := dbmodel.New(w.dbtxf(ctx))
	var err error
	for botWithTargets := range w.botsQueue {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "exiting from worker by context done")
			return
		default:
		}

		if botWithTargets == nil {
			logger.Error(ctx, "got nil bot with targets, skpping it")
			continue
		}

		startTime := time.Now()
		taskCtx := logger.WithKV(ctx, "bot_account", botWithTargets.Username)

		logger.Debug(taskCtx, "got account for processing")

		err = nil

		targetsLen := len(botWithTargets.Targets)
		if targetsLen != postsPerBot*targetsPerPost {
			logger.Warnf(taskCtx, "got %d targets, expected %d", targetsLen, postsPerBot*targetsPerPost)
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

		err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.StartedBotStatus, ID: botWithTargets.ID})
		if err != nil {
			logger.Errorf(taskCtx, "failed to set bot status to 'started': %v", err)
			continue
		}

		var (
			i, postsDone int
			shouldBreak  = false
			targetIds    []uuid.UUID
		)

		for i = 0; i < postsPerBot; i++ {
			rightBorderOfTargets := (i + 1) * targetsPerPost
			if rightBorderOfTargets >= targetsLen {
				rightBorderOfTargets = targetsLen - 1
				shouldBreak = true
			}

			targetsBatch := botWithTargets.Targets[i*targetsPerPost : rightBorderOfTargets]
			targetIds = domain.Ids(targetsBatch)
			caption := w.preparePostCaption(w.task.TextTemplate, targetsBatch)

			err = w.cli.MakePost(taskCtx, botWithTargets.Headers.AuthData.SessionID, caption, w.generator.Next(taskCtx))
			if err != nil {
				logger.Errorf(taskCtx, "failed to create post [%d]: %v", i, err)
				err = q.SetTargetsStatus(taskCtx, dbmodel.SetTargetsStatusParams{Status: dbmodel.FailedTargetStatus, Ids: targetIds})
				if err != nil {
					logger.Errorf(taskCtx, "failed to set targets statuses to 'failed' for targets '%v': %v", targetIds, err)
				}

				continue
			}

			postsDone++

			err = q.SetTargetsStatus(taskCtx, dbmodel.SetTargetsStatusParams{Status: dbmodel.NotifiedTargetStatus, Ids: targetIds})
			if err != nil {
				logger.Errorf(taskCtx, "failed to set targets statuses to 'notified' for targets '%v': %v", targetIds, err)
				continue
			}

			// тегнули уже всех пользователей, больше постов не нужно
			if shouldBreak {
				break
			}

			select {
			case <-ctx.Done():
				logger.Warnf(taskCtx, "exiting from worker by context done, created %d posts", i)
				err = q.SetBotPostsCount(taskCtx, dbmodel.SetBotPostsCountParams{PostsCount: int16(postsDone), ID: botWithTargets.ID})
				if err != nil {
					logger.Errorf(taskCtx, "failed to mark bot account as completed: %v", err)
				}
				return
			case <-time.After(3 * time.Second):
			}
		}

		logger.Infof(taskCtx, "made %d posts, saving results time elapsed: %s", postsDone, time.Since(startTime))

		err = q.SetBotPostsCount(taskCtx, dbmodel.SetBotPostsCountParams{PostsCount: int16(postsDone), ID: botWithTargets.ID})
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

func (w *worker) preparePostCaption(template string, targetUsers []dbmodel.TargetUser) string {
	b := strings.Builder{}
	b.WriteString(template)

	for _, user := range targetUsers {
		b.WriteByte('@')
		b.WriteString(user.Username)
		b.WriteByte(' ')
	}

	return b.String()
}

// '{"upload_id":"1664888837874","xsharing_nonces":{},"status":"ok"}' webp
// got '{"upload_id":"1664889100793","xsharing_nonces":{},"status":"ok"}' from body jpeg
