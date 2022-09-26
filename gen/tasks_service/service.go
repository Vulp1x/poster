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
	UploadFile(context.Context, *UploadFilePayload) (res []*UploadError, err error)
	// начать выполнение задачи
	StartTask(context.Context, *StartTaskPayload) (err error)
	// остановить выполнение задачи
	StopTask(context.Context, *StopTaskPayload) (err error)
	// получить задачу по id
	GetTask(context.Context, *GetTaskPayload) (err error)
	// получить все задачи для текущего пользователя
	ListTasks(context.Context, *ListTasksPayload) (err error)
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
var MethodNames = [6]string{"create task draft", "upload file", "start task", "stop task", "get task", "list tasks"}

type BotAccount struct {
	// login
	Username string
	// login
	Password string
	// user agent header
	UserAgent string `json:"user_agent"`
	// main id, ex: android-0d735e1f4db26782
	DeviceID string `json:"device_id"`
	UUID     string
	// phone_id
	PhoneID string `json:"phone_id"`
	// adv id
	AdvertisingID  string `json:"advertising_id"`
	FamilyDeviceID string `json:"family_device_id"`
	Headers        map[string]string
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

type Proxy struct {
	// адрес прокси
	Host string
	// номер порта
	Port     int64
	Login    string
	Password string
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

type TargetUser struct {
	// instagram username
	Username string
	// instagram user id
	UserID int64 `json:"user_id"`
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

// UploadFilePayload is the payload type of the tasks_service service upload
// file method.
type UploadFilePayload struct {
	// JWT used for authentication
	Token string
	// id задачи, в которую загружаем пользователей/прокси
	TaskID string `json:"task_id"`
	// список ботов
	Bots []*BotAccount
	// список проксей для использования
	Proxies []*Proxy
	// список аккаунтов, которым показать надо рекламу
	Targets []*TargetUser
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
