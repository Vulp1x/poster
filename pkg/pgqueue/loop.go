package pgqueue

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/inst-api/poster/internal/mw"
	"github.com/inst-api/poster/pkg/ctxutil"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue/internal/db"
	"github.com/inst-api/poster/pkg/pgqueue/internal/workerpool"
	"github.com/inst-api/poster/pkg/pgqueue/pkg/status"
)

// loop расширяет Queue логикой по обработке задач.
type loop struct {
	*Queue
	jobs
}

type jobs struct {
	janitor chan func()
	fetcher chan func()
}

func newLoop(queue *Queue) loop {
	return loop{
		Queue: queue,
		jobs: jobs{
			janitor: make(chan func()),
			fetcher: make(chan func()),
		},
	}
}

// Watch запускает слушателя каналов start/stop.
func (l loop) Watch(ctx context.Context) <-chan struct{} {
	defer logger.Warnf(ctx, "started watcher")

	done := make(chan struct{})
	go func() {
		defer close(done)
		defer l.shutdown()
		defer logger.Warnf(ctx, "shutting down...")

		l.watch(ctx)
	}()

	return done
}

// watch слушает каналы start/stop и запускает/останавливает обработчики задач.
func (l loop) watch(parentCtx context.Context) {
	var wg sync.WaitGroup
	defer wg.Wait()

	stopped := true
	var ctx context.Context
	var cancel context.CancelFunc

	for {
		select {
		case <-l.start:
			if stopped {
				ctx, cancel = context.WithCancel(parentCtx)
				defer cancel()
				l.run(ctx, &wg)

				stopped = false
				logger.Warnf(ctx, "started loops")
			}
		case <-l.stop:
			if !stopped {
				cancel()
				wg.Wait()

				stopped = true
				logger.Warnf(ctx, "stopped loops")
			}
		case <-parentCtx.Done():
			return
		}
	}
}

func (l loop) run(ctx context.Context, wg *sync.WaitGroup) {
	for kind, kd := range l.kds {
		func(ctx context.Context, kind int16, kd kindData) {
			ctx = logger.WithKV(ctx, "pgqueue_name", kd.opts.Name)
			runWithWG(func() { l.runJanitor(ctx, kind, kd.opts) }, wg)
			runWithWG(func() { l.runFetcher(ctx, kind, kd) }, wg)
		}(ctx, kind, kd)
	}

	runWithWG(func() { runJobWatcher(ctx, l.jobs.janitor) }, wg)
	runWithWG(func() { runJobWatcher(ctx, l.jobs.fetcher) }, wg)
}

func (l loop) runJanitor(ctx context.Context, kind int16, opts KindOptions) {
	janitor := func() {
		if err := l.storage(ctx).CleanupTasks(ctx, db.CleanupTasksParams{
			Kind:    kind,
			Timeout: db.SQLInterval(opts.TerminalTasksTTL()),
			Period:  db.SQLInterval(opts.Loop.JanitorPeriod()),
		}); err == nil {
			logger.Warnf(ctx, "CleanupTasks succeeded for [%v]", opts.Name)
		} else if !errors.Is(err, db.ErrRegistryIsBusy) {
			logger.Errorf(ctx, "CleanupTasks failed: %v", err)
		}
	}

	runWithPeriod(ctx, l.jobs.janitor, janitor, opts.Loop.JanitorPeriod)
}

func (l loop) runFetcher(ctx context.Context, kind int16, kd kindData) {
	fetcher := func() {
		resizeWP(ctx, kd.wp, int32(kd.opts.WorkerCount()))
		if kd.wp.RestingCount() > 0 {
			tasks, err := l.storage(ctx).GetOpenTasks(ctx, db.GetOpenTasksParams{
				Kind:          kind,
				Limit:         kd.wp.RestingCount(),
				UntilDeadline: db.SQLInterval(2 * (kd.opts.AttemptTimeout + storageBackoffMaxElapsedTime)),
			})
			if err != nil {
				logger.Errorf(ctx, "GetOpenTasks failed: %v", err)
				return
			}

			for _, task := range tasks {
				kd.wp.Push(l.wrapTask(ctx, kd, task))
			}
			// mc.CollectTasksInStatus(kd.opts.Name, status.Processing, int32(len(tasks)))
		}
	}

	runWithPeriod(ctx, l.jobs.fetcher, fetcher, kd.opts.Loop.FetcherPeriod)
}

func (l loop) wrapTask(parentCtx context.Context, kd kindData, task Task) func() {
	return func() {
		ctx := logger.WithKV(ctxutil.Detach(parentCtx), "task_id", task.id)
		ctx = mw.GenerateRequestID(ctx)

		elapsedTime, executionErr := task.wrapExecution(ctx, kd.handler, kd.opts.AttemptTimeout)
		collectExecutionResults(kd.opts.Name, elapsedTime, errToStatus(executionErr))

		switch {
		case executionErr == nil:
			if kd.opts.TerminalTasksTTL() == 0 {
				if err := l.storage(ctx).DeleteTask(ctx, task.id); err != nil {
					logger.Errorf(ctx, "DeleteTask failed: %v", err)
					// mc.CollectTasksInStatus(kd.opts.Name, status.Lost, 1)
				}
			} else {
				if err := l.storage(ctx).CompleteTask(ctx, task.id); err != nil {
					logger.Errorf(ctx, "CompleteTask failed: %v", err)
					// // mc.CollectTasksInStatus(kd.opts.Name, status.Lost, 1)
				}
			}
		case errors.Is(executionErr, ErrMustCancelTask):
			logger.Warnf(ctx, "cancelling, spent [%v], error: %v", elapsedTime, executionErr)
			if kd.opts.TerminalTasksTTL() == 0 {
				if err := l.storage(ctx).DeleteTask(ctx, task.id); err != nil {
					logger.Errorf(ctx, "DeleteTask failed: %v", err)
					// mc.CollectTasksInStatus(kd.opts.Name, status.Lost, 1)
				}
			} else {
				params := db.CancelTaskParams{ID: task.id, Reason: executionErr.Error()}
				if err := l.storage(ctx).CancelTask(ctx, params); err != nil {
					logger.Errorf(ctx, "CancelTask failed: %v", err)
					// mc.CollectTasksInStatus(kd.opts.Name, status.Lost, 1)
				}
			}
		default:
			delay := kd.opts.Delayer(task.attemptsElapsed - 1)
			if task.attemptsLeft == 0 {
				// mc.CollectTasksInStatus(kd.opts.Name, status.NoAttemptsLeft, 1)
				logger.Errorf(ctx, "no attempts left, spent [%v], error: %v", elapsedTime, executionErr)
			} else {
				logger.Warnf(ctx, "attempts left [%v], spent [%v], delay [%v], error: %v", task.attemptsLeft, elapsedTime, delay, executionErr)
			}

			if err := l.storage(ctx).RefuseTask(ctx, db.RefuseTaskParams{
				ID:            task.id,
				Reason:        executionErr.Error(),
				Delay:         db.SQLInterval(delay),
				MessagesLimit: kd.opts.MaxTaskErrorMessages,
			}); err != nil {
				logger.Errorf(ctx, "RefuseTask failed: %v", err)
				// mc.CollectTasksInStatus(kd.opts.Name, status.Lost, 1)
			}
		}
	}
}

func runWithPeriod(ctx context.Context, jobs chan<- func(), job func(), period ValueProvider[time.Duration]) {
	for {
		select {
		case <-time.After(period()):
			jobs <- job
		case <-ctx.Done():
			return
		}
	}
}

// JobWatcher читает задачи из канала jobs и исполняет их.
// Нужен, чтобы распределить нагрузку на БД, запуская запросы последовательно.
func runJobWatcher(ctx context.Context, jobs <-chan func()) {
	for {
		select {
		case job := <-jobs:
			job()
		case <-ctx.Done():
			return
		}
	}
}

func runWithWG(fn func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
}

func resizeWP(ctx context.Context, wp *workerpool.Pool, size int32) {
	if s := wp.Size(); s != size {
		wp.Resize(size)
		logger.Warnf(ctx, "resized worker pool from [%v] to [%v]", s, size)
	}
}

func collectExecutionResults(name string, elapsedTime time.Duration, status string) {
	// mc.CollectTasksInStatus(name, status, 1)
	// mc.CollectTaskDuration(name, elapsedTime, status)
}

func errToStatus(err error) string {
	switch {
	case err == nil:
		return status.Succeeded
	case errors.Is(err, ErrMustCancelTask):
		return status.Cancelled
	case errors.Is(err, ErrMustIgnore):
		return status.Ignored
	default:
		return status.Failed
	}
}
