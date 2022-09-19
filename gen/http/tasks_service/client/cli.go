// Code generated by goa v3.8.5, DO NOT EDIT.
//
// tasks_service HTTP client CLI support package
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package client

import (
	"encoding/json"
	"fmt"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	goa "goa.design/goa/v3/pkg"
)

// BuildCreateTaskPayload builds the payload for the tasks_service create task
// endpoint from CLI flags.
func BuildCreateTaskPayload(tasksServiceCreateTaskBody string, tasksServiceCreateTaskToken string) (*tasksservice.CreateTaskPayload, error) {
	var err error
	var body CreateTaskRequestBody
	{
		err = json.Unmarshal([]byte(tasksServiceCreateTaskBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"description\": \"Ipsam expedita libero eum et.\",\n      \"tittle\": \"Dolor culpa temporibus sit.\"\n   }'")
		}
	}
	var token string
	{
		token = tasksServiceCreateTaskToken
	}
	v := &tasksservice.CreateTaskPayload{
		Tittle:      body.Tittle,
		Description: body.Description,
	}
	v.Token = token

	return v, nil
}

// BuildUploadFilePayload builds the payload for the tasks_service upload file
// endpoint from CLI flags.
func BuildUploadFilePayload(tasksServiceUploadFileBody string, tasksServiceUploadFileTaskID string, tasksServiceUploadFileToken string) (*tasksservice.UploadFilePayload, error) {
	var err error
	var body UploadFileRequestBody
	{
		err = json.Unmarshal([]byte(tasksServiceUploadFileBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"bytes\": \"RXZlbmlldCBtb2xlc3RpYWUgc2ludCByZXJ1bSBldCBvZGl0Lg==\",\n      \"file_type\": 1\n   }'")
		}
		if !(body.FileType == 1 || body.FileType == 2 || body.FileType == 3) {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("body.file_type", body.FileType, []interface{}{1, 2, 3}))
		}
		if err != nil {
			return nil, err
		}
	}
	var taskID string
	{
		taskID = tasksServiceUploadFileTaskID
	}
	var token string
	{
		token = tasksServiceUploadFileToken
	}
	v := &tasksservice.UploadFilePayload{
		FileType: body.FileType,
		Bytes:    body.Bytes,
	}
	v.TaskID = taskID
	v.Token = token

	return v, nil
}

// BuildStartTaskPayload builds the payload for the tasks_service start task
// endpoint from CLI flags.
func BuildStartTaskPayload(tasksServiceStartTaskTaskID string, tasksServiceStartTaskToken string) (*tasksservice.StartTaskPayload, error) {
	var taskID string
	{
		taskID = tasksServiceStartTaskTaskID
	}
	var token string
	{
		token = tasksServiceStartTaskToken
	}
	v := &tasksservice.StartTaskPayload{}
	v.TaskID = taskID
	v.Token = token

	return v, nil
}

// BuildGetTaskPayload builds the payload for the tasks_service get task
// endpoint from CLI flags.
func BuildGetTaskPayload(tasksServiceGetTaskTaskID string, tasksServiceGetTaskToken string) (*tasksservice.GetTaskPayload, error) {
	var taskID string
	{
		taskID = tasksServiceGetTaskTaskID
	}
	var token string
	{
		token = tasksServiceGetTaskToken
	}
	v := &tasksservice.GetTaskPayload{}
	v.TaskID = taskID
	v.Token = token

	return v, nil
}

// BuildListTasksPayload builds the payload for the tasks_service list tasks
// endpoint from CLI flags.
func BuildListTasksPayload(tasksServiceListTasksToken string) (*tasksservice.ListTasksPayload, error) {
	var token string
	{
		token = tasksServiceListTasksToken
	}
	v := &tasksservice.ListTasksPayload{}
	v.Token = token

	return v, nil
}
