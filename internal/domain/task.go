package domain

import (
	"context"

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
