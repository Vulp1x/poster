// Code generated by goa v3.8.5, DO NOT EDIT.
//
// HTTP request path constructors for the tasks_service service.
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package client

import (
	"fmt"
)

// CreateTaskDraftTasksServicePath returns the URL path to the tasks_service service create task draft HTTP endpoint.
func CreateTaskDraftTasksServicePath() string {
	return "/api/tasks/draft"
}

// UploadFileTasksServicePath returns the URL path to the tasks_service service upload file HTTP endpoint.
func UploadFileTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/upload", taskID)
}

// StartTaskTasksServicePath returns the URL path to the tasks_service service start task HTTP endpoint.
func StartTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/start", taskID)
}

// StopTaskTasksServicePath returns the URL path to the tasks_service service stop task HTTP endpoint.
func StopTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/stop", taskID)
}

// GetTaskTasksServicePath returns the URL path to the tasks_service service get task HTTP endpoint.
func GetTaskTasksServicePath(taskID string) string {
	return fmt.Sprintf("/api/tasks/%v/", taskID)
}

// ListTasksTasksServicePath returns the URL path to the tasks_service service list tasks HTTP endpoint.
func ListTasksTasksServicePath() string {
	return "/api/tasks/"
}
