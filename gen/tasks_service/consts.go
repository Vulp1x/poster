package tasksservice

const (
	// 1 - список ботов
	BotAccountUploadErrorType int = iota + 1
	// 2 - список прокси
	ProxiesUploadErrorType
	// 3 - список получателей рекламы
	TargetsUploadErrorType
)
