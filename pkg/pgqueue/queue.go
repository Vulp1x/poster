package pgqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/inst-api/poster/pkg/pgqueue/internal/db"
	"github.com/inst-api/poster/pkg/pgqueue/internal/workerpool"
	"github.com/inst-api/poster/pkg/pgqueue/pkg/executor"
)

var mc MetricsCollector

// MetricsCollector - сборщик метрик.
type MetricsCollector interface {
	// CollectTasksInStatus регистрирует [amount] задач типа [name] в статусе [status].
	CollectTasksInStatus(name string, status string, amount int32)
	// CollectTaskDuration регистрирует длительность выполнения задачи типа [name].
	CollectTaskDuration(name string, duration time.Duration, status string)
	// RegisterKind регистрирует тип задачи, создавая заранее всевозможные метрики.
	// Необходимо для решения проблемы, связанной с тем, что счетчик стартует не с нуля.
	// Подробнее: https://www.section.io/blog/beware-prometheus-counters-that-do-not-begin-at-zero.
	RegisterKind(name string)
}

// Cluster - кластер БД с методом получения подключения к мастеру.
// Используйте пакеты в adapter для обертки database-go и database-pg.
type Cluster interface {
	Master(ctx context.Context) executor.Executor
}

// Queue - очередь задач.
type Queue struct {
	executor executor.Executor
	kds      map[int16]kindData

	start chan struct{}
	stop  chan struct{}
	done  <-chan struct{}
}

// NewQueue создает очередь задач.
func NewQueue(ctx context.Context, executor executor.Executor) *Queue {
	queue := &Queue{
		executor: executor,
		kds:      make(map[int16]kindData),
		start:    make(chan struct{}),
		stop:     make(chan struct{}),
	}
	queue.done = newLoop(queue).Watch(ctx)
	return queue
}

// Start запускает обработчики задач. Вызов блокирующий.
// Повторные и асинхронные вызовы безопасны.
func (q *Queue) Start() {
	q.start <- struct{}{}
}

// Stop останавливает обработчики задач. Вызов блокирующий.
// Повторные и асинхронные вызовы безопасны.
func (q *Queue) Stop() {
	q.stop <- struct{}{}
}

// Done возвращает канал завершения очереди.
func (q *Queue) Done() <-chan struct{} {
	return q.done
}

// RegisterKind регистрирует вид задач в очереди.
func (q *Queue) RegisterKind(kind int16, handler TaskHandler, opts KindOptions) *Queue {
	if _, ok := q.kds[kind]; ok {
		panic(fmt.Sprintf("kind = [%v] already exists", kind))
	}

	checkAndSetDefaultKindOptions(kind, &opts)
	q.kds[kind] = kindData{handler: handler, opts: opts, wp: workerpool.New()}
	// mc.RegisterKind(opts.Name)

	return q
}

// RegisterKindFunc регистрирует вид задач в очереди.
func (q *Queue) RegisterKindFunc(kind int16, handler func(ctx context.Context, task Task) error, opts KindOptions) *Queue {
	return q.RegisterKind(kind, handleTaskFunc(handler), opts)
}

// PushTask публикует задачу в очередь.
func (q *Queue) PushTask(ctx context.Context, task Task, opts ...PushOption) error {
	return q.storage(ctx).PushTasks(ctx, q.kds, []Task{task}, opts...)
}

// PushTaskTx публикует задачу в очередь в переданной транзакции.
func (q *Queue) PushTaskTx(ctx context.Context, tx executor.Tx, task Task, opts ...PushOption) error {
	return q.storage(ctx, tx).PushTasks(ctx, q.kds, []Task{task}, opts...)
}

// PushTasks публикует задачи в очередь batch'ом.
func (q *Queue) PushTasks(ctx context.Context, tasks []Task, opts ...PushOption) error {
	return q.storage(ctx).PushTasks(ctx, q.kds, tasks, opts...)
}

// PushTasksTx публикует задачи в очередь batch'ом в переданной транзакции.
func (q *Queue) PushTasksTx(ctx context.Context, tx executor.Tx, tasks []Task, opts ...PushOption) error {
	return q.storage(ctx, tx).PushTasks(ctx, q.kds, tasks, opts...)
}

// CancelTaskByKey закрывает открытую задачу по ключу идемпотентности.
// Если задача не была найдена в открытом статусе, ошибка не возвращается.
func (q *Queue) CancelTaskByKey(ctx context.Context, kind int16, key string) error {
	_, ok := q.kds[kind]
	if !ok {
		return ErrUnexpectedTaskKind
	}

	// mc.CollectTasksInStatus(kd.opts.Name, status.Cancelled, 1)
	params := db.CancelTaskByKeyParams{Reason: "CancelTaskByKey", Kind: kind, ExternalKey: db.NullString(key)}
	return q.storage(ctx).CancelTaskByKey(ctx, params)
}

// RetryTasks восстанавливает не более [limit] задач типа [kind], у которых закончились попытки.
// Задачи восстанавливаются в порядке добавления в очередь.
func (q *Queue) RetryTasks(ctx context.Context, kind, attempts int16, limit int32) error {
	if _, ok := q.kds[kind]; !ok {
		return ErrUnexpectedTaskKind
	}

	params := db.RetryTasksParams{Kind: kind, AttemptsLeft: attempts, Limit: limit}
	if err := q.storage(ctx).RetryTasks(ctx, params); err != nil {
		return fmt.Errorf("RetryTasks failed: %w", err)
	}

	return nil
}

func (q *Queue) storage(ctx context.Context, db ...executor.Executor) *storage {
	if len(db) > 0 {
		return newStorage(db[0])
	}
	return newStorage(q.executor)
}

func (q *Queue) shutdown() {
	for _, kd := range q.kds {
		kd.wp.CloseAndWait()
	}
}
