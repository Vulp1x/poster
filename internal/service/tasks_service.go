package service

import (
	"context"
	"encoding/base64"

	routesservice "github.com/SimpleRouting/RoutingAppService/gen/routes_service"
	"github.com/google/uuid"
	authservice "github.com/inst-api/poster/gen/auth_service"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/pkg/logger"
	"goa.design/goa/v3/security"
)

type taskStore interface {
	CreateDraftTask(ctx context.Context, userID uuid.UUID, title, textTemplate string, image []byte) (uuid.UUID, error)
	StartTask(ctx context.Context, taskID uuid.UUID) error
	StopTask(ctx context.Context, taskID uuid.UUID) error
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
		return ctx, routesservice.Unauthorized("internal error")
	}

	return s.auth.JWTAuth(ctx, token, scheme)
}

// CreateTask создаёт драфт задачи
func (s *tasksServicesrvc) CreateTaskDraft(ctx context.Context, p *tasksservice.CreateTaskDraftPayload) (string, error) {
	logger.Debug(ctx, "starting CreateTask with payload %#v", p)

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
	return nil, nil
}
