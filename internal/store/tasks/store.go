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

type Store struct {
	tasksChan   chan domain.TaskWithCtx
	taskCancels map[uuid.UUID]func()
	taskMu      *sync.Mutex
	pushTimeout time.Duration
	dbtxf       dbmodel.DBTXFunc
}

func NewStore(timeout time.Duration, dbtxFunc dbmodel.DBTXFunc) *Store {
	return &Store{
		tasksChan:   make(chan domain.TaskWithCtx, 10),
		taskCancels: make(map[uuid.UUID]func()),
		pushTimeout: timeout,
		dbtxf:       dbtxFunc,
	}
}

const workersPerTask = 10

func (s *Store) CreateDraftTask(ctx context.Context, userID uuid.UUID, title, textTemplate string, image []byte) (uuid.UUID, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	taskID, err := q.CreateDraftTask(ctx, dbmodel.CreateDraftTaskParams{
		ManagerID:    userID,
		TextTemplate: textTemplate,
		Image:        image,
		Title:        title,
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create task draft: %w", err)
	}

	return taskID, nil
}

func (s *Store) StartTask(ctx context.Context, taskID uuid.UUID) error {
	q := dbmodel.New(s.dbtxf(ctx))

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
			dbtxf:          s.dbtxf,
			cli:            transport.InitHTTPClient(),
			processorIndex: int64(i),
		})
	}

	select {
	case s.tasksChan <- taskWithCtx:
		s.taskCancels[task.ID] = taskCancel
		return nil

	case <-time.After(s.pushTimeout):
		logger.Debugf(ctx, "waited for %s, failed to push task to queue")
		break
	}

	taskCancel()

	return fmt.Errorf("failed to push task to queue")
}

func (s *Store) StopTask(ctx context.Context, taskID uuid.UUID) error {
	logger.Infof(ctx, "stopping task '%s'", taskID)
	cancel, ok := s.taskCancels[taskID]
	if !ok {
		return fmt.Errorf("failed to find task '%s' in tasks: %#v", taskID, s.taskCancels)
	}

	cancel()

	s.taskMu.Lock()
	defer s.taskMu.Unlock()

	delete(s.taskCancels, taskID)

	return nil
}
