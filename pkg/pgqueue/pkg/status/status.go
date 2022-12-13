// Статусы задач. Используются только для сбора метрик.
package status

const (
	// Задача добавлена в очередь и еще не взята в обработку.
	New = "new"
	// Задача неуспешно обработана и ожидает следующей попытки.
	Failed = "failed"
	// Задача взята в обработку.
	Processing = "processing"
	// Попытки выполнения задачи закончились.
	NoAttemptsLeft = "no_attempts_left"
	// Задача закрыта пользователем.
	Cancelled = "cancelled"
	// Задача успешно обработана.
	Succeeded = "succeeded"
	// Цикл обработки задачи был прерван.
	// Возможно из-за ошибок при работе с БД.
	Lost = "lost"
	// Пользовательский обработчик задач вернул ошибку ErrMustIgnore.
	Ignored = "ignored"
)

var All = []string{New, Failed, Processing, NoAttemptsLeft, Cancelled, Succeeded, Lost, Ignored}
