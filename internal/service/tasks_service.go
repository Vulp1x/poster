package service

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	authservice "github.com/inst-api/poster/gen/auth_service"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/store/tasks"
	"github.com/inst-api/poster/pkg/logger"
	"goa.design/goa/v3/security"
)

type taskStore interface {
	CreateDraftTask(ctx context.Context, userID uuid.UUID, title, textTemplate string, image []byte) (uuid.UUID, error)
	StartTask(ctx context.Context, taskID uuid.UUID) error
	StopTask(ctx context.Context, taskID uuid.UUID) error
	PrepareTask(ctx context.Context, taskID uuid.UUID, botAccounts domain.BotAccounts, proxies domain.Proxies, targets domain.TargetUsers) error
	ForceDelete(ctx context.Context, taskID uuid.UUID) error
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
		logger.Error(ctx, "routes service has nil auther")
		return ctx, tasksservice.Unauthorized("internal error")
	}

	return s.auth.JWTAuth(ctx, token, scheme)
}

// CreateTask создаёт драфт задачи
func (s *tasksServicesrvc) CreateTaskDraft(ctx context.Context, p *tasksservice.CreateTaskDraftPayload) (string, error) {
	logger.Debug(ctx, "starting CreateTask")

	userID, err := UserIDFromContext(ctx)
	if err != nil {
		logger.Errorf(ctx, "failed to get user id from context: %v", err)
		return "", tasksservice.InternalError("")
	}

	imageDecodedBytes, err := base64.StdEncoding.DecodeString(p.PostImage)
	if err != nil {
		logger.Errorf(ctx, "failed to decode base 64 string: %v", err)
		return "", tasksservice.BadRequest("invalid image")
	}

	taskID, err := s.store.CreateDraftTask(ctx, userID, p.Title, p.TextTemplate, imageDecodedBytes)
	if err != nil {
		logger.Errorf(ctx, "failed to create task: %v", err)
		return "", tasksservice.InternalError("")
	}

	return taskID.String(), nil
}

// StartTask начать выполнение задачи
func (s *tasksServicesrvc) StartTask(ctx context.Context, p *tasksservice.StartTaskPayload) error {
	logger.Debug(ctx, "starting StartTask with payload %#v", p)

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return tasksservice.BadRequest("invalid task_id")
	}

	err = s.store.StartTask(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to start task: %v", err)
		return tasksservice.InternalError("")
	}

	return nil
}

// StopTask остановить выполнение задачи
func (s *tasksServicesrvc) StopTask(ctx context.Context, p *tasksservice.StopTaskPayload) error {
	logger.Debug(ctx, "starting StopTask with payload %#v", p)

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task id from '%s': %v", p.TaskID, err)
		return tasksservice.BadRequest("invalid task_id")
	}

	err = s.store.StopTask(ctx, taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to stop task: %v", err)
		return tasksservice.InternalError("")
	}

	return nil
}

// получить задачу по id
func (s *tasksServicesrvc) GetTask(ctx context.Context, p *tasksservice.GetTaskPayload) (err error) {
	logger.Debug(ctx, "starting GetTask with payload %#v", p)
	return
}

// получить все задачи для текущего пользователя
func (s *tasksServicesrvc) ListTasks(ctx context.Context, p *tasksservice.ListTasksPayload) (err error) {
	logger.Debug(ctx, "starting ListTasks with payload %#v", p)
	return
}

func (s *tasksServicesrvc) UploadFile(ctx context.Context, p *tasksservice.UploadFilePayload) ([]*tasksservice.UploadError, error) {
	logger.Debug(ctx, "starting UploadFile with payload %#v", p)

	taskID, err := uuid.Parse(p.TaskID)
	if err != nil {
		logger.Errorf(ctx, "failed to parse task_id from '%s': %v", p.TaskID, err)
		return nil, tasksservice.BadRequest("bad task_id")
	}

	ctx = logger.WithKV(ctx, "task_id", taskID.String())

	domainAccounts, uploadErrors := domain.ParseBotAccounts(p.Bots)
	previousLen := len(uploadErrors)
	logger.Infof(ctx, "got %d bots and %d errors from %d inputs", len(domainAccounts), previousLen, len(p.Bots))

	domainProxies := domain.ParseProxies(p.Proxies, uploadErrors)
	logger.Infof(ctx, "got %d proxies and %d errors from %d inputs",
		len(domainProxies), len(uploadErrors)-previousLen, len(p.Proxies),
	)

	previousLen = len(uploadErrors)
	domainTargets := domain.ParseTargetUsers(p.Targets, uploadErrors)
	logger.Infof(ctx, "got %d targets and %d errors from %d inputs",
		len(domainTargets), len(uploadErrors)-previousLen, len(p.Targets),
	)

	err = s.store.PrepareTask(ctx, taskID, domainAccounts, domainProxies, domainTargets)
	if err != nil {
		logger.Errorf(ctx, "failed to prepare task: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			return uploadErrors, tasksservice.TaskNotFound("")
		}

		if errors.Is(err, tasks.ErrTaskInvalidStatus) {
			return uploadErrors, tasksservice.BadRequest("invalid task status")
		}

		return uploadErrors, tasksservice.InternalError("")
	}

	return uploadErrors, nil
}

func (s *tasksServicesrvc) AssignProxies(ctx context.Context, payload *tasksservice.AssignProxiesPayload) (err error) {
	// TODO implement me
	panic("implement me")
}

func (s *tasksServicesrvc) ForceDelete(ctx context.Context, p *tasksservice.ForceDeletePayload) error {
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

		return tasksservice.InternalError("")
	}

	return nil
}
