package dbmodel

// targetStatus - статус уведомления для конкретного пользователя
type targetStatus int16

const (
	UnusedTargetStatus     targetStatus = 1 // UnusedTargetStatus пользователь не упоминался в постах
	InProgressTargetStatus targetStatus = 2 // InProgressTargetStatus пользователь используется для
	FailedTargetStatus     targetStatus = 3 // FailedTargetStatus произошла ошибка при публикации
	NotifiedTargetStatus   targetStatus = 4 // NotifiedTargetStatus пользователь успешно упомянут в публикации
)
