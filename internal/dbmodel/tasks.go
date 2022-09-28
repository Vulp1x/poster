package dbmodel

type taskStatus int16

const (
	DraftTaskStatus        taskStatus = 1 // DraftTaskStatus задача только создана, нужно загрузить список ботов, прокси и получателей
	DataUploadedTaskStatus taskStatus = 2 // DataUploadedTaskStatus в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	ReadyTaskStatus        taskStatus = 3 // ReadyTaskStatus задача готова к запуску
	StartedTaskStatus      taskStatus = 4
	DoneTaskStatus         taskStatus = 5
)
