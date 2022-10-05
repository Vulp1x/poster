package tasks

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/pkg/logger"
)

type instagrapiClient interface {
	MakePost(ctx context.Context, sessionID, caption string, image []byte) error
	InitBot(ctx context.Context, bot domain.BotAccount) error
}

type worker struct {
	tasksQueue     chan *domain.BotWithTargets
	dbtxf          dbmodel.DBTXFunc
	cli            instagrapiClient
	task           domain.Task
	generator      images.Generator
	processorIndex int64
	captionFormat  string
}

const (
	postsPerBot    = 15
	targetsPerPost = 15
)

func (w *worker) run(ctx context.Context) {
	ctx = logger.WithKV(ctx, "processor_index", w.processorIndex)
	q := dbmodel.New(w.dbtxf(ctx))
	var err error
	for botWithTargets := range w.tasksQueue {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "exiting from worker by context done")
			return
		default:
		}

		startTime := time.Now()
		taskCtx := logger.WithKV(ctx, "bot_account", botWithTargets.Username)

		logger.Debug(taskCtx, "got account for processing")

		err = nil

		targetsLen := len(botWithTargets.Targets)
		if targetsLen != postsPerBot*targetsPerPost {
			logger.Warnf(taskCtx, "got %d targets, expected %d", targetsLen, postsPerBot*targetsPerPost)
		}

		var i int
		var shouldBreak = false
		for i = 0; i < postsPerBot; i++ {
			rightBorderOfTargets := (i + 1) * targetsPerPost
			if rightBorderOfTargets >= targetsLen {
				rightBorderOfTargets = targetsLen - 1
				shouldBreak = true
			}

			targetsBatch := botWithTargets.Targets[i*targetsPerPost : rightBorderOfTargets]
			caption := w.preparePostCaption(w.task.TextTemplate, targetsBatch)
			err = w.cli.MakePost(taskCtx, botWithTargets.Headers.AuthData.SessionID, caption, w.generator.Next(taskCtx))
			if err != nil {
				logger.Errorf(taskCtx, "failed to create post [%d]: %v", i, err)
				break
			}

			// тегнули уже всех пользователей, больше постов не нужно
			if shouldBreak {
				break
			}

			select {
			case <-ctx.Done():
				logger.Warnf(ctx, "exiting from worker by context done, created %d posts", i)
				return
			case <-time.After(3 * time.Second):
			}
		}

		logger.Info(taskCtx, "make %d posts, saving results time elapsed: %s", i, time.Since(startTime))

		err = q.MarkBotAsCompleted(taskCtx, botWithTargets.ID)
		if err != nil {
			logger.Errorf(taskCtx, "failed to mark bot account as completed: %v", err)
		}

		err = q.MarkTargetsAsNotified(taskCtx, w.task.ID)
		if err != nil {
			logger.Errorf(taskCtx, "failed to mark targets as completed: %v", err)
		}
	}
}

type APIResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	ErrorType     string `json:"error_type"`
	ExceptionName string `json:"exception_name"`
}

func (w *worker) saveRequest(ctx context.Context, reqBody []byte, resp *http.Response) error {

	return nil
}

func (w *worker) createPost(ctx context.Context, botAccount domain.BotAccount, targetUsers []dbmodel.TargetUser) error {

	return nil
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
