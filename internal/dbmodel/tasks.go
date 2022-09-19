package dbmodel

type taskStatus int16

const (
	DraftTaskStatus   taskStatus = 1
	ReadyTaskStatus   taskStatus = 2
	StartedTaskStatus taskStatus = 3
	DoneTaskStatus    taskStatus = 3
)
