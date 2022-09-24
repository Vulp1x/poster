// Code generated by goa v3.8.5, DO NOT EDIT.
//
// tasks_service HTTP server
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package server

import (
	"context"
	"mime/multipart"
	"net/http"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Server lists the tasks_service service endpoint HTTP handlers.
type Server struct {
	Mounts          []*MountPoint
	CreateTaskDraft http.Handler
	UploadFile      http.Handler
	StartTask       http.Handler
	StopTask        http.Handler
	GetTask         http.Handler
	ListTasks       http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// TasksServiceUploadFileDecoderFunc is the type to decode multipart request
// for the "tasks_service" service "upload file" endpoint.
type TasksServiceUploadFileDecoderFunc func(*multipart.Reader, **tasksservice.UploadFilePayload) error

// New instantiates HTTP handlers for all the tasks_service service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *tasksservice.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
	tasksServiceUploadFileDecoderFn TasksServiceUploadFileDecoderFunc,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"CreateTaskDraft", "POST", "/api/tasks/draft"},
			{"UploadFile", "POST", "/api/tasks/{task_id}/upload"},
			{"StartTask", "POST", "/api/tasks/{task_id}/start"},
			{"StopTask", "POST", "/api/tasks/{task_id}/stop"},
			{"GetTask", "GET", "/api/tasks/{task_id}/"},
			{"ListTasks", "GET", "/api/tasks/"},
		},
		CreateTaskDraft: NewCreateTaskDraftHandler(e.CreateTaskDraft, mux, decoder, encoder, errhandler, formatter),
		UploadFile:      NewUploadFileHandler(e.UploadFile, mux, NewTasksServiceUploadFileDecoder(mux, tasksServiceUploadFileDecoderFn), encoder, errhandler, formatter),
		StartTask:       NewStartTaskHandler(e.StartTask, mux, decoder, encoder, errhandler, formatter),
		StopTask:        NewStopTaskHandler(e.StopTask, mux, decoder, encoder, errhandler, formatter),
		GetTask:         NewGetTaskHandler(e.GetTask, mux, decoder, encoder, errhandler, formatter),
		ListTasks:       NewListTasksHandler(e.ListTasks, mux, decoder, encoder, errhandler, formatter),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "tasks_service" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.CreateTaskDraft = m(s.CreateTaskDraft)
	s.UploadFile = m(s.UploadFile)
	s.StartTask = m(s.StartTask)
	s.StopTask = m(s.StopTask)
	s.GetTask = m(s.GetTask)
	s.ListTasks = m(s.ListTasks)
}

// Mount configures the mux to serve the tasks_service endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountCreateTaskDraftHandler(mux, h.CreateTaskDraft)
	MountUploadFileHandler(mux, h.UploadFile)
	MountStartTaskHandler(mux, h.StartTask)
	MountStopTaskHandler(mux, h.StopTask)
	MountGetTaskHandler(mux, h.GetTask)
	MountListTasksHandler(mux, h.ListTasks)
}

// Mount configures the mux to serve the tasks_service endpoints.
func (s *Server) Mount(mux goahttp.Muxer) {
	Mount(mux, s)
}

// MountCreateTaskDraftHandler configures the mux to serve the "tasks_service"
// service "create task draft" endpoint.
func MountCreateTaskDraftHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/tasks/draft", f)
}

// NewCreateTaskDraftHandler creates a HTTP handler which loads the HTTP
// request and calls the "tasks_service" service "create task draft" endpoint.
func NewCreateTaskDraftHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeCreateTaskDraftRequest(mux, decoder)
		encodeResponse = EncodeCreateTaskDraftResponse(encoder)
		encodeError    = EncodeCreateTaskDraftError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "create task draft")
		ctx = context.WithValue(ctx, goa.ServiceKey, "tasks_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountUploadFileHandler configures the mux to serve the "tasks_service"
// service "upload file" endpoint.
func MountUploadFileHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/tasks/{task_id}/upload", f)
}

// NewUploadFileHandler creates a HTTP handler which loads the HTTP request and
// calls the "tasks_service" service "upload file" endpoint.
func NewUploadFileHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUploadFileRequest(mux, decoder)
		encodeResponse = EncodeUploadFileResponse(encoder)
		encodeError    = EncodeUploadFileError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "upload file")
		ctx = context.WithValue(ctx, goa.ServiceKey, "tasks_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountStartTaskHandler configures the mux to serve the "tasks_service"
// service "start task" endpoint.
func MountStartTaskHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/tasks/{task_id}/start", f)
}

// NewStartTaskHandler creates a HTTP handler which loads the HTTP request and
// calls the "tasks_service" service "start task" endpoint.
func NewStartTaskHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeStartTaskRequest(mux, decoder)
		encodeResponse = EncodeStartTaskResponse(encoder)
		encodeError    = EncodeStartTaskError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "start task")
		ctx = context.WithValue(ctx, goa.ServiceKey, "tasks_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountStopTaskHandler configures the mux to serve the "tasks_service" service
// "stop task" endpoint.
func MountStopTaskHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/tasks/{task_id}/stop", f)
}

// NewStopTaskHandler creates a HTTP handler which loads the HTTP request and
// calls the "tasks_service" service "stop task" endpoint.
func NewStopTaskHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeStopTaskRequest(mux, decoder)
		encodeResponse = EncodeStopTaskResponse(encoder)
		encodeError    = EncodeStopTaskError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "stop task")
		ctx = context.WithValue(ctx, goa.ServiceKey, "tasks_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountGetTaskHandler configures the mux to serve the "tasks_service" service
// "get task" endpoint.
func MountGetTaskHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/tasks/{task_id}/", f)
}

// NewGetTaskHandler creates a HTTP handler which loads the HTTP request and
// calls the "tasks_service" service "get task" endpoint.
func NewGetTaskHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeGetTaskRequest(mux, decoder)
		encodeResponse = EncodeGetTaskResponse(encoder)
		encodeError    = EncodeGetTaskError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "get task")
		ctx = context.WithValue(ctx, goa.ServiceKey, "tasks_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountListTasksHandler configures the mux to serve the "tasks_service"
// service "list tasks" endpoint.
func MountListTasksHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/tasks/", f)
}

// NewListTasksHandler creates a HTTP handler which loads the HTTP request and
// calls the "tasks_service" service "list tasks" endpoint.
func NewListTasksHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeListTasksRequest(mux, decoder)
		encodeResponse = EncodeListTasksResponse(encoder)
		encodeError    = EncodeListTasksError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "list tasks")
		ctx = context.WithValue(ctx, goa.ServiceKey, "tasks_service")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
