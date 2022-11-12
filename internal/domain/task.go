package domain

import (
	"context"
	"encoding/base64"
	"strings"

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

	var videFilenamePointer *string
	if t.VideoFilename != nil {
		parts := strings.SplitN(*t.VideoFilename, "_", 3)
		if len(parts) == 0 {
			videFilenamePointer = t.VideoFilename
		} else {
			videFilenamePointer = &parts[len(parts)-1]
		}
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
		FollowTargets:              t.FollowTargets,
		NeedPhotoTags:              t.NeedPhotoTags,
		PerPostSleepSeconds:        uint(t.PerPostSleepSeconds),
		PhotoTagsDelaySeconds:      uint(t.PhotoTagsDelaySeconds),
		PostsPerBot:                uint(t.PostsPerBot),
		TargetsPerPost:             uint(t.TargetsPerPost),
		Type:                       tasksservice.TaskType(t.Type),
		VideoFilename:              videFilenamePointer,
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

type TaskProgress struct {
	BotsProgress   []dbmodel.GetBotsProgressRow
	TargetCounters dbmodel.GetTaskTargetsCountRow
	Done           bool
}

func (p TaskProgress) ToProto() *tasksservice.TaskProgress {
	botsMap := make(map[string]*tasksservice.BotsProgress, len(p.BotsProgress))
	for _, progress := range p.BotsProgress {
		botsMap[progress.Username] = &tasksservice.BotsProgress{
			UserName:   progress.Username,
			PostsCount: int(progress.PostsCount),
			Status:     int(progress.Status),
		}
	}

	return &tasksservice.TaskProgress{
		BotsProgresses:  botsMap,
		TargetsNotified: int(p.TargetCounters.NotifiedTargets),
		TargetsFailed:   int(p.TargetCounters.FailedTargets),
		TargetsWaiting:  int(p.TargetCounters.UnusedTargets),
		Done:            p.Done,
	}
}
