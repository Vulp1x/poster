# pgqueue | postgres job queue

[[_TOC_]]

## Важное

Документация доступна локально:
```shell
make docs
```

### Задача

Описывается видом `Kind`, данными `Payload` и обработчиком `TaskHandler`.
При регистрации вида задач `Queue.RegisterKind` задаются важные настройки `KindOptions`.

В `Payload` рекомендуется сохранять небольшой объем данных (e.g. `id`).
Это снижает нагрузку на очередь и включает оптимизацию за счет `fillfactor`.

### Идемпотентность

Вид задачи `Kind` и опциональный ключ `ExternalKey` задают ключ идемпотентности,
по которому задача добавляется в очередь один раз (в рамках [Retention Policy](#retention-policy)).

### Retention Policy

Из очереди удаляются только задачи, находящиеся в одном из терминальных состояний:
  - `cancelled`: задача закрыта по ключу идемпотентности `Queue.CancelTaskByKey` или ошибкой `ErrMustCancelTask`,
  - `succeeded`: задача успешно выполнена.

Время `TerminalTasksTTL`, в течение которого хранятся задачи в терминальном статусе, задается в `KindOptions`.

## Использование

### 1. Миграции

Необходимо перенести миграции из [migrations](migrations) в свой проект.

### 2. Адаптеры

Для работы с различными драйверами/библиотеками в pgqueue используется абстракция [executor](pkg/executor/executor.go).

Адаптеры под [executor](pkg/executor/executor.go) для `database-go` и `database-pg` 
доступны в [pkg/adapter/go](pkg/adapter/go) и [pkg/adapter/pg](pkg/adapter/pg).

### 3. Метрики

Предлагается использовать пакет [pkg/metrics](pkg/metrics).
Метрики собираются по числу задач в каждом из статусов (см. [status](pkg/status/status.go)) и времени выполнения.

### 4. Пример

Полный пример: [example](example).

---

#### Создание очереди

##### database-go

```go
import (
	pgqueuego "github.com/inst-api/poster/pkg/pgqueuepkg/adapter/go"
	pgqueuemetrics "github.com/inst-api/poster/pkg/pgqueuepkg/metrics"
)

queueCtx, queueCancel := context.WithCancel(logger.WithName(context.Background(), "pgqueue"))
queue := pgqueuego.NewQueue(queueCtx, pgqueuego.AdaptBalancer(balancer), pgqueuemetrics.NewCollector())
```

##### database-pg

```go
import (
	pgqueuepg "github.com/inst-api/poster/pkg/pgqueuepkg/adapter/pg"
	pgqueuemetrics "github.com/inst-api/poster/pkg/pgqueuepkg/metrics"
)

queueCtx, queueCancel := context.WithCancel(logger.WithName(context.Background(), "pgqueue"))
queue := pgqueuepg.NewQueue(queueCtx, pgqueuepg.AdaptRoleCluster(roleCluster), pgqueuemetrics.NewCollector())
// или
queue := pgqueuepg.NewQueue(queueCtx, pgqueuepg.AdaptShardsCluster(shardsCluster, "shard_key"), pgqueuemetrics.NewCollector())
```

#### Регистрация типов задач

При регистрации типов задач задаются важные настройки `KindOptions`.
```go
queue.RegisterKind(taskKind, taskHandler, pgqueue.KindOptions{
	// Название типа задач.
	// Отображается в метриках и логах, поэтому важно использовать уникальное и понятное значение.
	Name: "short-descriptive-name",
	// Число воркеров на конкретный под на данный тип задач
	// Примерная оценка оптимального значения (закон Литтла):
	// workerCount = alpha * tasksPerSecond * averageTaskDuration / replicaCount
	// Значение alpha нужно подобрать таким, чтобы иметь достаточный запас воркеров
	// на случай крактовременного увеличения трафика.
	WorkerCount: pgqueue.NewConstProvider(int16(100)),
	// Число попыток выполнения задачи. Попытки можно добавить методом RetryTasks.
	MaxAttempts: 10,
	// Таймаут на попытку выполнения задачи.
	AttemptTimeout: time.Second,
	// Максимальное число последних сообщений ошибок выполнения задачи, которые сохраняются.
	MaxTaskErrorMessages: 3,
	// Рассчитывает по номеру попытки время задержки перед следующей попыткой выполнить задачу.
	Delayer: delayer.NewJitterDelayer(delayer.EqualJitter, 6*time.Second),
	// Время, в течение которого задачи в терминальном статусе хранятся в БД.
	// Чтобы задачи в терминальном статусе удалялись сразу, передавайте time.Duration(0).
	TerminalTasksTTL: pgqueue.NewConstProvider(time.Duration(0)),
	// Внутренние настройки очереди. Будьте аккуратны при их изменении.
	Loop: pgqueue.LoopOptions{
		// ...
	},
})
```

#### Запуск обработчиков задач

```go
queue.Start()
```

#### Добавление задачи в очередь

- `PushTask` вне транзакции
- `PushTaskTx` в транзакции
- `PushTasks` вне транзакции батчом
- `PushTasksTx` в транзакции батчом

```go
err := queue.PushTask(ctx, pgqueue.Task{
		Kind: taskKind,
		Payload: payload,
		ExternalKey: externalKey,
	},
	pgqueue.WithDelay(time.Hour),
)
```

#### Остановка/возобновление обработки задач

```go
config.WatchValue(ctx, config.EnablePgqueue, func(_, n realtimeconfig.Variable) {
	if n.Value().Bool() {
		queue.Start()
	} else {
		queue.Stop()
	}
})
```

#### Закрытие очереди

```go
queueCancel()
<-queue.Done()
```

## Статусная модель

![image info](docs/pgqueue.svg)
