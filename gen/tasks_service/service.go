// Code generated by goa v3.8.5, DO NOT EDIT.
//
// tasks_service service
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package tasksservice

import (
	"context"

	"goa.design/goa/v3/security"
)

// сервис для создания, редактирования и работы с задачами (рекламными
// компаниями)
type Service interface {
	// создать драфт задачи
	CreateTaskDraft(context.Context, *CreateTaskDraftPayload) (res string, err error)
	// обновить информацию о задаче. Не меняет статус задачи, можно вызывать
	// сколько угодно раз.
	// Нельзя вызвать для задачи, которая уже выполняется, для этого надо сначала
	// остановить выполнение.
	UpdateTask(context.Context, *UpdateTaskPayload) (res *Task, err error)
	// загрузить файл с пользователями, прокси
	UploadVideo(context.Context, *UploadVideoPayload) (res *UploadVideoResult, err error)
	// загрузить файл с пользователями, прокси
	UploadFiles(context.Context, *UploadFilesPayload) (res *UploadFilesResult, err error)
	// присвоить ботам прокси
	AssignProxies(context.Context, *AssignProxiesPayload) (res *AssignProxiesResult, err error)
	// удалить задачу и все связанные с ней сущности. Использовать только для тестов
	ForceDelete(context.Context, *ForceDeletePayload) (err error)
	// начать выполнение задачи
	StartTask(context.Context, *StartTaskPayload) (res *StartTaskResult, err error)
	// начать выполнение задачи для конкретных ботов
	PartialStartTask(context.Context, *PartialStartTaskPayload) (res *PartialStartTaskResult, err error)
	// остановить выполнение задачи
	StopTask(context.Context, *StopTaskPayload) (res *StopTaskResult, err error)
	// получить задачу по id
	GetTask(context.Context, *GetTaskPayload) (res *Task, err error)
	// получить статус выполнения задачи по id
	GetProgress(context.Context, *GetProgressPayload) (res *TaskProgress, err error)
	// получить все задачи для текущего пользователя
	ListTasks(context.Context, *ListTasksPayload) (res []*Task, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "tasks_service"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [12]string{"create task draft", "update task", "upload video", "upload files", "assign proxies", "force delete", "start task", "partial start task", "stop task", "get task", "get progress", "list tasks"}

// AssignProxiesPayload is the payload type of the tasks_service service assign
// proxies method.
type AssignProxiesPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

// AssignProxiesResult is the result type of the tasks_service service assign
// proxies method.
type AssignProxiesResult struct {
	// количество аккаунтов с проксями, которые будут использованы для текущей
	// задачи
	BotsNumber int `json:"bots_number"`
	Status     TaskStatus
	// id задачи
	TaskID string `json:"task_id"`
}

type BotAccountRecord struct {
	Record []string
	// номер строки в исходном файле
	LineNumber int `json:"line_number"`
}

type BotsProgress struct {
	// имя пользователя бота
	UserName string `json:"user_name"`
	// количество выложенных постов
	PostsCount int32 `json:"posts_count"`
	// текущий статус бота, будут ли выкладываться посты
	Status int32
	// количество аккаунтов, которых упомянули в постах
	DescriptionTargetsNotified int32 `json:"description_targets_notified"`
	// количество аккаунтов, которых упомянули в постах на фото
	PhotoTargetsNotified int32 `json:"photo_targets_notified"`
	// номер бота в загруженном файле
	FileOrder int32 `json:"file_order"`
}

// CreateTaskDraftPayload is the payload type of the tasks_service service
// create task draft method.
type CreateTaskDraftPayload struct {
	// JWT used for authentication
	Token string
	// название задачи
	Title string
	// шаблон для подписи под постом
	TextTemplate string `json:"text_template"`
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
	Type       TaskType
}

// ForceDeletePayload is the payload type of the tasks_service service force
// delete method.
type ForceDeletePayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

// GetProgressPayload is the payload type of the tasks_service service get
// progress method.
type GetProgressPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
	// размер страницы для пагинации
	PageSize uint32 `json:"page_size"`
	// номер страницы для пагинации
	Page uint32
	Sort string
	// сортировать по убыванию или нет
	SortDescending bool
}

// GetTaskPayload is the payload type of the tasks_service service get task
// method.
type GetTaskPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

// ListTasksPayload is the payload type of the tasks_service service list tasks
// method.
type ListTasksPayload struct {
	// JWT used for authentication
	Token string
}

// PartialStartTaskPayload is the payload type of the tasks_service service
// partial start task method.
type PartialStartTaskPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
	// список имен ботов, которых нужно запустить
	Usernames []string
}

// PartialStartTaskResult is the result type of the tasks_service service
// partial start task method.
type PartialStartTaskResult struct {
	// id задачи
	TaskID string `json:"task_id"`
	// список успешных имен ботов
	Succeeded []string
	// имена живых аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
}

type ProxyRecord struct {
	Record []string
	// номер строки в исходном файле
	LineNumber int `json:"line_number"`
}

// StartTaskPayload is the payload type of the tasks_service service start task
// method.
type StartTaskPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

// StartTaskResult is the result type of the tasks_service service start task
// method.
type StartTaskResult struct {
	Status TaskStatus
	// id задачи
	TaskID string `json:"task_id"`
	// имена живых аккаунтов, на которых ведем трафик
	LandingAccounts []string `json:"landing_accounts"`
}

// StopTaskPayload is the payload type of the tasks_service service stop task
// method.
type StopTaskPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

// StopTaskResult is the result type of the tasks_service service stop task
// method.
type StopTaskResult struct {
	Status TaskStatus
	// id задачи
	TaskID string `json:"task_id"`
}

type TargetUserRecord struct {
	Record []string
	// номер строки в исходном файле
	LineNumber int `json:"line_number"`
}

// Task is the result type of the tasks_service service update task method.
type Task struct {
	ID   string
	Type TaskType
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
	Status  TaskStatus
	// название задачи
	Title string
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

type TaskFileNames struct {
	// название файла, из которого брали ботов
	BotsFilename string `json:"bots_filename"`
	// название файла, из которого брали резидентские прокси
	ResidentialProxiesFilename string `json:"residential_proxies_filename"`
	// название файла, из которого брали дешёвые прокси
	CheapProxiesFilename string `json:"cheap_proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename string `json:"targets_filename"`
}

// TaskProgress is the result type of the tasks_service service get progress
// method.
type TaskProgress struct {
	// результат работы по каждому боту
	BotsProgresses []*BotsProgress `json:"bots_progresses"`
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
	Done bool
}

// 1 - задача только создана, нужно загрузить список ботов, прокси и получателей
// 2- в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
// 3- задача готова к запуску
// 4- задача запущена
// 5 - задача остановлена
// 6 - задача завершена
type TaskStatus int

// 1 - загружаем изображения
// 2- загружаем видео в рилсы
type TaskType int

// UpdateTaskPayload is the payload type of the tasks_service service update
// task method.
type UpdateTaskPayload struct {
	// JWT used for authentication
	Token string
	// id задачи, которую хотим обновить
	TaskID string `json:"task_id"`
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
	Title *string
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

type UploadError struct {
	// 1 - список ботов
	// 2 - список прокси
	// 3 - список получателей рекламы
	Type int
	Line int
	// номер порта
	Input  string
	Reason string
}

// UploadFilesPayload is the payload type of the tasks_service service upload
// files method.
type UploadFilesPayload struct {
	// JWT used for authentication
	Token string
	// id задачи, в которую загружаем пользователей/прокси
	TaskID    string `json:"task_id"`
	Filenames *TaskFileNames
	// список ботов
	Bots []*BotAccountRecord
	// список проксей для использования
	ResidentialProxies []*ProxyRecord
	// список дешёвых проксей для загрузки фото
	CheapProxies []*ProxyRecord
	// список аккаунтов, которым показать надо рекламу
	Targets []*TargetUserRecord
}

// UploadFilesResult is the result type of the tasks_service service upload
// files method.
type UploadFilesResult struct {
	// ошибки, которые возникли при загрузке файлов
	UploadErrors []*UploadError `json:"upload_errors"`
	Status       TaskStatus
}

// UploadVideoPayload is the payload type of the tasks_service service upload
// video method.
type UploadVideoPayload struct {
	// JWT used for authentication
	Token string
	// id задачи, в которую загружаем пользователей/прокси
	TaskID string `json:"task_id"`
	// не нужно присылать руками, подставится автоматом
	Filename *string
	Video    []byte
}

// UploadVideoResult is the result type of the tasks_service service upload
// video method.
type UploadVideoResult struct {
	Status TaskStatus
}

// Invalid request
type BadRequest string

// internal error
type InternalError string

// Not found
type TaskNotFound string

// Credentials are invalid
type Unauthorized string

// Error returns an error description.
func (e BadRequest) Error() string {
	return "Invalid request"
}

// ErrorName returns "bad request".
func (e BadRequest) ErrorName() string {
	return "bad request"
}

// Error returns an error description.
func (e InternalError) Error() string {
	return "internal error"
}

// ErrorName returns "internal error".
func (e InternalError) ErrorName() string {
	return "internal error"
}

// Error returns an error description.
func (e TaskNotFound) Error() string {
	return "Not found"
}

// ErrorName returns "task not found".
func (e TaskNotFound) ErrorName() string {
	return "task not found"
}

// Error returns an error description.
func (e Unauthorized) Error() string {
	return "Credentials are invalid"
}

// ErrorName returns "unauthorized".
func (e Unauthorized) ErrorName() string {
	return "unauthorized"
}
