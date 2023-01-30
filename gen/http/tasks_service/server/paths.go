// Code generated by goa v3.8.5, DO NOT EDIT.
//
// HTTP request path constructors for the tasks_service service.
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package server

import (
	"fmt"
)

// CreateTaskDraftTasksServicePath returns the URL path to the tasks_service service create task draft HTTP endpoint.
func CreateTaskDraftTasksServicePath() string {
	return "/api/tasks/draft/"
}

// UpdateTaskTasksServicePath returns the URL path to the tasks_service service update task HTTP endpoint.
func UpdateTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/", taskID)
}

// UploadVideoTasksServicePath returns the URL path to the tasks_service service upload video HTTP endpoint.
func UploadVideoTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/upload/video/", taskID)
}

// UploadFilesTasksServicePath returns the URL path to the tasks_service service upload files HTTP endpoint.
func UploadFilesTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/upload/", taskID)
}

// AssignProxiesTasksServicePath returns the URL path to the tasks_service service assign proxies HTTP endpoint.
func AssignProxiesTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/assign/", taskID)
}

// ForceDeleteTasksServicePath returns the URL path to the tasks_service service force delete HTTP endpoint.
func ForceDeleteTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/force/", taskID)
}

// StartTaskTasksServicePath returns the URL path to the tasks_service service start task HTTP endpoint.
func StartTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/start/", taskID)
}

// PartialStartTaskTasksServicePath returns the URL path to the tasks_service service partial start task HTTP endpoint.
func PartialStartTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/start/partial/", taskID)
}

// UpdatePostContentsTasksServicePath returns the URL path to the tasks_service service update post contents HTTP endpoint.
func UpdatePostContentsTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/start/post-contents/", taskID)
}

// StopTaskTasksServicePath returns the URL path to the tasks_service service stop task HTTP endpoint.
func StopTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/stop/", taskID)
}

// GetTaskTasksServicePath returns the URL path to the tasks_service service get task HTTP endpoint.
func GetTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/", taskID)
}

// GetProgressTasksServicePath returns the URL path to the tasks_service service get progress HTTP endpoint.
func GetProgressTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/progress", taskID)
}

// GetEditingProgressTasksServicePath returns the URL path to the tasks_service service get editing progress HTTP endpoint.
func GetEditingProgressTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/reprogress", taskID)
}

// ListTasksTasksServicePath returns the URL path to the tasks_service service list tasks HTTP endpoint.
func ListTasksTasksServicePath() string {
	return "/api/tasks/"
}

// DownloadTargetsTasksServicePath returns the URL path to the tasks_service service download targets HTTP endpoint.
func DownloadTargetsTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/targets/download/", taskID)
}

// DownloadBotsTasksServicePath returns the URL path to the tasks_service service download bots HTTP endpoint.
func DownloadBotsTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/bots/download/", taskID)
}
