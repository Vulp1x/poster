package dbmodel

type taskStatus int16

const (
	// DraftTaskStatus задача только создана, нужно загрузить список ботов, прокси и получателей
	DraftTaskStatus taskStatus = 1
	// DataUploadedTaskStatus в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	DataUploadedTaskStatus taskStatus = 2
	// ReadyTaskStatus задача готова к запуску
	ReadyTaskStatus   taskStatus = 3
	StartedTaskStatus taskStatus = 4
	// StoppedTaskStatus - задача остановлена
	StoppedTaskStatus taskStatus = 5
	// DoneTaskStatus задача выполнена
	DoneTaskStatus taskStatus = 6
	// UpdatingPostContentsTaskStatus обновляем описания под публикациями
	UpdatingPostContentsTaskStatus taskStatus = 7
	// AllDoneTaskStatus все этапы задачи выполнены, конечный статус
	AllDoneTaskStatus taskStatus = 8
)

type taskType int16

const (
	// PhotoTaskType задача только создана, нужно загрузить список ботов, прокси и получателей
	PhotoTaskType taskType = 1
	// ReelsTaskType в задачу загрузили необходимые списки, нужно присвоить прокси для ботов
	ReelsTaskType taskType = 2
)
