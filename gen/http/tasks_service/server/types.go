// Code generated by goa v3.8.5, DO NOT EDIT.
//
// tasks_service HTTP server types
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package server

import (
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	goa "goa.design/goa/v3/pkg"
)

// CreateTaskDraftRequestBody is the type of the "tasks_service" service
// "create task draft" endpoint HTTP request body.
type CreateTaskDraftRequestBody struct {
	// название задачи
	Title *string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
	// шаблон для подписи под постом
	TextTemplate *string `json:"text_template"`
	// фотография для постов
	PostImage *string `json:"post_image"`
}

// UploadFileRequestBody is the type of the "tasks_service" service "upload
// file" endpoint HTTP request body.
type UploadFileRequestBody struct {
	// список ботов
	Bots []*BotAccountRecordRequestBody `form:"bots,omitempty" json:"bots,omitempty" xml:"bots,omitempty"`
	// список проксей для использования
	Proxies []*ProxyRecordRequestBody `form:"proxies,omitempty" json:"proxies,omitempty" xml:"proxies,omitempty"`
	// список аккаунтов, которым показать надо рекламу
	Targets []*TargetUserRecordRequestBody `form:"targets,omitempty" json:"targets,omitempty" xml:"targets,omitempty"`
}

// UploadFileResponseBody is the type of the "tasks_service" service "upload
// file" endpoint HTTP response body.
type UploadFileResponseBody []*UploadErrorResponse

// UploadErrorResponse is used to define fields on response body types.
type UploadErrorResponse struct {
	// 1 - список ботов
	// 2 - список прокси
	// 3 - список получателей рекламы
	Type int `form:"type" json:"type" xml:"type"`
	Line int `form:"line" json:"line" xml:"line"`
	// номер порта
	Input  string `form:"input" json:"input" xml:"input"`
	Reason string `form:"reason" json:"reason" xml:"reason"`
}

// BotAccountRecordRequestBody is used to define fields on request body types.
type BotAccountRecordRequestBody struct {
	Record []string `form:"record,omitempty" json:"record,omitempty" xml:"record,omitempty"`
	// номер строки в исходном файле
	LineNumber *int `json:"line_number"`
}

// ProxyRecordRequestBody is used to define fields on request body types.
type ProxyRecordRequestBody struct {
	Record []string `form:"record,omitempty" json:"record,omitempty" xml:"record,omitempty"`
	// номер строки в исходном файле
	LineNumber *int `json:"line_number"`
}

// TargetUserRecordRequestBody is used to define fields on request body types.
type TargetUserRecordRequestBody struct {
	Record []string `form:"record,omitempty" json:"record,omitempty" xml:"record,omitempty"`
	// номер строки в исходном файле
	LineNumber *int `json:"line_number"`
}

// NewUploadFileResponseBody builds the HTTP response body from the result of
// the "upload file" endpoint of the "tasks_service" service.
func NewUploadFileResponseBody(res []*tasksservice.UploadError) UploadFileResponseBody {
	body := make([]*UploadErrorResponse, len(res))
	for i, val := range res {
		body[i] = marshalTasksserviceUploadErrorToUploadErrorResponse(val)
	}
	return body
}

// NewCreateTaskDraftPayload builds a tasks_service service create task draft
// endpoint payload.
func NewCreateTaskDraftPayload(body *CreateTaskDraftRequestBody, token string) *tasksservice.CreateTaskDraftPayload {
	v := &tasksservice.CreateTaskDraftPayload{
		Title:        *body.Title,
		TextTemplate: *body.TextTemplate,
		PostImage:    *body.PostImage,
	}
	v.Token = token

	return v
}

// NewUploadFilePayload builds a tasks_service service upload file endpoint
// payload.
func NewUploadFilePayload(body *UploadFileRequestBody, taskID string, token string) *tasksservice.UploadFilePayload {
	v := &tasksservice.UploadFilePayload{}
	v.Bots = make([]*tasksservice.BotAccountRecord, len(body.Bots))
	for i, val := range body.Bots {
		v.Bots[i] = unmarshalBotAccountRecordRequestBodyToTasksserviceBotAccountRecord(val)
	}
	v.Proxies = make([]*tasksservice.ProxyRecord, len(body.Proxies))
	for i, val := range body.Proxies {
		v.Proxies[i] = unmarshalProxyRecordRequestBodyToTasksserviceProxyRecord(val)
	}
	v.Targets = make([]*tasksservice.TargetUserRecord, len(body.Targets))
	for i, val := range body.Targets {
		v.Targets[i] = unmarshalTargetUserRecordRequestBodyToTasksserviceTargetUserRecord(val)
	}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewStartTaskPayload builds a tasks_service service start task endpoint
// payload.
func NewStartTaskPayload(taskID string, token string) *tasksservice.StartTaskPayload {
	v := &tasksservice.StartTaskPayload{}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewStopTaskPayload builds a tasks_service service stop task endpoint payload.
func NewStopTaskPayload(taskID string, token string) *tasksservice.StopTaskPayload {
	v := &tasksservice.StopTaskPayload{}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewGetTaskPayload builds a tasks_service service get task endpoint payload.
func NewGetTaskPayload(taskID string, token string) *tasksservice.GetTaskPayload {
	v := &tasksservice.GetTaskPayload{}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewListTasksPayload builds a tasks_service service list tasks endpoint
// payload.
func NewListTasksPayload(token string) *tasksservice.ListTasksPayload {
	v := &tasksservice.ListTasksPayload{}
	v.Token = token

	return v
}

// ValidateCreateTaskDraftRequestBody runs the validations defined on Create
// Task DraftRequestBody
func ValidateCreateTaskDraftRequestBody(body *CreateTaskDraftRequestBody) (err error) {
	if body.Title == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("title", "body"))
	}
	if body.TextTemplate == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("text_template", "body"))
	}
	if body.PostImage == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("post_image", "body"))
	}
	return
}

// ValidateUploadFileRequestBody runs the validations defined on Upload
// FileRequestBody
func ValidateUploadFileRequestBody(body *UploadFileRequestBody) (err error) {
	if body.Bots == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bots", "body"))
	}
	if body.Proxies == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("proxies", "body"))
	}
	if body.Targets == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("targets", "body"))
	}
	for _, e := range body.Bots {
		if e != nil {
			if err2 := ValidateBotAccountRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.Proxies {
		if e != nil {
			if err2 := ValidateProxyRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.Targets {
		if e != nil {
			if err2 := ValidateTargetUserRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateBotAccountRecordRequestBody runs the validations defined on
// BotAccountRecordRequestBody
func ValidateBotAccountRecordRequestBody(body *BotAccountRecordRequestBody) (err error) {
	if body.Record == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("record", "body"))
	}
	if body.LineNumber == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("line_number", "body"))
	}
	return
}

// ValidateProxyRecordRequestBody runs the validations defined on
// ProxyRecordRequestBody
func ValidateProxyRecordRequestBody(body *ProxyRecordRequestBody) (err error) {
	if body.Record == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("record", "body"))
	}
	if body.LineNumber == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("line_number", "body"))
	}
	return
}

// ValidateTargetUserRecordRequestBody runs the validations defined on
// TargetUserRecordRequestBody
func ValidateTargetUserRecordRequestBody(body *TargetUserRecordRequestBody) (err error) {
	if body.Record == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("record", "body"))
	}
	if body.LineNumber == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("line_number", "body"))
	}
	return
}
