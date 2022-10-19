package domain

import (
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
)

type TasksWithCounters []TaskWithCounters

func (t TasksWithCounters) ToProto() []*tasksservice.Task {
	protoTasks := make([]*tasksservice.Task, len(t))
	for i, task := range t {
		task.AccountProfileImages = nil
		task.Images = nil
		protoTasks[i] = task.ToProto()
	}

	return protoTasks
}
