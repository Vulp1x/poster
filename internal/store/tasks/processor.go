package tasks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/internal/requests"
	"github.com/inst-api/poster/internal/transport"
	"github.com/inst-api/poster/pkg/logger"
)

type worker struct {
	tasksQueue     chan *domain.TaskPerBot
	dbtxf          dbmodel.DBTXFunc
	cli            *http.Client
	generator      images.Generator
	processorIndex int64
	captionFormat  string
}

const postsPerBot = 15

func (p *worker) run(ctx context.Context) {
	ctx = logger.WithKV(ctx, "processor_index", p.processorIndex)
	q := dbmodel.New(p.dbtxf(ctx))
	var err error
	for task := range p.tasksQueue {
		select {
		case <-ctx.Done():
			logger.Infof(ctx, "exiting from worker by context done")
			return
		default:
		}

		startTime := time.Now()
		taskCtx := logger.WithKV(ctx, "bot_account", task.Username)

		logger.Debug(taskCtx, "got account for processing")

		err = nil

		userTargetsBatchSize := len(task.Targets) / postsPerBot

		for i := 0; i < postsPerBot; i++ {
			err = p.createPost(taskCtx, task.BotAccount, task.Targets[i*userTargetsBatchSize:(i+1)*userTargetsBatchSize])
			if err != nil {
				logger.Errorf(taskCtx, "failed to create post [%d]: %v", i, err)
				break
			}

			time.Sleep(3 * time.Second)
		}

		if err != nil {
			continue // один из этапов упал с ошибкой, переходим к следующему аккаунту
		}

		logger.Info(taskCtx, "all stages succeeded, saving resultsm time elapsed: %s", time.Since(startTime))

		err = q.SetAccountAsCompleted(taskCtx, task.ID)
		if err != nil {
			logger.Errorf(taskCtx, "failed to set account as completed")
		}
	}
}

func (p *worker) Login(ctx context.Context, account domain.BotAccount) error {
	if account.Headers.AuthData.SessionID != "" {
		logger.Debugf(ctx, "account %s already logged in", account.Username)
		return nil
	}

	err := p.preLoginFlow(ctx, account)
	if err != nil {
		return err
	}

	err = p.login(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (p *worker) login(ctx context.Context, account domain.BotAccount) error {
	loginReq, err := requests.PrepareLoginRequest(account)
	if err != nil {
		return fmt.Errorf("failed to prepare login request: %v", err)
	}

	resp, err := p.cli.Do(loginReq)
	if err != nil {
		return fmt.Errorf("failed to send login request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.saveFailedRequest(loginReq, resp)
		return fmt.Errorf("got %d response code for login request", resp.StatusCode)
	}

	logger.Debugf(ctx, "account '%s' successfully logged in", account.Username)

	return nil
}

func (p *worker) preLoginFlow(ctx context.Context, account domain.BotAccount) error {
	contactPrefillReq := requests.PrepareContactPointPrefillRequest(account)

	resp, err := p.cli.Do(contactPrefillReq)
	if err != nil {
		return fmt.Errorf("failed to sent ContactPointPrefillRequest: %v", err)
	}

	if resp.StatusCode != 200 {
		p.saveFailedRequest(contactPrefillReq, resp)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Errorf(ctx, "failed to close response body: %v", err)
	}

	syncLauncherReq := requests.PrepareSyncLauncherRequest(account, true)

	resp, err = p.cli.Do(syncLauncherReq)
	if err != nil {
		return fmt.Errorf("failed to sent SyncLauncher request: %v", err)
	}

	if resp.StatusCode != 200 {

	}
	return nil
}

func generateJazoest(phoneId uuid.UUID) string {
	var sum int32
	for _, s := range phoneId.String() {
		sum += s
	}

	return strconv.FormatInt(int64(sum), 10)
}

type APIResponse struct {
	Status        string `json:"status"`
	Message       string `json:"message"`
	ErrorType     string `json:"error_type"`
	ExceptionName string `json:"exception_name"`
}

func (p *worker) saveFailedRequest(r *http.Request, resp *http.Response) {

}

func (p *worker) createPost(ctx context.Context, botAccount domain.BotAccount, targetUsers []dbmodel.TargetUser) error {
	// req, _ := http.NewRequestWithContext(
	// 	transport.ContextWithProxy(context.Background(), botAccount.ProxyURL()),
	// 	"GET",
	// 	"https://2ip.ru",
	// 	nil,
	// )
	//
	// checkresp, err := p.cli.Do(req)
	// if err != nil {
	// 	return err
	// }
	// bytes, err := io.ReadAll(checkresp.Body)
	// if err != nil {
	// 	return fmt.Errorf("failed to read body: %v", err)
	// }
	//
	// fmt.Printf("got '%s' from body\n", string(bytes))

	photoUploadReq, err := requests.PrepareUploadRequest(ctx, botAccount, p.generator.Next(ctx))
	if err != nil {
		return err
	}

	pr, err := transport.FromContext()(photoUploadReq)
	if err != nil {
		return err
	}

	fmt.Println(pr)

	resp, err := p.cli.Do(photoUploadReq)
	if err != nil {
		return fmt.Errorf("failed to upload photo: %v", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}

	fmt.Printf("got '%s' from body\n", string(bytes))

	return nil
}

// '{"upload_id":"1664888837874","xsharing_nonces":{},"status":"ok"}' webp
// got '{"upload_id":"1664889100793","xsharing_nonces":{},"status":"ok"}' from body jpeg
