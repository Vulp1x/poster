package service

import (
	"context"
	"fmt"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/sessions"
	"github.com/inst-api/poster/pkg/logger"
	"goa.design/goa/v3/security"
)

// tasks_service service example implementation.
// The example methods log the requests and return zero values.
type tasksServicesrvc struct {
}

// NewTasksService returns the tasks_service service implementation.
func NewTasksService(dbmodel.DBTXFunc, sessions.Configuration) tasksservice.Service {
	return &tasksServicesrvc{}
}

// JWTAuth implements the authorization logic for service "tasks_service" for
// the "jwt" security scheme.
func (s *tasksServicesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	//
	// TBD: add authorization logic.
	//
	// In case of authorization failure this function should return
	// one of the generated error structs, e.g.:
	//
	//    return ctx, myservice.MakeUnauthorizedError("invalid token")
	//
	// Alternatively this function may return an instance of
	// goa.ServiceError with a Name field value that matches one of
	// the design error names, e.g:
	//
	//    return ctx, goa.PermanentError("unauthorized", "invalid token")
	//
	return ctx, fmt.Errorf("not implemented")
}

// создать драфт задачи
func (s *tasksServicesrvc) CreateTask(ctx context.Context, p *tasksservice.CreateTaskPayload) (res string, err error) {
	logger.Debug(ctx, "starting CreateTask with payload %#v", p)
	return
}

// начать выполнение задачи
func (s *tasksServicesrvc) StartTask(ctx context.Context, p *tasksservice.StartTaskPayload) (err error) {
	logger.Debug(ctx, "starting StartTask with payload %#v", p)
	return
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

func (s *tasksServicesrvc) UploadFile(ctx context.Context, p *tasksservice.UploadFilePayload) (err error) {
	logger.Debug(ctx, "starting UploadFile with payload %#v", p)
	return nil
}
