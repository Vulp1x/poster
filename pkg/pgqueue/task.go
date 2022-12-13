package pgqueue

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

// Task структура задачи для обработки.
type Task struct {
	// Тип задачи.
	Kind int16
	// Payload задачи.
	Payload []byte
	// Необязательный внешний ключ.
	// В связке с типом задачи он задает ключ идемпотентности.
	// В рамках [KindOptions] задача по нему добавляется в очередь лишь раз
	// и может быть закрыта с помощью [CancelTaskByKey].
	ExternalKey string

	id              int64
	attemptsLeft    int16
	attemptsElapsed int16
}

// TaskHandler - интерфейс пользовательского обработчика задачи.
type TaskHandler interface {
	HandleTask(ctx context.Context, task Task) error
}

// handleTaskFunc обертка над методом интерфейса TaskHandler.
// Нужна для анонимной реализации TaskHandler.
type handleTaskFunc func(ctx context.Context, task Task) error

// HandleTask реализует TaskHandler.HandleTask.
func (h handleTaskFunc) HandleTask(ctx context.Context, task Task) error {
	return h(ctx, task)
}

func (t Task) wrapExecution(ctx context.Context, handler TaskHandler, timeout time.Duration) (time.Duration, error) {
	handlerCtx, handlerCtxCancel := context.WithTimeout(ctx, timeout)
	defer handlerCtxCancel()

	startTime := time.Now()
	err := t.executeWithRecover(handlerCtx, handler)
	elapsedTime := time.Since(startTime)

	return elapsedTime, err
}

func (t Task) executeWithRecover(ctx context.Context, handler TaskHandler) (err error) {
	defer func() {
		if p := recover(); p != nil {
			trace := string(debug.Stack())
			err = fmt.Errorf("panic in task handler: %v, stack trace: %v", p, trace)
		}
	}()
	return handler.HandleTask(ctx, t)
}
