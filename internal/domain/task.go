package domain

import (
	"context"
	"encoding/base64"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
)

type TaskWithCtx struct {
	dbmodel.Task
	Ctx context.Context
}

type TaskPerBot struct {
	BotAccount
	Targets []dbmodel.TargetUser
}

// PostingsPipe общий интерфейс для создания постов
type PostingsPipe interface {
	Process(ctx context.Context, account *TaskPerBot) error
}

//
// func (t *TaskWithCancel) Cancel() {
// 	t.cancel()
// }

type TaskWithCounters struct {
	dbmodel.Task
	dbmodel.SelectCountsForTaskRow
}

func (t TaskWithCounters) ToProto() *tasksservice.Task {
	return &tasksservice.Task{
		ID:              t.ID.String(),
		TextTemplate:    t.TextTemplate,
		Image:           base64.StdEncoding.EncodeToString(t.Image),
		Status:          int(t.Status),
		Title:           t.Title,
		BotsNum:         int(t.BotsCount),
		ProxiesNum:      int(t.ProxiesCount),
		TargetsNum:      int(t.TargetsCount),
		BotsFilename:    t.BotsFilename,
		ProxiesFilename: t.ProxiesFilename,
		TargetsFilename: t.TargetsFilename,
	}
}
