package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/mw"
	"github.com/inst-api/poster/internal/store/tasks"
	"github.com/inst-api/poster/pkg/logger"
)

type progressStore interface {
	TaskProgress(ctx context.Context, taskID uuid.UUID) (domain.TaskProgress, error)
}

const (
	bufferSize     = 1024
	maxMessageSize = 512
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{ //nolint
	ReadBufferSize:  bufferSize,
	WriteBufferSize: bufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Hub struct {
	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	taskProgressStore progressStore

	mu *sync.Mutex
	// ключ - task_id
	clients             map[uuid.UUID][]*Client
	getProgressEndpoint func(ctx context.Context, p *tasksservice.GetProgressPayload) (*tasksservice.TaskProgress, error)
}

type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	taskID uuid.UUID

	ctx context.Context

	cancel context.CancelFunc
}

func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		logger.Debugf(c.ctx, "returning from ws write func")
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			logger.Warnf(c.ctx, "got error on closing ws connection: %v", err)
		}
	}()

	logger.Debugf(c.ctx, "start writing goroutine")

	for {
		select {
		case <-c.ctx.Done():
			logger.Debugf(c.ctx, "context done, returning in write")
			return
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				return
			}
			if err = c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Warnf(c.ctx, "failed to write ping message")
				return
			}
		}
	}

}

func NewWebSocketHub(driverGeoStore progressStore) *Hub {
	return &Hub{
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		taskProgressStore: driverGeoStore,
	}
}

func (h *Hub) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/tasks/{taskID}/progress", h.serveTaskProgress)
	})

	return r
}

// serveTaskProgress handles websocket requests for one driver's geo updates.
func (h *Hub) serveTaskProgress(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	logger.Debug(ctx, "upgrading ws connection")

	taskID, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		logger.Errorf(ctx, "failed to parse taskID from url: %v", err)
		http.Error(w, "Bad task_id provided", http.StatusBadRequest)

		return
	}

	progress, err := h.taskProgressStore.TaskProgress(r.Context(), taskID)
	if err != nil {
		logger.Errorf(ctx, "failed to get task progress: %v", err)
		if errors.Is(err, tasks.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
		}

		if errors.Is(err, tasks.ErrTaskInvalidStatus) {
			http.Error(w, "invalid status", http.StatusBadRequest)
		}

		mw.InternalError(w, r, err.Error())
	}

	proto := progress.ToProto()
	progressBytes, err := json.Marshal(proto)
	if err != nil {
		logger.Errorf(ctx, "failed to marshal task progress '%+v': %v", proto, err)
		mw.InternalError(w, r, "failed to marshal task progress")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf(ctx, "Failed to upgrade websocket: %v", err)
		mw.InternalError(w, r, "failed to upgrade ws connection")
		return
	}

	ctx, cancel := context.WithCancel(logger.WithFields(ctx, logger.Fields{"driver_id": taskID}))

	client := &Client{
		hub:    h,
		conn:   conn,
		taskID: taskID,
		ctx:    ctx,
		cancel: cancel,
	}

	logger.Debug(ctx, "registered new client")

	_, _ = w.Write(progressBytes)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
}

func (h Hub) PushBotProgress(ctx context.Context, taskID uuid.UUID) error {
	return nil
}
