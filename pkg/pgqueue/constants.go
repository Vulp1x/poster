package pgqueue

import "errors"

var (
	// ErrUnexpectedTaskKind возвращается, если был передан неизвестный тип задачи.
	ErrUnexpectedTaskKind = errors.New("unexpected task kind")
	// ErrMustCancelTask возвращается пользовательским обработчиком,
	// если задачу необходимо закрыть, несмотря на оставшиеся попытки.
	ErrMustCancelTask = errors.New("must cancel")
	// ErrMustIgnore возвращается пользовательским обработчиком,
	// если задача не может быть выполнена в данный момент и это ожидаемое поведение.
	// По этой ошибке собирается отдельная метрика. Попытка задачи считается использованной.
	ErrMustIgnore = errors.New("must ignore")
)
