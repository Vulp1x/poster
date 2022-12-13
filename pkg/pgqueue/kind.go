package pgqueue

import (
	"fmt"
	"math"
	"time"

	"github.com/inst-api/poster/pkg/pgqueue/internal/workerpool"
	"github.com/inst-api/poster/pkg/pgqueue/pkg/delayer"
)

type kindData struct {
	handler TaskHandler
	opts    KindOptions
	wp      *workerpool.Pool
}

// ValueProvider - провайдер значений.
type ValueProvider[T any] func() T

// NewConstProvider возвращает провайдер константы.
func NewConstProvider[T any](value T) ValueProvider[T] {
	return func() T {
		return value
	}
}

// KindOptions содержит настройки типа задач.
type KindOptions struct {
	// Название типа задач.
	// Отображается в метриках и логах, поэтому важно использовать уникальное и понятное значение.
	// По умолчанию: fmt.Sprintf("pgqueue-%v", kind).
	Name string
	// Число воркеров (горутин), выполняющих задания.
	// Примерная оценка оптимального значения (закон Литтла):
	// workerCount = alpha * tasksPerSecond * averageTaskDuration / replicaCount
	// Значение alpha нужно подобрать таким, чтобы иметь достаточный запас воркеров
	// на случай крактовременного увеличения трафика.
	// По умолчанию: NewConstProvider(int16(1)).
	WorkerCount ValueProvider[int16]
	// Максимальное число попыток выполнить задачу.
	// По умолчанию: math.MaxInt16.
	MaxAttempts int16
	// Максимальное число последних сообщений ошибок выполнения задачи, которые сохраняются.
	// По умолчанию: 1.
	MaxTaskErrorMessages int16
	// Таймаут на попытку выполнения задачи.
	// По умолчанию: паника.
	AttemptTimeout time.Duration
	// Рассчитывает по номеру попытки время задержки перед следующей попыткой выполнить задачу.
	// Предлагается использовать пакет [delayer].
	// По умолчанию: паника.
	Delayer delayer.Delayer
	// Время, в течение которого задачи в терминальном статусе хранятся в БД.
	// Чтобы задачи в терминальном статусе удалялись сразу, передавайте time.Duration(0).
	// По умолчанию: NewConstProvider(3 * 24 * time.Hour).
	TerminalTasksTTL ValueProvider[time.Duration]
	// Внутренние настройки очереди. Будьте аккуратны при их изменении.
	Loop LoopOptions
}

// LoopOptions внутренние настройки очереди.
type LoopOptions struct {
	// Период, с которым запускается удаление задач в терминальном статусе.
	// По умолчанию: NewConstProvider(time.Hour).
	JanitorPeriod ValueProvider[time.Duration]
	// Период, с которым запускаются поиск и отправка на выполнение задач.
	// Не рекомендуется использовать период меньше секунды, так как нагрузка на БД может сильно вырасти.
	// По умолчанию: NewConstProvider(1500 * time.Millisecond).
	FetcherPeriod ValueProvider[time.Duration]
}

func checkAndSetDefaultKindOptions(kind int16, opts *KindOptions) {
	if opts.AttemptTimeout == 0 {
		panic("KindOptions.AttemptTimeout must not be zero")
	}

	if opts.Delayer == nil {
		panic("KindOptions.Delayer must not be nil")
	}

	if len(opts.Name) == 0 {
		opts.Name = fmt.Sprintf("pgqueue-%v", kind)
	}

	if opts.WorkerCount == nil {
		opts.WorkerCount = NewConstProvider(int16(1))
	}

	if opts.MaxAttempts == 0 {
		opts.MaxAttempts = math.MaxInt16
	}

	if opts.MaxTaskErrorMessages == 0 {
		opts.MaxTaskErrorMessages = 1
	}

	if opts.TerminalTasksTTL == nil {
		opts.TerminalTasksTTL = NewConstProvider(3 * 24 * time.Hour)
	}

	if opts.Loop.JanitorPeriod == nil {
		opts.Loop.JanitorPeriod = NewConstProvider(time.Hour)
	}

	if opts.Loop.FetcherPeriod == nil {
		opts.Loop.FetcherPeriod = NewConstProvider(1500 * time.Millisecond)
	}
}
