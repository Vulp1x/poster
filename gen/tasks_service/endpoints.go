// Code generated by goa v3.8.5, DO NOT EDIT.
//
// tasks_service endpoints
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package tasksservice

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "tasks_service" service endpoints.
type Endpoints struct {
	CreateTaskDraft goa.Endpoint
	UploadFile      goa.Endpoint
	StartTask       goa.Endpoint
	StopTask        goa.Endpoint
	GetTask         goa.Endpoint
	ListTasks       goa.Endpoint
}

// NewEndpoints wraps the methods of the "tasks_service" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		CreateTaskDraft: NewCreateTaskDraftEndpoint(s, a.JWTAuth),
		UploadFile:      NewUploadFileEndpoint(s, a.JWTAuth),
		StartTask:       NewStartTaskEndpoint(s, a.JWTAuth),
		StopTask:        NewStopTaskEndpoint(s, a.JWTAuth),
		GetTask:         NewGetTaskEndpoint(s, a.JWTAuth),
		ListTasks:       NewListTasksEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "tasks_service" service
// endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.CreateTaskDraft = m(e.CreateTaskDraft)
	e.UploadFile = m(e.UploadFile)
	e.StartTask = m(e.StartTask)
	e.StopTask = m(e.StopTask)
	e.GetTask = m(e.GetTask)
	e.ListTasks = m(e.ListTasks)
}

// NewCreateTaskDraftEndpoint returns an endpoint function that calls the
// method "create task draft" of service "tasks_service".
func NewCreateTaskDraftEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*CreateTaskDraftPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.CreateTaskDraft(ctx, p)
	}
}

// NewUploadFileEndpoint returns an endpoint function that calls the method
// "upload file" of service "tasks_service".
func NewUploadFileEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UploadFilePayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.UploadFile(ctx, p)
	}
}

// NewStartTaskEndpoint returns an endpoint function that calls the method
// "start task" of service "tasks_service".
func NewStartTaskEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*StartTaskPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.StartTask(ctx, p)
	}
}

// NewStopTaskEndpoint returns an endpoint function that calls the method "stop
// task" of service "tasks_service".
func NewStopTaskEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*StopTaskPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.StopTask(ctx, p)
	}
}

// NewGetTaskEndpoint returns an endpoint function that calls the method "get
// task" of service "tasks_service".
func NewGetTaskEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetTaskPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.GetTask(ctx, p)
	}
}

// NewListTasksEndpoint returns an endpoint function that calls the method
// "list tasks" of service "tasks_service".
func NewListTasksEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ListTasksPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"driver", "admin"},
			RequiredScopes: []string{},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return nil, s.ListTasks(ctx, p)
	}
}
