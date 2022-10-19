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
	// имена аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
	// имена для аккаунтов-ботов
	BotNames []string `json:"bot_names"`
	// фамилии для аккаунтов-ботов
	BotLastNames []string `json:"bot_last_names"`
	// аватарки для ботов
	BotImages []string `json:"bot_images"`
	// список фотографий для постов
	PostImages []string `json:"post_images"`
}

// UpdateTaskRequestBody is the type of the "tasks_service" service "update
// task" endpoint HTTP request body.
type UpdateTaskRequestBody struct {
	// название задачи
	Title *string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
	// шаблон для подписи под постом
	TextTemplate *string `json:"text_template"`
	// фотография для постов
	PostImages []string `json:"post_images"`
}

// UploadFilesRequestBody is the type of the "tasks_service" service "upload
// files" endpoint HTTP request body.
type UploadFilesRequestBody struct {
	Filenames *TaskFileNamesRequestBody `form:"filenames,omitempty" json:"filenames,omitempty" xml:"filenames,omitempty"`
	// список ботов
	Bots []*BotAccountRecordRequestBody `form:"bots,omitempty" json:"bots,omitempty" xml:"bots,omitempty"`
	// список проксей для использования
	ResidentialProxies []*ProxyRecordRequestBody `form:"residential_proxies,omitempty" json:"residential_proxies,omitempty" xml:"residential_proxies,omitempty"`
	// список дешёвых проксей для загрузки фото
	CheapProxies []*ProxyRecordRequestBody `form:"cheap_proxies,omitempty" json:"cheap_proxies,omitempty" xml:"cheap_proxies,omitempty"`
	// список аккаунтов, которым показать надо рекламу
	Targets []*TargetUserRecordRequestBody `form:"targets,omitempty" json:"targets,omitempty" xml:"targets,omitempty"`
}

// UpdateTaskOKResponseBody is the type of the "tasks_service" service "update
// task" endpoint HTTP response body.
type UpdateTaskOKResponseBody struct {
	ID string `form:"id" json:"id" xml:"id"`
	// описание под постом
	TextTemplate string `json:"text_template"`
	// список base64 строк картинок
	Images []string `form:"images" json:"images" xml:"images"`
	Status int      `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество прокси в задаче
	ProxiesNum int `json:"proxies_num"`
	// количество целевых пользователей в задаче
	TargetsNum int `json:"targets_num"`
	// название файла, из которого брали ботов
	BotsFilename *string `json:"bots_filename"`
	// название файла, из которого брали резидентские прокси
	ResidentialProxiesFilename *string `json:"residential_proxies_filename"`
	// название файла, из которого брали дешёвые прокси
	CheapProxiesFilename *string `json:"cheap_proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename *string `json:"targets_filename"`
}

// UploadFilesOKResponseBody is the type of the "tasks_service" service "upload
// files" endpoint HTTP response body.
type UploadFilesOKResponseBody struct {
	// ошибки, которые возникли при загрузке файлов
	UploadErrors []*UploadErrorResponseBody `json:"upload_errors"`
	Status       int                        `form:"status" json:"status" xml:"status"`
}

// AssignProxiesOKResponseBody is the type of the "tasks_service" service
// "assign proxies" endpoint HTTP response body.
type AssignProxiesOKResponseBody struct {
	// количество аккаунтов с проксями, которые будут использованы для текущей
	// задачи
	BotsNumber int `json:"bots_number"`
	Status     int `form:"status" json:"status" xml:"status"`
	// id задачи
	TaskID string `json:"task_id"`
}

// StartTaskOKResponseBody is the type of the "tasks_service" service "start
// task" endpoint HTTP response body.
type StartTaskOKResponseBody struct {
	Status int `form:"status" json:"status" xml:"status"`
	// id задачи
	TaskID string `json:"task_id"`
}

// StopTaskOKResponseBody is the type of the "tasks_service" service "stop
// task" endpoint HTTP response body.
type StopTaskOKResponseBody struct {
	Status int `form:"status" json:"status" xml:"status"`
	// id задачи
	TaskID string `json:"task_id"`
}

// GetTaskOKResponseBody is the type of the "tasks_service" service "get task"
// endpoint HTTP response body.
type GetTaskOKResponseBody struct {
	ID string `form:"id" json:"id" xml:"id"`
	// описание под постом
	TextTemplate string `json:"text_template"`
	// список base64 строк картинок
	Images []string `form:"images" json:"images" xml:"images"`
	Status int      `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество прокси в задаче
	ProxiesNum int `json:"proxies_num"`
	// количество целевых пользователей в задаче
	TargetsNum int `json:"targets_num"`
	// название файла, из которого брали ботов
	BotsFilename *string `json:"bots_filename"`
	// название файла, из которого брали резидентские прокси
	ResidentialProxiesFilename *string `json:"residential_proxies_filename"`
	// название файла, из которого брали дешёвые прокси
	CheapProxiesFilename *string `json:"cheap_proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename *string `json:"targets_filename"`
}

// GetProgressResponseBody is the type of the "tasks_service" service "get
// progress" endpoint HTTP response body.
type GetProgressResponseBody []*BotsProgressResponse

// ListTasksResponseBody is the type of the "tasks_service" service "list
// tasks" endpoint HTTP response body.
type ListTasksResponseBody []*TaskResponse

// UploadErrorResponseBody is used to define fields on response body types.
type UploadErrorResponseBody struct {
	// 1 - список ботов
	// 2 - список прокси
	// 3 - список получателей рекламы
	Type int `form:"type" json:"type" xml:"type"`
	Line int `form:"line" json:"line" xml:"line"`
	// номер порта
	Input  string `form:"input" json:"input" xml:"input"`
	Reason string `form:"reason" json:"reason" xml:"reason"`
}

// BotsProgressResponse is used to define fields on response body types.
type BotsProgressResponse struct {
	// имя пользователя бота
	UserName string `json:"user_name"`
	// количество выложенных постов
	PostsCount int `json:"posts_count"`
	// текущий статус бота, будут ли выкладываться посты
	Status int `form:"status" json:"status" xml:"status"`
}

// TaskResponse is used to define fields on response body types.
type TaskResponse struct {
	ID string `form:"id" json:"id" xml:"id"`
	// описание под постом
	TextTemplate string `json:"text_template"`
	// список base64 строк картинок
	Images []string `form:"images" json:"images" xml:"images"`
	Status int      `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество прокси в задаче
	ProxiesNum int `json:"proxies_num"`
	// количество целевых пользователей в задаче
	TargetsNum int `json:"targets_num"`
	// название файла, из которого брали ботов
	BotsFilename *string `json:"bots_filename"`
	// название файла, из которого брали резидентские прокси
	ResidentialProxiesFilename *string `json:"residential_proxies_filename"`
	// название файла, из которого брали дешёвые прокси
	CheapProxiesFilename *string `json:"cheap_proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename *string `json:"targets_filename"`
}

// TaskFileNamesRequestBody is used to define fields on request body types.
type TaskFileNamesRequestBody struct {
	// название файла, из которого брали ботов
	BotsFilename *string `json:"bots_filename"`
	// название файла, из которого брали резидентские прокси
	ResidentialProxiesFilename *string `json:"residential_proxies_filename"`
	// название файла, из которого брали дешёвые прокси
	CheapProxiesFilename *string `json:"cheap_proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename *string `json:"targets_filename"`
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

// NewUpdateTaskOKResponseBody builds the HTTP response body from the result of
// the "update task" endpoint of the "tasks_service" service.
func NewUpdateTaskOKResponseBody(res *tasksservice.Task) *UpdateTaskOKResponseBody {
	body := &UpdateTaskOKResponseBody{
		ID:                         res.ID,
		TextTemplate:               res.TextTemplate,
		Status:                     int(res.Status),
		Title:                      res.Title,
		BotsNum:                    res.BotsNum,
		ProxiesNum:                 res.ProxiesNum,
		TargetsNum:                 res.TargetsNum,
		BotsFilename:               res.BotsFilename,
		ResidentialProxiesFilename: res.ResidentialProxiesFilename,
		CheapProxiesFilename:       res.CheapProxiesFilename,
		TargetsFilename:            res.TargetsFilename,
	}
	if res.Images != nil {
		body.Images = make([]string, len(res.Images))
		for i, val := range res.Images {
			body.Images[i] = val
		}
	}
	return body
}

// NewUploadFilesOKResponseBody builds the HTTP response body from the result
// of the "upload files" endpoint of the "tasks_service" service.
func NewUploadFilesOKResponseBody(res *tasksservice.UploadFilesResult) *UploadFilesOKResponseBody {
	body := &UploadFilesOKResponseBody{
		Status: int(res.Status),
	}
	if res.UploadErrors != nil {
		body.UploadErrors = make([]*UploadErrorResponseBody, len(res.UploadErrors))
		for i, val := range res.UploadErrors {
			body.UploadErrors[i] = marshalTasksserviceUploadErrorToUploadErrorResponseBody(val)
		}
	}
	return body
}

// NewAssignProxiesOKResponseBody builds the HTTP response body from the result
// of the "assign proxies" endpoint of the "tasks_service" service.
func NewAssignProxiesOKResponseBody(res *tasksservice.AssignProxiesResult) *AssignProxiesOKResponseBody {
	body := &AssignProxiesOKResponseBody{
		BotsNumber: res.BotsNumber,
		Status:     int(res.Status),
		TaskID:     res.TaskID,
	}
	return body
}

// NewStartTaskOKResponseBody builds the HTTP response body from the result of
// the "start task" endpoint of the "tasks_service" service.
func NewStartTaskOKResponseBody(res *tasksservice.StartTaskResult) *StartTaskOKResponseBody {
	body := &StartTaskOKResponseBody{
		Status: int(res.Status),
		TaskID: res.TaskID,
	}
	return body
}

// NewStopTaskOKResponseBody builds the HTTP response body from the result of
// the "stop task" endpoint of the "tasks_service" service.
func NewStopTaskOKResponseBody(res *tasksservice.StopTaskResult) *StopTaskOKResponseBody {
	body := &StopTaskOKResponseBody{
		Status: int(res.Status),
		TaskID: res.TaskID,
	}
	return body
}

// NewGetTaskOKResponseBody builds the HTTP response body from the result of
// the "get task" endpoint of the "tasks_service" service.
func NewGetTaskOKResponseBody(res *tasksservice.Task) *GetTaskOKResponseBody {
	body := &GetTaskOKResponseBody{
		ID:                         res.ID,
		TextTemplate:               res.TextTemplate,
		Status:                     int(res.Status),
		Title:                      res.Title,
		BotsNum:                    res.BotsNum,
		ProxiesNum:                 res.ProxiesNum,
		TargetsNum:                 res.TargetsNum,
		BotsFilename:               res.BotsFilename,
		ResidentialProxiesFilename: res.ResidentialProxiesFilename,
		CheapProxiesFilename:       res.CheapProxiesFilename,
		TargetsFilename:            res.TargetsFilename,
	}
	if res.Images != nil {
		body.Images = make([]string, len(res.Images))
		for i, val := range res.Images {
			body.Images[i] = val
		}
	}
	return body
}

// NewGetProgressResponseBody builds the HTTP response body from the result of
// the "get progress" endpoint of the "tasks_service" service.
func NewGetProgressResponseBody(res []*tasksservice.BotsProgress) GetProgressResponseBody {
	body := make([]*BotsProgressResponse, len(res))
	for i, val := range res {
		body[i] = marshalTasksserviceBotsProgressToBotsProgressResponse(val)
	}
	return body
}

// NewListTasksResponseBody builds the HTTP response body from the result of
// the "list tasks" endpoint of the "tasks_service" service.
func NewListTasksResponseBody(res []*tasksservice.Task) ListTasksResponseBody {
	body := make([]*TaskResponse, len(res))
	for i, val := range res {
		body[i] = marshalTasksserviceTaskToTaskResponse(val)
	}
	return body
}

// NewCreateTaskDraftPayload builds a tasks_service service create task draft
// endpoint payload.
func NewCreateTaskDraftPayload(body *CreateTaskDraftRequestBody, token string) *tasksservice.CreateTaskDraftPayload {
	v := &tasksservice.CreateTaskDraftPayload{
		Title:        *body.Title,
		TextTemplate: *body.TextTemplate,
	}
	v.LandingAccounts = make([]string, len(body.LandingAccounts))
	for i, val := range body.LandingAccounts {
		v.LandingAccounts[i] = val
	}
	if body.BotNames != nil {
		v.BotNames = make([]string, len(body.BotNames))
		for i, val := range body.BotNames {
			v.BotNames[i] = val
		}
	}
	if body.BotLastNames != nil {
		v.BotLastNames = make([]string, len(body.BotLastNames))
		for i, val := range body.BotLastNames {
			v.BotLastNames[i] = val
		}
	}
	if body.BotImages != nil {
		v.BotImages = make([]string, len(body.BotImages))
		for i, val := range body.BotImages {
			v.BotImages[i] = val
		}
	}
	v.PostImages = make([]string, len(body.PostImages))
	for i, val := range body.PostImages {
		v.PostImages[i] = val
	}
	v.Token = token

	return v
}

// NewUpdateTaskPayload builds a tasks_service service update task endpoint
// payload.
func NewUpdateTaskPayload(body *UpdateTaskRequestBody, taskID string, token string) *tasksservice.UpdateTaskPayload {
	v := &tasksservice.UpdateTaskPayload{
		Title:        body.Title,
		TextTemplate: body.TextTemplate,
	}
	if body.PostImages != nil {
		v.PostImages = make([]string, len(body.PostImages))
		for i, val := range body.PostImages {
			v.PostImages[i] = val
		}
	}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewUploadFilesPayload builds a tasks_service service upload files endpoint
// payload.
func NewUploadFilesPayload(body *UploadFilesRequestBody, taskID string, token string) *tasksservice.UploadFilesPayload {
	v := &tasksservice.UploadFilesPayload{}
	v.Filenames = unmarshalTaskFileNamesRequestBodyToTasksserviceTaskFileNames(body.Filenames)
	v.Bots = make([]*tasksservice.BotAccountRecord, len(body.Bots))
	for i, val := range body.Bots {
		v.Bots[i] = unmarshalBotAccountRecordRequestBodyToTasksserviceBotAccountRecord(val)
	}
	v.ResidentialProxies = make([]*tasksservice.ProxyRecord, len(body.ResidentialProxies))
	for i, val := range body.ResidentialProxies {
		v.ResidentialProxies[i] = unmarshalProxyRecordRequestBodyToTasksserviceProxyRecord(val)
	}
	v.CheapProxies = make([]*tasksservice.ProxyRecord, len(body.CheapProxies))
	for i, val := range body.CheapProxies {
		v.CheapProxies[i] = unmarshalProxyRecordRequestBodyToTasksserviceProxyRecord(val)
	}
	v.Targets = make([]*tasksservice.TargetUserRecord, len(body.Targets))
	for i, val := range body.Targets {
		v.Targets[i] = unmarshalTargetUserRecordRequestBodyToTasksserviceTargetUserRecord(val)
	}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewAssignProxiesPayload builds a tasks_service service assign proxies
// endpoint payload.
func NewAssignProxiesPayload(taskID string, token string) *tasksservice.AssignProxiesPayload {
	v := &tasksservice.AssignProxiesPayload{}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewForceDeletePayload builds a tasks_service service force delete endpoint
// payload.
func NewForceDeletePayload(taskID string, token string) *tasksservice.ForceDeletePayload {
	v := &tasksservice.ForceDeletePayload{}
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

// NewGetProgressPayload builds a tasks_service service get progress endpoint
// payload.
func NewGetProgressPayload(taskID string, token string) *tasksservice.GetProgressPayload {
	v := &tasksservice.GetProgressPayload{}
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
	if body.PostImages == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("post_images", "body"))
	}
	if body.LandingAccounts == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("landing_accounts", "body"))
	}
	return
}

// ValidateUploadFilesRequestBody runs the validations defined on Upload
// FilesRequestBody
func ValidateUploadFilesRequestBody(body *UploadFilesRequestBody) (err error) {
	if body.Bots == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bots", "body"))
	}
	if body.ResidentialProxies == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("residential_proxies", "body"))
	}
	if body.CheapProxies == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("cheap_proxies", "body"))
	}
	if body.Targets == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("targets", "body"))
	}
	if body.Filenames == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("filenames", "body"))
	}
	if body.Filenames != nil {
		if err2 := ValidateTaskFileNamesRequestBody(body.Filenames); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	for _, e := range body.Bots {
		if e != nil {
			if err2 := ValidateBotAccountRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.ResidentialProxies {
		if e != nil {
			if err2 := ValidateProxyRecordRequestBody(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	for _, e := range body.CheapProxies {
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

// ValidateTaskFileNamesRequestBody runs the validations defined on
// TaskFileNamesRequestBody
func ValidateTaskFileNamesRequestBody(body *TaskFileNamesRequestBody) (err error) {
	if body.BotsFilename == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("bots_filename", "body"))
	}
	if body.ResidentialProxiesFilename == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("residential_proxies_filename", "body"))
	}
	if body.CheapProxiesFilename == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("cheap_proxies_filename", "body"))
	}
	if body.TargetsFilename == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("targets_filename", "body"))
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
	if len(body.Record) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, true))
	}
	if len(body.Record) > 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, false))
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
	if len(body.Record) < 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, true))
	}
	if len(body.Record) > 4 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 4, false))
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
	if len(body.Record) < 2 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 2, true))
	}
	if len(body.Record) > 2 {
		err = goa.MergeErrors(err, goa.InvalidLengthError("body.record", body.Record, len(body.Record), 2, false))
	}
	return
}
