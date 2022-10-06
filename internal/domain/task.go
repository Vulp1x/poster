package domain

import (
	"context"
	"encoding/base64"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
)

type Task dbmodel.Task

func (t Task) ToProto() *tasksservice.Task {
	return &tasksservice.Task{
		ID:              t.ID.String(),
		TextTemplate:    t.TextTemplate,
		Image:           base64.StdEncoding.EncodeToString(t.Image),
		Status:          int(t.Status),
		Title:           t.Title,
		BotsNum:         -1,
		ProxiesNum:      -1,
		TargetsNum:      -1,
		BotsFilename:    t.BotsFilename,
		ProxiesFilename: t.ProxiesFilename,
		TargetsFilename: t.TargetsFilename,
	}
}

type BotWithTargets struct {
	BotAccount
	Targets []dbmodel.TargetUser
}

// PostingsPipe общий интерфейс для создания постов
type PostingsPipe interface {
	Process(ctx context.Context, account *BotWithTargets) error
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
