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
	// ссылки для описания у ботов
	BotUrls []string `json:"bot_images"`
	// список фотографий для постов
	PostImages []string `json:"post_images"`
	Type       *int     `form:"type,omitempty" json:"type,omitempty" xml:"type,omitempty"`
}

// UpdateTaskRequestBody is the type of the "tasks_service" service "update
// task" endpoint HTTP request body.
type UpdateTaskRequestBody struct {
	// описание под постом
	TextTemplate *string `json:"text_template"`
	// имена аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
	// имена для аккаунтов-ботов
	BotNames []string `json:"bot_names"`
	// фамилии для аккаунтов-ботов
	BotLastNames []string `json:"bot_last_names"`
	// ссылки для описания у ботов
	BotUrls []string `json:"bot_urls"`
	// название задачи
	Title *string `form:"title,omitempty" json:"title,omitempty" xml:"title,omitempty"`
	// нужно ли подписываться на аккаунты
	FollowTargets *bool `json:"follow_targets"`
	// делать отметки на фотографии
	NeedPhotoTags *bool `json:"need_photo_tags"`
	// задержка между постами
	PerPostSleepSeconds *uint `json:"per_post_sleep_seconds"`
	// задержка между загрузкой фотографии и проставлением отметок (в секундах)
	PhotoTagsDelaySeconds *uint `json:"photo_tags_delay_seconds"`
	// количество постов для каждого бота
	PostsPerBot *uint `json:"posts_per_bot"`
	// количество постов с отметками на фото для каждого бота
	PhotoTagsPostsPerBot *uint `json:"photo_tags_posts_per_bot"`
	// количество упоминаний под каждым постом
	TargetsPerPost *uint `json:"targets_per_post"`
	// количество упоминаний на фото у каждого поста
	PhotoTargetsPerPost *uint `json:"photo_targets_per_post"`
	// список base64 строк картинок
	PostImages []string `json:"post_images"`
	// аватарки для ботов
	BotImages []string `json:"bot_images"`
}

// UploadVideoRequestBody is the type of the "tasks_service" service "upload
// video" endpoint HTTP request body.
type UploadVideoRequestBody struct {
	// не нужно присылать руками, подставится автоматом
	Filename *string `form:"filename,omitempty" json:"filename,omitempty" xml:"filename,omitempty"`
	Video    []byte  `form:"video,omitempty" json:"video,omitempty" xml:"video,omitempty"`
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

// PartialStartTaskRequestBody is the type of the "tasks_service" service
// "partial start task" endpoint HTTP request body.
type PartialStartTaskRequestBody struct {
	// список имен ботов, которых нужно запустить
	Usernames []string `form:"usernames,omitempty" json:"usernames,omitempty" xml:"usernames,omitempty"`
}

// UpdateTaskOKResponseBody is the type of the "tasks_service" service "update
// task" endpoint HTTP response body.
type UpdateTaskOKResponseBody struct {
	ID   string `form:"id" json:"id" xml:"id"`
	Type int    `form:"type" json:"type" xml:"type"`
	// описание под постом
	TextTemplate string `json:"text_template"`
	// имена аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
	// имена для аккаунтов-ботов
	BotNames []string `json:"bot_names"`
	// фамилии для аккаунтов-ботов
	BotLastNames []string `json:"bot_last_names"`
	// ссылки для описания у ботов
	BotUrls []string `json:"bot_urls"`
	Status  int      `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество резидентских прокси в задаче
	ResidentialProxiesNum int `json:"residential_proxies_num"`
	// количество дешёвых прокси в задаче
	CheapProxiesNum int `json:"cheap_proxies_num"`
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
	// название файла с видео, если тип задачи - рилсы
	VideoFilename *string `json:"video_filename"`
	// нужно ли подписываться на аккаунты
	FollowTargets bool `json:"follow_targets"`
	// делать отметки на фотографии
	NeedPhotoTags bool `json:"need_photo_tags"`
	// задержка между постами
	PerPostSleepSeconds uint `json:"per_post_sleep_seconds"`
	// задержка между загрузкой фотографии и проставлением отметок (в секундах)
	PhotoTagsDelaySeconds uint `json:"photo_tags_delay_seconds"`
	// количество постов для каждого бота
	PostsPerBot uint `json:"posts_per_bot"`
	// количество постов с отметками на фото для каждого бота
	PhotoTagsPostsPerBot uint `json:"photo_tags_posts_per_bot"`
	// количество упоминаний под каждым постом
	TargetsPerPost uint `json:"targets_per_post"`
	// количество упоминаний на фото у каждого поста
	PhotoTargetsPerPost uint `json:"photo_targets_per_post"`
	// список base64 строк картинок
	PostImages []string `json:"post_images"`
	// аватарки для ботов
	BotImages []string `json:"bot_images"`
}

// UploadVideoOKResponseBody is the type of the "tasks_service" service "upload
// video" endpoint HTTP response body.
type UploadVideoOKResponseBody struct {
	Status int `form:"status" json:"status" xml:"status"`
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
	// имена живых аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
}

// PartialStartTaskOKResponseBody is the type of the "tasks_service" service
// "partial start task" endpoint HTTP response body.
type PartialStartTaskOKResponseBody struct {
	// id задачи
	TaskID string `json:"task_id"`
	// список успешных имен ботов
	Succeeded []string `form:"succeeded" json:"succeeded" xml:"succeeded"`
	// имена живых аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
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
	ID   string `form:"id" json:"id" xml:"id"`
	Type int    `form:"type" json:"type" xml:"type"`
	// описание под постом
	TextTemplate string `json:"text_template"`
	// имена аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
	// имена для аккаунтов-ботов
	BotNames []string `json:"bot_names"`
	// фамилии для аккаунтов-ботов
	BotLastNames []string `json:"bot_last_names"`
	// ссылки для описания у ботов
	BotUrls []string `json:"bot_urls"`
	Status  int      `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество резидентских прокси в задаче
	ResidentialProxiesNum int `json:"residential_proxies_num"`
	// количество дешёвых прокси в задаче
	CheapProxiesNum int `json:"cheap_proxies_num"`
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
	// название файла с видео, если тип задачи - рилсы
	VideoFilename *string `json:"video_filename"`
	// нужно ли подписываться на аккаунты
	FollowTargets bool `json:"follow_targets"`
	// делать отметки на фотографии
	NeedPhotoTags bool `json:"need_photo_tags"`
	// задержка между постами
	PerPostSleepSeconds uint `json:"per_post_sleep_seconds"`
	// задержка между загрузкой фотографии и проставлением отметок (в секундах)
	PhotoTagsDelaySeconds uint `json:"photo_tags_delay_seconds"`
	// количество постов для каждого бота
	PostsPerBot uint `json:"posts_per_bot"`
	// количество постов с отметками на фото для каждого бота
	PhotoTagsPostsPerBot uint `json:"photo_tags_posts_per_bot"`
	// количество упоминаний под каждым постом
	TargetsPerPost uint `json:"targets_per_post"`
	// количество упоминаний на фото у каждого поста
	PhotoTargetsPerPost uint `json:"photo_targets_per_post"`
	// список base64 строк картинок
	PostImages []string `json:"post_images"`
	// аватарки для ботов
	BotImages []string `json:"bot_images"`
}

// GetProgressOKResponseBody is the type of the "tasks_service" service "get
// progress" endpoint HTTP response body.
type GetProgressOKResponseBody struct {
	// результат работы по каждому боту
	BotsProgresses []*BotsProgressResponseBody `json:"bots_progresses"`
	// количество аккаунтов, которых упомянули в постах
	TargetsNotified int `json:"targets_notified"`
	// количество аккаунтов, которых упомянули в постах на фото
	PhotoTargetsNotified int `json:"photo_targets_notified"`
	// количество аккаунтов, которых не получилось упомянуть, при перезапуске
	// задачи будут использованы заново
	TargetsFailed int `json:"targets_failed"`
	// количество аккаунтов, которых не выбрали для постов
	TargetsWaiting int `json:"targets_waiting"`
	// закончена ли задача
	Done bool `form:"done" json:"done" xml:"done"`
}

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

// BotsProgressResponseBody is used to define fields on response body types.
type BotsProgressResponseBody struct {
	// имя пользователя бота
	Username string `form:"username" json:"username" xml:"username"`
	// количество выложенных постов
	PostsCount int32 `json:"posts_count"`
	// текущий статус бота, будут ли выкладываться посты
	Status int32 `form:"status" json:"status" xml:"status"`
	// количество аккаунтов, которых упомянули в постах
	DescriptionTargetsNotified int32 `json:"description_targets_notified"`
	// количество аккаунтов, которых упомянули в постах на фото
	PhotoTargetsNotified int32 `json:"photo_targets_notified"`
	// номер бота в загруженном файле
	FileOrder int32 `json:"file_order"`
}

// TaskResponse is used to define fields on response body types.
type TaskResponse struct {
	ID   string `form:"id" json:"id" xml:"id"`
	Type int    `form:"type" json:"type" xml:"type"`
	// описание под постом
	TextTemplate string `json:"text_template"`
	// имена аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
	// имена для аккаунтов-ботов
	BotNames []string `json:"bot_names"`
	// фамилии для аккаунтов-ботов
	BotLastNames []string `json:"bot_last_names"`
	// ссылки для описания у ботов
	BotUrls []string `json:"bot_urls"`
	Status  int      `form:"status" json:"status" xml:"status"`
	// название задачи
	Title string `form:"title" json:"title" xml:"title"`
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество резидентских прокси в задаче
	ResidentialProxiesNum int `json:"residential_proxies_num"`
	// количество дешёвых прокси в задаче
	CheapProxiesNum int `json:"cheap_proxies_num"`
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
	// название файла с видео, если тип задачи - рилсы
	VideoFilename *string `json:"video_filename"`
	// нужно ли подписываться на аккаунты
	FollowTargets bool `json:"follow_targets"`
	// делать отметки на фотографии
	NeedPhotoTags bool `json:"need_photo_tags"`
	// задержка между постами
	PerPostSleepSeconds uint `json:"per_post_sleep_seconds"`
	// задержка между загрузкой фотографии и проставлением отметок (в секундах)
	PhotoTagsDelaySeconds uint `json:"photo_tags_delay_seconds"`
	// количество постов для каждого бота
	PostsPerBot uint `json:"posts_per_bot"`
	// количество постов с отметками на фото для каждого бота
	PhotoTagsPostsPerBot uint `json:"photo_tags_posts_per_bot"`
	// количество упоминаний под каждым постом
	TargetsPerPost uint `json:"targets_per_post"`
	// количество упоминаний на фото у каждого поста
	PhotoTargetsPerPost uint `json:"photo_targets_per_post"`
	// список base64 строк картинок
	PostImages []string `json:"post_images"`
	// аватарки для ботов
	BotImages []string `json:"bot_images"`
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
		Type:                       int(res.Type),
		TextTemplate:               res.TextTemplate,
		Status:                     int(res.Status),
		Title:                      res.Title,
		BotsNum:                    res.BotsNum,
		ResidentialProxiesNum:      res.ResidentialProxiesNum,
		CheapProxiesNum:            res.CheapProxiesNum,
		TargetsNum:                 res.TargetsNum,
		BotsFilename:               res.BotsFilename,
		ResidentialProxiesFilename: res.ResidentialProxiesFilename,
		CheapProxiesFilename:       res.CheapProxiesFilename,
		TargetsFilename:            res.TargetsFilename,
		VideoFilename:              res.VideoFilename,
		FollowTargets:              res.FollowTargets,
		NeedPhotoTags:              res.NeedPhotoTags,
		PerPostSleepSeconds:        res.PerPostSleepSeconds,
		PhotoTagsDelaySeconds:      res.PhotoTagsDelaySeconds,
		PostsPerBot:                res.PostsPerBot,
		PhotoTagsPostsPerBot:       res.PhotoTagsPostsPerBot,
		TargetsPerPost:             res.TargetsPerPost,
		PhotoTargetsPerPost:        res.PhotoTargetsPerPost,
	}
	if res.LandingAccounts != nil {
		body.LandingAccounts = make([]string, len(res.LandingAccounts))
		for i, val := range res.LandingAccounts {
			body.LandingAccounts[i] = val
		}
	}
	if res.BotNames != nil {
		body.BotNames = make([]string, len(res.BotNames))
		for i, val := range res.BotNames {
			body.BotNames[i] = val
		}
	}
	if res.BotLastNames != nil {
		body.BotLastNames = make([]string, len(res.BotLastNames))
		for i, val := range res.BotLastNames {
			body.BotLastNames[i] = val
		}
	}
	if res.BotUrls != nil {
		body.BotUrls = make([]string, len(res.BotUrls))
		for i, val := range res.BotUrls {
			body.BotUrls[i] = val
		}
	}
	if res.PostImages != nil {
		body.PostImages = make([]string, len(res.PostImages))
		for i, val := range res.PostImages {
			body.PostImages[i] = val
		}
	}
	if res.BotImages != nil {
		body.BotImages = make([]string, len(res.BotImages))
		for i, val := range res.BotImages {
			body.BotImages[i] = val
		}
	}
	return body
}

// NewUploadVideoOKResponseBody builds the HTTP response body from the result
// of the "upload video" endpoint of the "tasks_service" service.
func NewUploadVideoOKResponseBody(res *tasksservice.UploadVideoResult) *UploadVideoOKResponseBody {
	body := &UploadVideoOKResponseBody{
		Status: int(res.Status),
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
	if res.LandingAccounts != nil {
		body.LandingAccounts = make([]string, len(res.LandingAccounts))
		for i, val := range res.LandingAccounts {
			body.LandingAccounts[i] = val
		}
	}
	return body
}

// NewPartialStartTaskOKResponseBody builds the HTTP response body from the
// result of the "partial start task" endpoint of the "tasks_service" service.
func NewPartialStartTaskOKResponseBody(res *tasksservice.PartialStartTaskResult) *PartialStartTaskOKResponseBody {
	body := &PartialStartTaskOKResponseBody{
		TaskID: res.TaskID,
	}
	if res.Succeeded != nil {
		body.Succeeded = make([]string, len(res.Succeeded))
		for i, val := range res.Succeeded {
			body.Succeeded[i] = val
		}
	}
	if res.LandingAccounts != nil {
		body.LandingAccounts = make([]string, len(res.LandingAccounts))
		for i, val := range res.LandingAccounts {
			body.LandingAccounts[i] = val
		}
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
		Type:                       int(res.Type),
		TextTemplate:               res.TextTemplate,
		Status:                     int(res.Status),
		Title:                      res.Title,
		BotsNum:                    res.BotsNum,
		ResidentialProxiesNum:      res.ResidentialProxiesNum,
		CheapProxiesNum:            res.CheapProxiesNum,
		TargetsNum:                 res.TargetsNum,
		BotsFilename:               res.BotsFilename,
		ResidentialProxiesFilename: res.ResidentialProxiesFilename,
		CheapProxiesFilename:       res.CheapProxiesFilename,
		TargetsFilename:            res.TargetsFilename,
		VideoFilename:              res.VideoFilename,
		FollowTargets:              res.FollowTargets,
		NeedPhotoTags:              res.NeedPhotoTags,
		PerPostSleepSeconds:        res.PerPostSleepSeconds,
		PhotoTagsDelaySeconds:      res.PhotoTagsDelaySeconds,
		PostsPerBot:                res.PostsPerBot,
		PhotoTagsPostsPerBot:       res.PhotoTagsPostsPerBot,
		TargetsPerPost:             res.TargetsPerPost,
		PhotoTargetsPerPost:        res.PhotoTargetsPerPost,
	}
	if res.LandingAccounts != nil {
		body.LandingAccounts = make([]string, len(res.LandingAccounts))
		for i, val := range res.LandingAccounts {
			body.LandingAccounts[i] = val
		}
	}
	if res.BotNames != nil {
		body.BotNames = make([]string, len(res.BotNames))
		for i, val := range res.BotNames {
			body.BotNames[i] = val
		}
	}
	if res.BotLastNames != nil {
		body.BotLastNames = make([]string, len(res.BotLastNames))
		for i, val := range res.BotLastNames {
			body.BotLastNames[i] = val
		}
	}
	if res.BotUrls != nil {
		body.BotUrls = make([]string, len(res.BotUrls))
		for i, val := range res.BotUrls {
			body.BotUrls[i] = val
		}
	}
	if res.PostImages != nil {
		body.PostImages = make([]string, len(res.PostImages))
		for i, val := range res.PostImages {
			body.PostImages[i] = val
		}
	}
	if res.BotImages != nil {
		body.BotImages = make([]string, len(res.BotImages))
		for i, val := range res.BotImages {
			body.BotImages[i] = val
		}
	}
	return body
}

// NewGetProgressOKResponseBody builds the HTTP response body from the result
// of the "get progress" endpoint of the "tasks_service" service.
func NewGetProgressOKResponseBody(res *tasksservice.TaskProgress) *GetProgressOKResponseBody {
	body := &GetProgressOKResponseBody{
		TargetsNotified:      res.TargetsNotified,
		PhotoTargetsNotified: res.PhotoTargetsNotified,
		TargetsFailed:        res.TargetsFailed,
		TargetsWaiting:       res.TargetsWaiting,
		Done:                 res.Done,
	}
	if res.BotsProgresses != nil {
		body.BotsProgresses = make([]*BotsProgressResponseBody, len(res.BotsProgresses))
		for i, val := range res.BotsProgresses {
			body.BotsProgresses[i] = marshalTasksserviceBotsProgressToBotsProgressResponseBody(val)
		}
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
		Type:         tasksservice.TaskType(*body.Type),
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
	if body.BotUrls != nil {
		v.BotUrls = make([]string, len(body.BotUrls))
		for i, val := range body.BotUrls {
			v.BotUrls[i] = val
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
		TextTemplate:          body.TextTemplate,
		Title:                 body.Title,
		FollowTargets:         body.FollowTargets,
		NeedPhotoTags:         body.NeedPhotoTags,
		PerPostSleepSeconds:   body.PerPostSleepSeconds,
		PhotoTagsDelaySeconds: body.PhotoTagsDelaySeconds,
		PostsPerBot:           body.PostsPerBot,
		PhotoTagsPostsPerBot:  body.PhotoTagsPostsPerBot,
		TargetsPerPost:        body.TargetsPerPost,
		PhotoTargetsPerPost:   body.PhotoTargetsPerPost,
	}
	if body.LandingAccounts != nil {
		v.LandingAccounts = make([]string, len(body.LandingAccounts))
		for i, val := range body.LandingAccounts {
			v.LandingAccounts[i] = val
		}
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
	if body.BotUrls != nil {
		v.BotUrls = make([]string, len(body.BotUrls))
		for i, val := range body.BotUrls {
			v.BotUrls[i] = val
		}
	}
	if body.PostImages != nil {
		v.PostImages = make([]string, len(body.PostImages))
		for i, val := range body.PostImages {
			v.PostImages[i] = val
		}
	}
	if body.BotImages != nil {
		v.BotImages = make([]string, len(body.BotImages))
		for i, val := range body.BotImages {
			v.BotImages[i] = val
		}
	}
	v.TaskID = taskID
	v.Token = token

	return v
}

// NewUploadVideoPayload builds a tasks_service service upload video endpoint
// payload.
func NewUploadVideoPayload(body *UploadVideoRequestBody, taskID string, token string) *tasksservice.UploadVideoPayload {
	v := &tasksservice.UploadVideoPayload{
		Filename: body.Filename,
		Video:    body.Video,
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

// NewPartialStartTaskPayload builds a tasks_service service partial start task
// endpoint payload.
func NewPartialStartTaskPayload(body *PartialStartTaskRequestBody, taskID string, token string) *tasksservice.PartialStartTaskPayload {
	v := &tasksservice.PartialStartTaskPayload{}
	if body.Usernames != nil {
		v.Usernames = make([]string, len(body.Usernames))
		for i, val := range body.Usernames {
			v.Usernames[i] = val
		}
	}
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
func NewGetProgressPayload(taskID string, pageSize uint32, page uint32, sort string, sortDescending bool, token string) *tasksservice.GetProgressPayload {
	v := &tasksservice.GetProgressPayload{}
	v.TaskID = taskID
	v.PageSize = pageSize
	v.Page = page
	v.Sort = sort
	v.SortDescending = sortDescending
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

// NewDownloadTargetsPayload builds a tasks_service service download targets
// endpoint payload.
func NewDownloadTargetsPayload(taskID string, format int, token string) *tasksservice.DownloadTargetsPayload {
	v := &tasksservice.DownloadTargetsPayload{}
	v.TaskID = taskID
	v.Format = format
	v.Token = token

	return v
}

// NewDownloadBotsPayload builds a tasks_service service download bots endpoint
// payload.
func NewDownloadBotsPayload(taskID string, proxies bool, token string) *tasksservice.DownloadBotsPayload {
	v := &tasksservice.DownloadBotsPayload{}
	v.TaskID = taskID
	v.Proxies = proxies
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
	if body.Type == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("type", "body"))
	}
	if body.Type != nil {
		if !(*body.Type == 1 || *body.Type == 2) {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.type", *body.Type, []interface{}{1, 2}))
		}
	}
	return
}

// ValidateUploadVideoRequestBody runs the validations defined on Upload
// VideoRequestBody
func ValidateUploadVideoRequestBody(body *UploadVideoRequestBody) (err error) {
	if body.Video == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("video", "body"))
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
