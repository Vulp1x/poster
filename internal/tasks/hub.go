package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/transport"
	"github.com/inst-api/poster/pkg/logger"
)

type Hub struct {
	tasksChan   chan domain.TaskWithCtx
	taskCancels map[uuid.UUID]func()
	taskMu      *sync.Mutex
	pushTimeout time.Duration
	dbtxf       dbmodel.DBTXFunc
}

func NewHub(timeout time.Duration) *Hub {
	return &Hub{
		tasksChan:   make(chan domain.TaskWithCtx, 10),
		taskCancels: make(map[uuid.UUID]func()),
		pushTimeout: timeout,
	}
}

const workersPerTask = 10

func (h *Hub) StartTask(ctx context.Context, taskID uuid.UUID) error {
	q := dbmodel.New(h.dbtxf(ctx))

	task, err := q.StartTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	taskCtx, taskCancel := context.WithCancel(ctx)

	taskWithCtx := domain.TaskWithCtx{
		Task: task,
		Ctx:  taskCtx,
	}

	botsChan := make(chan *domain.TaskPerBot, 20)

	var workers []*worker
	for i := 0; i < workersPerTask; i++ {
		workers = append(workers, &worker{
			tasksQueue:     botsChan,
			dbtxf:          h.dbtxf,
			cli:            transport.InitHTTPClient(),
			postPipe:       nil,
			processorIndex: 0,
		})
	}

	select {
	case h.tasksChan <- taskWithCtx:
		h.taskCancels[task.ID] = taskCancel
		return nil

	case <-time.After(h.pushTimeout):
		logger.Debugf(ctx, "waited for %s, failed to push task to queue")
		break
	}

	taskCancel()

	return fmt.Errorf("failed to push task to queue")
}

func (h Hub) StopTask(taskID uuid.UUID) error {
	cancel, ok := h.taskCancels[taskID]
	if !ok {
		return fmt.Errorf("failed to find task '%s' in tasks: %#v", taskID, h.taskCancels)
	}

	cancel()

	h.taskMu.Lock()
	defer h.taskMu.Unlock()

	delete(h.taskCancels, taskID)

	return nil
}
