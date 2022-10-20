package domain

import (
	"context"
	"encoding/base64"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
)

type Task dbmodel.Task

func (t Task) ToProto() *tasksservice.Task {
	images := make([]string, len(t.Images))
	for i, image := range t.Images {
		images[i] = base64.StdEncoding.EncodeToString(image)
	}

	botsProfileImages := make([]string, len(t.AccountProfileImages))
	for i, image := range t.AccountProfileImages {
		botsProfileImages[i] = base64.StdEncoding.EncodeToString(image)
	}

	return &tasksservice.Task{
		ID:                         t.ID.String(),
		TextTemplate:               t.TextTemplate,
		PostImages:                 images,
		LandingAccounts:            t.LandingAccounts,
		BotNames:                   t.AccountNames,
		BotLastNames:               t.AccountLastNames,
		BotImages:                  botsProfileImages,
		BotUrls:                    t.AccountUrls,
		Status:                     tasksservice.TaskStatus(t.Status),
		Title:                      t.Title,
		BotsNum:                    -1,
		ResidentialProxiesNum:      -1,
		CheapProxiesNum:            -1,
		TargetsNum:                 -1,
		BotsFilename:               t.BotsFilename,
		ResidentialProxiesFilename: t.ResProxiesFilename,
		CheapProxiesFilename:       t.CheapProxiesFilename,
		TargetsFilename:            t.TargetsFilename,
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
	images := make([]string, len(t.Images))
	for i, image := range t.Images {
		images[i] = base64.StdEncoding.EncodeToString(image)
	}

	botsProfileImages := make([]string, len(t.AccountProfileImages))
	for i, image := range t.AccountProfileImages {
		botsProfileImages[i] = base64.StdEncoding.EncodeToString(image)
	}

	task := Task(t.Task).ToProto()

	task.BotsNum = int(t.BotsCount)
	task.ResidentialProxiesNum = int(t.ResidentialProxiesCount)
	task.CheapProxiesNum = int(t.CheapProxiesCount)
	task.TargetsNum = int(t.TargetsCount)

	return task
}

type TaskProgress []dbmodel.GetTaskProgressRow

func (p TaskProgress) ToProto() []*tasksservice.BotsProgress {
	protos := make([]*tasksservice.BotsProgress, len(p))

	for i, row := range p {
		protos[i] = &tasksservice.BotsProgress{
			UserName:   row.Username,
			PostsCount: int(row.PostsCount),
			Status:     int(row.Status),
		}
	}

	return protos
}
