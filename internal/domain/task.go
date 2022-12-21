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
		Type:                       tasksservice.TaskType(t.Type),
		TextTemplate:               t.TextTemplate,
		LandingAccounts:            t.LandingAccounts,
		BotNames:                   t.AccountNames,
		BotLastNames:               t.AccountLastNames,
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
		VideoFilename:              videFilenamePointer,
		FollowTargets:              t.FollowTargets,
		NeedPhotoTags:              t.NeedPhotoTags,
		PerPostSleepSeconds:        uint(t.PerPostSleepSeconds),
		PhotoTagsDelaySeconds:      uint(t.PhotoTagsDelaySeconds),
		PostsPerBot:                uint(t.PostsPerBot),
		PhotoTagsPostsPerBot:       uint(t.PhotoTagsPostsPerBot),
		TargetsPerPost:             uint(t.TargetsPerPost),
		PhotoTargetsPerPost:        uint(t.PhotoTargetsPerPost),
		PostImages:                 images,
		BotImages:                  botsProfileImages,
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
	BotsProgress   []BotProgress
	TargetCounters dbmodel.GetTaskTargetsCountRow
	Done           bool
}

type BotProgress struct {
	Username               string
	Status                 int32
	PostDescriptionTargets int32
	PhotoTaggedTargets     int32
	PostsCount             int32
	FileOrder              int32
}

func (p TaskProgress) ToProto() *tasksservice.TaskProgress {
	bots := make([]*tasksservice.BotsProgress, len(p.BotsProgress))
	for i, progress := range p.BotsProgress {
		bots[i] = &tasksservice.BotsProgress{
			UserName:                   progress.Username,
			PostsCount:                 progress.PostsCount,
			Status:                     progress.Status,
			DescriptionTargetsNotified: progress.PostDescriptionTargets,
			PhotoTargetsNotified:       progress.PhotoTaggedTargets,
			FileOrder:                  progress.FileOrder,
		}
	}

	return &tasksservice.TaskProgress{
		BotsProgresses:       bots,
		TargetsNotified:      int(p.TargetCounters.DescriptionNotifiedTargets),
		PhotoTargetsNotified: int(p.TargetCounters.PhotoNotifiedTargets),
		TargetsFailed:        int(p.TargetCounters.FailedTargets),
		TargetsWaiting:       int(p.TargetCounters.UnusedTargets),
		Done:                 p.Done,
	}
}
