package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/uuid"
	authservice "github.com/inst-api/poster/gen/auth_service"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/mw"
	"github.com/inst-api/poster/internal/pager"
	"github.com/inst-api/poster/internal/store/tasks"
	"github.com/inst-api/poster/internal/tracer"
	"github.com/inst-api/poster/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"goa.design/goa/v3/security"
)

type taskStore interface {
	CreateDraftTask(ctx context.Context, userID uuid.UUID, title, textTemplate string, accounts []string, images [][]byte, taskType tasksservice.TaskType, opts ...tasks.DraftOption) (uuid.UUID, error)
	UpdateTask(ctx context.Context, taskID uuid.UUID, opts ...tasks.UpdateOption) (domain.Task, error)
	StartTask(ctx context.Context, taskID uuid.UUID) ([]string, error)
	StopTask(ctx context.Context, taskID uuid.UUID) error
	PrepareTask(ctx context.Context, taskID uuid.UUID, botAccounts domain.BotAccounts, proxies domain.Proxies, cheapProxies domain.Proxies, targets domain.TargetUsers, filenames *tasksservice.TaskFileNames) error
	ForceDelete(ctx context.Context, taskID uuid.UUID) error
	AssignProxies(ctx context.Context, taskID uuid.UUID) (int, error)
	GetTask(ctx context.Context, taskID uuid.UUID) (domain.TaskWithCounters, error)
	ListTasks(ctx context.Context, userID uuid.UUID) (domain.TasksWithCounters, error)
	TaskProgress(ctx context.Context, taskID uuid.UUID, pager *pager.Pager) (domain.TaskProgress, error)
	SaveVideo(ctx context.Context, taskID uuid.UUID, video []byte, filename string) (domain.Task, error)
	StartBots(ctx context.Context, taskID uuid.UUID, usernames []string) ([]string, error)
	TaskTargets(ctx context.Context, taskID uuid.UUID) (domain.Targets, error)
	TaskBots(ctx context.Context, taskID uuid.UUID) (domain.BotAccounts, error)
	StartUpdatePostContents(ctx context.Context, taskID uuid.UUID) ([]string, error)
}

// tasks_service service example implementation.
// The example methods log the requests and return zero values.
type tasksServicesrvc struct {
	auth  authservice.Auther
	store taskStore
}

// NewTasksService returns the tasks_service service implementation.
func NewTasksService(auth authservice.Auther, store taskStore) tasksservice.Service {
	return &tasksServicesrvc{auth: auth, store: store}
}

// JWTAuth implements the authorization logic for service "tasks_service" for
// the "jwt" security scheme.
func (s *tasksServicesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	if s.auth == nil {
		logger.ErrorKV(ctx, "routes service has nil auther")
		return ctx, tasksservice.Unauthorized("internal error")
	}

	return s.auth.JWTAuth(ctx, token, scheme)
}

// CreateTaskDraft создаёт драфт задачи
func (s *tasksServicesrvc) CreateTaskDraft(ctx context.Context, p *tasksservice.CreateTaskDraftPayload) (string, error) {
	logger.DebugKV(ctx, "starting CreateTask")

	userID, err := UserIDFromContext(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to get user id from context: %v", err)
		return "", tasksservice.InternalError(err.Error())
	}

	imagesDecodedBytes := make([][]byte, len(p.PostImages))
	for i, image := range p.PostImages {
		imageDecodedBytes, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			logger.Errorf(ctx, "failed to decode base 64 string: %v", err)
			return "", tasksservice.BadRequest(fmt.Sprintf("invalid %d post image", i+1))
		}
		imagesDecodedBytes[i] = imageDecodedBytes
	}

	var botProfileImages [][]byte
	for i, image := range p.BotImages {
		imageDecodedBytes, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			logger.Errorf(ctx, "failed to decode base 64 string: %v", err)
			return "", tasksservice.BadRequest(fmt.Sprintf("invalid %d bot image", i+1))
		}

		botProfileImages = append(botProfileImages, imageDecodedBytes)
	}

	taskID, err := s.store.CreateDraftTask(ctx, userID, p.Title, p.TextTemplate, p.LandingAccounts, imagesDecodedBytes, p.Type,
		tasks.CreateDraftWithAccountNames(p.BotNames),
		tasks.CreateDraftWithAccountLastNames(p.BotLastNames),
		tasks.CreateDraftWithAccountProfileImages(botProfileImages),
		tasks.CreateDraftWithAccountURLs(p.BotUrls),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to create task: %v", err)
		return "", tasksservice.InternalError(err.Error())
	}

	return taskID.String(), nil
}

// UpdateTask обновляет информацию о задаче. Не меняет статус задачи, можно вызывать сколько угодно раз.
// Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала остановить выполнение.
func (s *tasksServicesrvc) UpdateTask(ctx context.Context, p *tasksservice.UpdateTaskPayload) (*tasksservice.Task, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	var postImagesDecodedBytes [][]byte

	for _, image := range p.PostImages {
		imageDecodedBytes, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			logger.Errorf(ctx, "failed to decode base 64 string: %v", err)
			return nil, tasksservice.BadRequest("invalid image")
		}

		postImagesDecodedBytes = append(postImagesDecodedBytes, imageDecodedBytes)
	}

	var botImagesDecodedBytes [][]byte

	for _, image := range p.BotImages {
		imageDecodedBytes, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			logger.Errorf(ctx, "failed to decode base 64 string: %v", err)
			return nil, tasksservice.BadRequest("invalid image")
		}

		botImagesDecodedBytes = append(botImagesDecodedBytes, imageDecodedBytes)
	}

	task, err := s.store.UpdateTask(ctx, taskID,
		tasks.WithImagesUpdateOption(postImagesDecodedBytes),
		tasks.WithTextTemplateUpdateOption(p.TextTemplate),
		tasks.WithBotImagesUpdateOption(botImagesDecodedBytes),
		tasks.WithBotLasNamesUpdateOption(p.BotLastNames),
		tasks.WithBotNamesUpdateOption(p.BotNames),
		tasks.WithBotURLsUpdateOption(p.BotUrls),
		tasks.WithLandingAccountsUpdateOption(p.LandingAccounts),
		tasks.WithTitleUpdateOption(p.Title),
		tasks.WithFollowTargets(p.FollowTargets),
		tasks.WithNeedPhotoTags(p.NeedPhotoTags),
		tasks.WithPerPostSleepSeconds(p.PerPostSleepSeconds),
		tasks.WithPhotoTagsDelaySeconds(p.PhotoTagsDelaySeconds),
		tasks.WithPostsPerBot(p.PostsPerBot),
		tasks.WithTargetsPerPost(p.TargetsPerPost),
		tasks.WithPhotoTargetsPerPost(p.PhotoTargetsPerPost),
		tasks.WithPhotoTagsPostsPerBot(p.PhotoTagsPostsPerBot),
		tasks.WithFixedTagUpdateOption(p.TestingTagUsername),
		tasks.WithFixedPhotoTagUpdateOption(p.TestingTagUserID),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to update task: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		if errors.Is(err, tasks.ErrTaskInvalidStatus) {
			return nil, tasksservice.BadRequest("invalid status")
		}

		return nil, tasksservice.InternalError(err.Error())
	}

	return task.ToProto(), nil
}

// StartTask начать выполнение задачи
func (s *tasksServicesrvc) StartTask(ctx context.Context, p *tasksservice.StartTaskPayload) (*tasksservice.StartTaskResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.DebugKV(ctx, "starting StartTask")

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	landingAccounts, err := s.store.StartTask(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to start task: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		if errors.Is(err, tasks.ErrTaskInvalidStatus) ||
			errors.Is(err, tasks.ErrTaskWithEmptyPostImages) ||
			errors.Is(err, tasks.ErrTaskWithEmptyPostsPerBot) ||
			errors.Is(err, tasks.ErrTaskWithEmptyTargetsPerPost) ||
			errors.Is(err, tasks.ErrTaskWithEmptyLandingAccounts) ||
			errors.Is(err, tasks.ErrTaskInvalidTextTemplate) {
			return nil, tasksservice.BadRequest(err.Error())
		}

		return nil, tasksservice.InternalError(err.Error())
	}

	return &tasksservice.StartTaskResult{
		Status:          tasksservice.TaskStatus(dbmodel.StartedTaskStatus),
		TaskID:          taskID.String(),
		LandingAccounts: landingAccounts,
	}, nil
}

// StopTask остановить выполнение задачи
func (s *tasksServicesrvc) StopTask(ctx context.Context, p *tasksservice.StopTaskPayload) (*tasksservice.StopTaskResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.DebugKV(ctx, "starting StopTask", "payload", p)

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	err = s.store.StopTask(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to stop task: %v", err)
		return nil, tasksservice.InternalError(err.Error())
	}

	return &tasksservice.StopTaskResult{
		Status: tasksservice.TaskStatus(dbmodel.StoppedTaskStatus),
		TaskID: taskID.String(),
	}, nil
}

// GetTask возвращает задачу по id
func (s *tasksServicesrvc) GetTask(ctx context.Context, p *tasksservice.GetTaskPayload) (*tasksservice.Task, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)

	// trace.WithAttributes(attribute.String("task_id", p.TaskID))
	var span trace.Span
	ctx, span = tracer.Start(ctx, "tasks.Get")
	defer span.End()

	logger.DebugKV(ctx, "starting GetTask")

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task_id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("bad task_id")
	}

	task, err := s.store.GetTask(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to find task: %v", err)

		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		return nil, tasksservice.InternalError(err.Error())
	}

	return task.ToProto(), nil
}

// ListTasks получить все задачи для текущего пользователя
func (s *tasksServicesrvc) ListTasks(ctx context.Context, p *tasksservice.ListTasksPayload) ([]*tasksservice.Task, error) {
	logger.DebugKV(ctx, "starting ListTasks")

	userID, err := UserIDFromContext(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to get user id from context: %v", err)
		return nil, tasksservice.InternalError(err.Error())
	}

	domainTasks, err := s.store.ListTasks(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "failed to find task: %v", err)

		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		return nil, tasksservice.InternalError(err.Error())
	}

	return domainTasks.ToProto(), nil
}

func (s *tasksServicesrvc) UploadFiles(ctx context.Context, p *tasksservice.UploadFilesPayload) (*tasksservice.UploadFilesResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	ctx, span := tracer.Start(ctx, "UploadFiles", trace.WithAttributes(
		attribute.String("bots", p.Filenames.BotsFilename),
		attribute.String("targets", p.Filenames.TargetsFilename),
		attribute.String("cheap_proxies", p.Filenames.CheapProxiesFilename),
		attribute.String("residential_proxies", p.Filenames.ResidentialProxiesFilename),
	))
	defer span.End()

	logger.Infof(ctx, "starting UploadFile with filenames %+v", p.Filenames)

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task_id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("bad task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	domainAccounts, uploadErrors := domain.ParseBotAccounts(ctx, p.Bots)
	previousLen := len(uploadErrors)
	logger.Infof(ctx, "got %d bots and %d errors from %d inputs", len(domainAccounts), previousLen, len(p.Bots))

	domainProxies := domain.ParseProxies(p.ResidentialProxies, uploadErrors)
	logger.Infof(ctx, "got %d residential proxies and %d errors from %d inputs",
		len(domainProxies), len(uploadErrors)-previousLen, len(p.ResidentialProxies),
	)

	cheapProxies := domain.ParseProxies(p.CheapProxies, uploadErrors)
	logger.Infof(ctx, "got %d cheap proxies and %d errors from %d inputs",
		len(cheapProxies), len(uploadErrors)-previousLen, len(p.CheapProxies),
	)

	previousLen = len(uploadErrors)
	domainTargets := domain.ParseTargetUsers(p.Targets, uploadErrors)
	logger.Infof(ctx, "got %d targets and %d errors from %d inputs",
		len(domainTargets), len(uploadErrors)-previousLen, len(p.Targets),
	)

	result := &tasksservice.UploadFilesResult{
		UploadErrors: uploadErrors,
		Status:       -1,
	}

	err = s.store.PrepareTask(ctx, taskID, domainAccounts, domainProxies, cheapProxies, domainTargets, p.Filenames)
	if err != nil {
		logger.Errorf(ctx, "failed to prepare task: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return result, tasksservice.TaskNotFound("")
		}

		if errors.Is(err, tasks.ErrTaskInvalidStatus) {
			return result, tasksservice.BadRequest("invalid task status")
		}

		return result, internalErr(err)
	}

	result.Status = tasksservice.TaskStatus(dbmodel.DataUploadedTaskStatus)

	return result, nil
}

func (s *tasksServicesrvc) UploadVideo(ctx context.Context, p *tasksservice.UploadVideoPayload) (*tasksservice.UploadVideoResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task_id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("bad task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	if p.Filename == nil {
		shortID := mw.ShortID()
		logger.Warnf(ctx, "got empty file name, using %s instead", shortID)
		p.Filename = &shortID
	}

	task, err := s.store.SaveVideo(ctx, taskID, p.Video, *p.Filename)
	if err != nil {
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")

		}

		if errors.Is(err, tasks.ErrUnexpectedTaskType) || errors.Is(err, tasks.ErrTaskInvalidStatus) {
			return nil, tasksservice.BadRequest(err.Error())
		}

		return nil, internalErr(err)
	}

	return &tasksservice.UploadVideoResult{Status: tasksservice.TaskStatus(task.Status)}, nil
}

func (s *tasksServicesrvc) AssignProxies(ctx context.Context, p *tasksservice.AssignProxiesPayload) (*tasksservice.AssignProxiesResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task_id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("bad task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	botAccounts, err := s.store.AssignProxies(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to assign proxies: %v", err)

		if errors.Is(err, tasks.ErrTaskInvalidStatus) {
			return nil, tasksservice.BadRequest("invalid task status")
		}

		return nil, internalErr(err)
	}

	return &tasksservice.AssignProxiesResult{
		BotsNumber: botAccounts,
		Status:     tasksservice.TaskStatus(dbmodel.ReadyTaskStatus),
		TaskID:     taskID.String(),
	}, nil

}

func (s *tasksServicesrvc) ForceDelete(ctx context.Context, p *tasksservice.ForceDeletePayload) error {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task_id from '%s': %v", p.TaskID, err)
		return tasksservice.BadRequest("bad task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	err = s.store.ForceDelete(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to force delete task: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return tasksservice.TaskNotFound("")
		}

		return internalErr(err)
	}

	return nil
}

func (s *tasksServicesrvc) GetProgress(ctx context.Context, p *tasksservice.GetProgressPayload) (*tasksservice.TaskProgress, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.Debugf(ctx, "get progress of task %s", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	per := pager.NewPagePer(p.Page, p.PageSize)
	per.SetSortColumns(p.Sort)
	if p.SortDescending {
		per.SetReverseOrderSort(p.Sort)
	}

	domainProgress, err := s.store.TaskProgress(ctx, taskID, per)
	if err != nil {
		logger.Errorf(ctx, "failed to get task progress: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		if errors.Is(err, tasks.ErrTaskInvalidStatus) {
			return nil, tasksservice.BadRequest("invalid status")
		}

		return nil, internalErr(err)
	}

	return domainProgress.ToProto(), nil
}

func (s *tasksServicesrvc) PartialStartTask(ctx context.Context, p *tasksservice.PartialStartTaskPayload) (*tasksservice.PartialStartTaskResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.Infof(ctx, "partial start task %s", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	usernames, err := s.store.StartBots(ctx, taskID, p.Usernames)
	if err != nil {
		logger.Errorf(ctx, "failed to start bots: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		return nil, internalErr(err)
	}

	return &tasksservice.PartialStartTaskResult{
		TaskID:    taskID.String(),
		Succeeded: usernames,
	}, nil
}

func (s *tasksServicesrvc) UpdatePostContents(ctx context.Context, p *tasksservice.UpdatePostContentsPayload) (*tasksservice.UpdatePostContentsResult, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.Infof(ctx, "starting to update post contents %s", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	landingAccounts, err := s.store.StartUpdatePostContents(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to start updating post contents: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		return nil, internalErr(err)
	}

	return &tasksservice.UpdatePostContentsResult{
		Status:          tasksservice.TaskStatus(dbmodel.UpdatingPostContentsTaskStatus),
		TaskID:          taskID.String(),
		LandingAccounts: landingAccounts,
	}, nil
}

func (s *tasksServicesrvc) DownloadTargets(ctx context.Context, p *tasksservice.DownloadTargetsPayload) ([]string, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.Infof(ctx, "download task targets %s", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	targets, err := s.store.TaskTargets(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to download targets: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		return nil, internalErr(err)
	}

	return targets.ToProto(p.Format), nil
}

func (s *tasksServicesrvc) DownloadBots(ctx context.Context, p *tasksservice.DownloadBotsPayload) ([]string, error) {
	ctx = logger.WithKV(ctx, "task_id", p.TaskID)
	logger.Infof(ctx, "download task bots %s", p.TaskID)
	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("invalid task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	bots, err := s.store.TaskBots(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to download bots: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return nil, tasksservice.TaskNotFound("")
		}

		return nil, internalErr(err)
	}

	return bots.ToProto(p.Proxies), nil
}

func internalErr(err error) tasksservice.InternalError {
	return tasksservice.InternalError(err.Error())
}
