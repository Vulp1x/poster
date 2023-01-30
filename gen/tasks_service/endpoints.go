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
	CreateTaskDraft    goa.Endpoint
	UpdateTask         goa.Endpoint
	UploadVideo        goa.Endpoint
	UploadFiles        goa.Endpoint
	AssignProxies      goa.Endpoint
	ForceDelete        goa.Endpoint
	StartTask          goa.Endpoint
	PartialStartTask   goa.Endpoint
	UpdatePostContents goa.Endpoint
	StopTask           goa.Endpoint
	GetTask            goa.Endpoint
	GetProgress        goa.Endpoint
	GetEditingProgress goa.Endpoint
	ListTasks          goa.Endpoint
	DownloadTargets    goa.Endpoint
	DownloadBots       goa.Endpoint
}

// NewEndpoints wraps the methods of the "tasks_service" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		CreateTaskDraft:    NewCreateTaskDraftEndpoint(s, a.JWTAuth),
		UpdateTask:         NewUpdateTaskEndpoint(s, a.JWTAuth),
		UploadVideo:        NewUploadVideoEndpoint(s, a.JWTAuth),
		UploadFiles:        NewUploadFilesEndpoint(s, a.JWTAuth),
		AssignProxies:      NewAssignProxiesEndpoint(s, a.JWTAuth),
		ForceDelete:        NewForceDeleteEndpoint(s, a.JWTAuth),
		StartTask:          NewStartTaskEndpoint(s, a.JWTAuth),
		PartialStartTask:   NewPartialStartTaskEndpoint(s, a.JWTAuth),
		UpdatePostContents: NewUpdatePostContentsEndpoint(s, a.JWTAuth),
		StopTask:           NewStopTaskEndpoint(s, a.JWTAuth),
		GetTask:            NewGetTaskEndpoint(s, a.JWTAuth),
		GetProgress:        NewGetProgressEndpoint(s, a.JWTAuth),
		GetEditingProgress: NewGetEditingProgressEndpoint(s, a.JWTAuth),
		ListTasks:          NewListTasksEndpoint(s, a.JWTAuth),
		DownloadTargets:    NewDownloadTargetsEndpoint(s, a.JWTAuth),
		DownloadBots:       NewDownloadBotsEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "tasks_service" service
// endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.CreateTaskDraft = m(e.CreateTaskDraft)
	e.UpdateTask = m(e.UpdateTask)
	e.UploadVideo = m(e.UploadVideo)
	e.UploadFiles = m(e.UploadFiles)
	e.AssignProxies = m(e.AssignProxies)
	e.ForceDelete = m(e.ForceDelete)
	e.StartTask = m(e.StartTask)
	e.PartialStartTask = m(e.PartialStartTask)
	e.UpdatePostContents = m(e.UpdatePostContents)
	e.StopTask = m(e.StopTask)
	e.GetTask = m(e.GetTask)
	e.GetProgress = m(e.GetProgress)
	e.GetEditingProgress = m(e.GetEditingProgress)
	e.ListTasks = m(e.ListTasks)
	e.DownloadTargets = m(e.DownloadTargets)
	e.DownloadBots = m(e.DownloadBots)
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

// NewUpdateTaskEndpoint returns an endpoint function that calls the method
// "update task" of service "tasks_service".
func NewUpdateTaskEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UpdateTaskPayload)
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
		return s.UpdateTask(ctx, p)
	}
}

// NewUploadVideoEndpoint returns an endpoint function that calls the method
// "upload video" of service "tasks_service".
func NewUploadVideoEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UploadVideoPayload)
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
		return s.UploadVideo(ctx, p)
	}
}

// NewUploadFilesEndpoint returns an endpoint function that calls the method
// "upload files" of service "tasks_service".
func NewUploadFilesEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UploadFilesPayload)
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
		return s.UploadFiles(ctx, p)
	}
}

// NewAssignProxiesEndpoint returns an endpoint function that calls the method
// "assign proxies" of service "tasks_service".
func NewAssignProxiesEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*AssignProxiesPayload)
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
		return s.AssignProxies(ctx, p)
	}
}

// NewForceDeleteEndpoint returns an endpoint function that calls the method
// "force delete" of service "tasks_service".
func NewForceDeleteEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*ForceDeletePayload)
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
		return nil, s.ForceDelete(ctx, p)
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
		return s.StartTask(ctx, p)
	}
}

// NewPartialStartTaskEndpoint returns an endpoint function that calls the
// method "partial start task" of service "tasks_service".
func NewPartialStartTaskEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*PartialStartTaskPayload)
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
		return s.PartialStartTask(ctx, p)
	}
}

// NewUpdatePostContentsEndpoint returns an endpoint function that calls the
// method "update post contents" of service "tasks_service".
func NewUpdatePostContentsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*UpdatePostContentsPayload)
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
		return s.UpdatePostContents(ctx, p)
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
		return s.StopTask(ctx, p)
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
		return s.GetTask(ctx, p)
	}
}

// NewGetProgressEndpoint returns an endpoint function that calls the method
// "get progress" of service "tasks_service".
func NewGetProgressEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetProgressPayload)
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
		return s.GetProgress(ctx, p)
	}
}

// NewGetEditingProgressEndpoint returns an endpoint function that calls the
// method "get editing progress" of service "tasks_service".
func NewGetEditingProgressEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*GetEditingProgressPayload)
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
		return s.GetEditingProgress(ctx, p)
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
		return s.ListTasks(ctx, p)
	}
}

// NewDownloadTargetsEndpoint returns an endpoint function that calls the
// method "download targets" of service "tasks_service".
func NewDownloadTargetsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DownloadTargetsPayload)
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
		return s.DownloadTargets(ctx, p)
	}
}

// NewDownloadBotsEndpoint returns an endpoint function that calls the method
// "download bots" of service "tasks_service".
func NewDownloadBotsEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*DownloadBotsPayload)
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
		return s.DownloadBots(ctx, p)
	}
}
