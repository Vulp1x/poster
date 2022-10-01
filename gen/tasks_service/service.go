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
	// загрузить файл с пользователями, прокси
	UploadFiles(context.Context, *UploadFilesPayload) (res []*UploadError, err error)
	// присвоить ботам прокси
	AssignProxies(context.Context, *AssignProxiesPayload) (res int, err error)
	// удалить задачу и все связанные с ней сущности. Использовать только для тестов
	ForceDelete(context.Context, *ForceDeletePayload) (err error)
	// начать выполнение задачи
	StartTask(context.Context, *StartTaskPayload) (err error)
	// остановить выполнение задачи
	StopTask(context.Context, *StopTaskPayload) (err error)
	// получить задачу по id
	GetTask(context.Context, *GetTaskPayload) (res *Task, err error)
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
var MethodNames = [8]string{"create task draft", "upload files", "assign proxies", "force delete", "start task", "stop task", "get task", "list tasks"}

// AssignProxiesPayload is the payload type of the tasks_service service assign
// proxies method.
type AssignProxiesPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

type BotAccountRecord struct {
	Record []string
	// номер строки в исходном файле
	LineNumber int `json:"line_number"`
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
	// фотография для постов
	PostImage string `json:"post_image"`
}

// ForceDeletePayload is the payload type of the tasks_service service force
// delete method.
type ForceDeletePayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
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

// StopTaskPayload is the payload type of the tasks_service service stop task
// method.
type StopTaskPayload struct {
	// JWT used for authentication
	Token string
	// id задачи
	TaskID string `json:"task_id"`
}

type TargetUserRecord struct {
	Record []string
	// номер строки в исходном файле
	LineNumber int `json:"line_number"`
}

// Task is the result type of the tasks_service service get task method.
type Task struct {
	ID string
	// описание под постом
	TextTemplate string `json:"text_template"`
	// base64 строка картинки
	Image  string
	Status int
	// название задачи
	Title string
	// количество ботов в задаче
	BotsNum int `json:"bots_num"`
	// количество прокси в задаче
	ProxiesNum int `json:"proxies_num"`
	// количество целевых пользователей в задаче
	TargetsNum int `json:"targets_num"`
	// название файла, из которого брали ботов
	BotsFilename *string `json:"bots_filename"`
	// название файла, из которого брали прокси
	ProxiesFilename *string `json:"proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename *string `json:"targets_filename"`
}

type TaskFileNames struct {
	// название файла, из которого брали ботов
	BotsFilename string `json:"bots_filename"`
	// название файла, из которого брали прокси
	ProxiesFilename string `json:"proxies_filename"`
	// название файла, из которого брали целевых пользователей
	TargetsFilename string `json:"targets_filename"`
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
	Proxies []*ProxyRecord
	// список аккаунтов, которым показать надо рекламу
	Targets []*TargetUserRecord
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
