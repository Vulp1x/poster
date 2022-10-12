// Code generated by goa v3.8.5, DO NOT EDIT.
//
// tasks_service HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package server

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeCreateTaskDraftResponse returns an encoder for responses returned by
// the tasks_service create task draft endpoint.
func EncodeCreateTaskDraftResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(string)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeCreateTaskDraftRequest returns a decoder for requests sent to the
// tasks_service create task draft endpoint.
func DecodeCreateTaskDraftRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body CreateTaskDraftRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}
		err = ValidateCreateTaskDraftRequestBody(&body)
		if err != nil {
			return nil, err
		}

		var (
			token string
		)
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewCreateTaskDraftPayload(&body, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeCreateTaskDraftError returns an encoder for errors returned by the
// create task draft tasks_service endpoint.
func EncodeCreateTaskDraftError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeUpdateTaskResponse returns an encoder for responses returned by the
// tasks_service update task endpoint.
func EncodeUpdateTaskResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*tasksservice.Task)
		enc := encoder(ctx, w)
		body := NewUpdateTaskOKResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeUpdateTaskRequest returns a decoder for requests sent to the
// tasks_service update task endpoint.
func DecodeUpdateTaskRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			body UpdateTaskRequestBody
			err  error
		)
		err = decoder(r).Decode(&body)
		if err != nil {
			if err == io.EOF {
				return nil, goa.MissingPayloadError()
			}
			return nil, goa.DecodePayloadError(err.Error())
		}

		var (
			taskID string
			token  string

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewUpdateTaskPayload(&body, taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeUpdateTaskError returns an encoder for errors returned by the update
// task tasks_service endpoint.
func EncodeUpdateTaskError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeUploadFilesResponse returns an encoder for responses returned by the
// tasks_service upload files endpoint.
func EncodeUploadFilesResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.([]*tasksservice.UploadError)
		enc := encoder(ctx, w)
		body := NewUploadFilesResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeUploadFilesRequest returns a decoder for requests sent to the
// tasks_service upload files endpoint.
func DecodeUploadFilesRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var payload *tasksservice.UploadFilesPayload
		if err := decoder(r).Decode(&payload); err != nil {
			return nil, goa.DecodePayloadError(err.Error())
		}
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// NewTasksServiceUploadFilesDecoder returns a decoder to decode the multipart
// request for the "tasks_service" service "upload files" endpoint.
func NewTasksServiceUploadFilesDecoder(mux goahttp.Muxer, tasksServiceUploadFilesDecoderFn TasksServiceUploadFilesDecoderFunc) func(r *http.Request) goahttp.Decoder {
	return func(r *http.Request) goahttp.Decoder {
		return goahttp.EncodingFunc(func(v interface{}) error {
			mr, merr := r.MultipartReader()
			if merr != nil {
				return merr
			}
			p := v.(**tasksservice.UploadFilesPayload)
			if err := tasksServiceUploadFilesDecoderFn(mr, p); err != nil {
				return err
			}

			var (
				taskID string
				token  string
				err    error

				params = mux.Vars(r)
			)
			taskID = params["task_id"]
			token = r.Header.Get("Authorization")
			if token == "" {
				err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
			}
			if err != nil {
				return err
			}
			(*p).TaskID = taskID
			(*p).Token = token
			return nil
		})
	}
}

// EncodeUploadFilesError returns an encoder for errors returned by the upload
// files tasks_service endpoint.
func EncodeUploadFilesError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeAssignProxiesResponse returns an encoder for responses returned by the
// tasks_service assign proxies endpoint.
func EncodeAssignProxiesResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(int)
		enc := encoder(ctx, w)
		body := res
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeAssignProxiesRequest returns a decoder for requests sent to the
// tasks_service assign proxies endpoint.
func DecodeAssignProxiesRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			token  string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewAssignProxiesPayload(taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeAssignProxiesError returns an encoder for errors returned by the
// assign proxies tasks_service endpoint.
func EncodeAssignProxiesError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeForceDeleteResponse returns an encoder for responses returned by the
// tasks_service force delete endpoint.
func EncodeForceDeleteResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// DecodeForceDeleteRequest returns a decoder for requests sent to the
// tasks_service force delete endpoint.
func DecodeForceDeleteRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			token  string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewForceDeletePayload(taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeForceDeleteError returns an encoder for errors returned by the force
// delete tasks_service endpoint.
func EncodeForceDeleteError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeStartTaskResponse returns an encoder for responses returned by the
// tasks_service start task endpoint.
func EncodeStartTaskResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// DecodeStartTaskRequest returns a decoder for requests sent to the
// tasks_service start task endpoint.
func DecodeStartTaskRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			token  string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewStartTaskPayload(taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeStartTaskError returns an encoder for errors returned by the start
// task tasks_service endpoint.
func EncodeStartTaskError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeStopTaskResponse returns an encoder for responses returned by the
// tasks_service stop task endpoint.
func EncodeStopTaskResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

// DecodeStopTaskRequest returns a decoder for requests sent to the
// tasks_service stop task endpoint.
func DecodeStopTaskRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			token  string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewStopTaskPayload(taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeStopTaskError returns an encoder for errors returned by the stop task
// tasks_service endpoint.
func EncodeStopTaskError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetTaskResponse returns an encoder for responses returned by the
// tasks_service get task endpoint.
func EncodeGetTaskResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*tasksservice.Task)
		enc := encoder(ctx, w)
		body := NewGetTaskOKResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetTaskRequest returns a decoder for requests sent to the
// tasks_service get task endpoint.
func DecodeGetTaskRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			token  string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetTaskPayload(taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeGetTaskError returns an encoder for errors returned by the get task
// tasks_service endpoint.
func EncodeGetTaskError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeGetProgressResponse returns an encoder for responses returned by the
// tasks_service get progress endpoint.
func EncodeGetProgressResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.([]*tasksservice.BotsProgress)
		enc := encoder(ctx, w)
		body := NewGetProgressResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeGetProgressRequest returns a decoder for requests sent to the
// tasks_service get progress endpoint.
func DecodeGetProgressRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			taskID string
			token  string
			err    error

			params = mux.Vars(r)
		)
		taskID = params["task_id"]
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewGetProgressPayload(taskID, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeGetProgressError returns an encoder for errors returned by the get
// progress tasks_service endpoint.
func EncodeGetProgressError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// EncodeListTasksResponse returns an encoder for responses returned by the
// tasks_service list tasks endpoint.
func EncodeListTasksResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.([]*tasksservice.Task)
		enc := encoder(ctx, w)
		body := NewListTasksResponseBody(res)
		w.WriteHeader(http.StatusOK)
		return enc.Encode(body)
	}
}

// DecodeListTasksRequest returns a decoder for requests sent to the
// tasks_service list tasks endpoint.
func DecodeListTasksRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			token string
			err   error
		)
		token = r.Header.Get("Authorization")
		if token == "" {
			err = goa.MergeErrors(err, goa.MissingFieldError("Authorization", "header"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewListTasksPayload(token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}

		return payload, nil
	}
}

// EncodeListTasksError returns an encoder for errors returned by the list
// tasks tasks_service endpoint.
func EncodeListTasksError(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder, formatter func(err error) goahttp.Statuser) func(context.Context, http.ResponseWriter, error) error {
	encodeError := goahttp.ErrorEncoder(encoder, formatter)
	return func(ctx context.Context, w http.ResponseWriter, v error) error {
		var en ErrorNamer
		if !errors.As(v, &en) {
			return encodeError(ctx, w, v)
		}
		switch en.ErrorName() {
		case "bad request":
			var res tasksservice.BadRequest
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusBadRequest)
			return enc.Encode(body)
		case "internal error":
			var res tasksservice.InternalError
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusInternalServerError)
			return enc.Encode(body)
		case "task not found":
			var res tasksservice.TaskNotFound
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusNotFound)
			return enc.Encode(body)
		case "unauthorized":
			var res tasksservice.Unauthorized
			errors.As(v, &res)
			enc := encoder(ctx, w)
			body := res
			w.Header().Set("goa-error", res.ErrorName())
			w.WriteHeader(http.StatusUnauthorized)
			return enc.Encode(body)
		default:
			return encodeError(ctx, w, v)
		}
	}
}

// unmarshalTaskFileNamesRequestBodyToTasksserviceTaskFileNames builds a value
// of type *tasksservice.TaskFileNames from a value of type
// *TaskFileNamesRequestBody.
func unmarshalTaskFileNamesRequestBodyToTasksserviceTaskFileNames(v *TaskFileNamesRequestBody) *tasksservice.TaskFileNames {
	res := &tasksservice.TaskFileNames{
		BotsFilename:    *v.BotsFilename,
		ProxiesFilename: *v.ProxiesFilename,
		TargetsFilename: *v.TargetsFilename,
	}

	return res
}

// unmarshalBotAccountRecordRequestBodyToTasksserviceBotAccountRecord builds a
// value of type *tasksservice.BotAccountRecord from a value of type
// *BotAccountRecordRequestBody.
func unmarshalBotAccountRecordRequestBodyToTasksserviceBotAccountRecord(v *BotAccountRecordRequestBody) *tasksservice.BotAccountRecord {
	res := &tasksservice.BotAccountRecord{
		LineNumber: *v.LineNumber,
	}
	res.Record = make([]string, len(v.Record))
	for i, val := range v.Record {
		res.Record[i] = val
	}

	return res
}

// unmarshalProxyRecordRequestBodyToTasksserviceProxyRecord builds a value of
// type *tasksservice.ProxyRecord from a value of type *ProxyRecordRequestBody.
func unmarshalProxyRecordRequestBodyToTasksserviceProxyRecord(v *ProxyRecordRequestBody) *tasksservice.ProxyRecord {
	res := &tasksservice.ProxyRecord{
		LineNumber: *v.LineNumber,
	}
	res.Record = make([]string, len(v.Record))
	for i, val := range v.Record {
		res.Record[i] = val
	}

	return res
}

// unmarshalTargetUserRecordRequestBodyToTasksserviceTargetUserRecord builds a
// value of type *tasksservice.TargetUserRecord from a value of type
// *TargetUserRecordRequestBody.
func unmarshalTargetUserRecordRequestBodyToTasksserviceTargetUserRecord(v *TargetUserRecordRequestBody) *tasksservice.TargetUserRecord {
	res := &tasksservice.TargetUserRecord{
		LineNumber: *v.LineNumber,
	}
	res.Record = make([]string, len(v.Record))
	for i, val := range v.Record {
		res.Record[i] = val
	}

	return res
}

// marshalTasksserviceUploadErrorToUploadErrorResponse builds a value of type
// *UploadErrorResponse from a value of type *tasksservice.UploadError.
func marshalTasksserviceUploadErrorToUploadErrorResponse(v *tasksservice.UploadError) *UploadErrorResponse {
	res := &UploadErrorResponse{
		Type:   v.Type,
		Line:   v.Line,
		Input:  v.Input,
		Reason: v.Reason,
	}

	return res
}

// marshalTasksserviceBotsProgressToBotsProgressResponse builds a value of type
// *BotsProgressResponse from a value of type *tasksservice.BotsProgress.
func marshalTasksserviceBotsProgressToBotsProgressResponse(v *tasksservice.BotsProgress) *BotsProgressResponse {
	res := &BotsProgressResponse{
		UserName:   v.UserName,
		PostsCount: v.PostsCount,
		Status:     v.Status,
	}

	return res
}

// marshalTasksserviceTaskToTaskResponse builds a value of type *TaskResponse
// from a value of type *tasksservice.Task.
func marshalTasksserviceTaskToTaskResponse(v *tasksservice.Task) *TaskResponse {
	res := &TaskResponse{
		ID:              v.ID,
		TextTemplate:    v.TextTemplate,
		Image:           v.Image,
		Status:          v.Status,
		Title:           v.Title,
		BotsNum:         v.BotsNum,
		ProxiesNum:      v.ProxiesNum,
		TargetsNum:      v.TargetsNum,
		BotsFilename:    v.BotsFilename,
		ProxiesFilename: v.ProxiesFilename,
		TargetsFilename: v.TargetsFilename,
	}

	return res
}
