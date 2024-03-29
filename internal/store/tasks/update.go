package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/store"
	"github.com/jackc/pgx/v4"
)

func (s *Store) UpdateTask(ctx context.Context, taskID uuid.UUID, opts ...UpdateOption) (domain.Task, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return domain.Task{}, store.ErrTransactionFail
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Task{}, ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if task.Status == dbmodel.AllDoneTaskStatus {
		return domain.Task{}, fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.AllDoneTaskStatus, task.Status)
	}

	for _, opt := range opts {
		opt(&task)
	}

	task, err = q.UpdateTask(ctx, dbmodel.UpdateTaskParams{
		TextTemplate:          task.TextTemplate,
		Title:                 task.Title,
		Images:                task.Images,
		AccountNames:          task.AccountNames,
		AccountLastNames:      task.AccountLastNames,
		AccountUrls:           task.AccountUrls,
		AccountProfileImages:  task.AccountProfileImages,
		LandingAccounts:       task.LandingAccounts,
		FollowTargets:         task.FollowTargets,
		NeedPhotoTags:         task.NeedPhotoTags,
		PerPostSleepSeconds:   task.PerPostSleepSeconds,
		PhotoTagsDelaySeconds: task.PhotoTagsDelaySeconds,
		PostsPerBot:           task.PostsPerBot,
		TargetsPerPost:        task.TargetsPerPost,
		PhotoTagsPostsPerBot:  task.PhotoTagsPostsPerBot,
		PhotoTargetsPerPost:   task.PhotoTargetsPerPost,
		FixedTag:              task.FixedTag,
		FixedPhotoTag:         task.FixedPhotoTag,
		ID:                    taskID,
	})
	if err != nil {
		return domain.Task{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Task{}, err
	}

	return domain.Task(task), nil
}

// UpdateOption позволяет добавить опциональные поля для создания драфтовой задачи
type UpdateOption func(params *dbmodel.Task)

func WithBotNamesUpdateOption(names []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(names) != 0 {
			task.AccountNames = names
		}
	}
}

func WithBotLasNamesUpdateOption(lastNames []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(lastNames) != 0 {
			task.AccountLastNames = lastNames
		}
	}
}

func WithBotURLsUpdateOption(urls []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(urls) != 0 {
			task.AccountUrls = urls
		}
	}
}

func WithBotImagesUpdateOption(images [][]byte) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(images) != 0 {
			task.AccountProfileImages = images
		}
	}
}

func WithImagesUpdateOption(images [][]byte) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(images) != 0 {
			task.Images = images
		}
	}
}

func WithLandingAccountsUpdateOption(landingAccounts []string) UpdateOption {
	return func(task *dbmodel.Task) {
		if len(landingAccounts) != 0 {
			task.LandingAccounts = landingAccounts
		}
	}
}

func WithTextTemplateUpdateOption(template *string) UpdateOption {
	return func(task *dbmodel.Task) {
		if template != nil {
			task.TextTemplate = *template
		}
	}
}

func WithTitleUpdateOption(title *string) UpdateOption {
	return func(task *dbmodel.Task) {
		if title != nil {
			task.Title = *title
		}
	}
}

// WithFollowTargets добавляет ссылки в профилях ботов
func WithFollowTargets(followTargets *bool) UpdateOption {
	return func(params *dbmodel.Task) {
		if followTargets != nil {
			params.FollowTargets = *followTargets
		}
	}
}

// WithNeedPhotoTags добавляет ссылки в профилях ботов
func WithNeedPhotoTags(needPhotoTags *bool) UpdateOption {
	return func(params *dbmodel.Task) {
		if needPhotoTags != nil {
			params.NeedPhotoTags = *needPhotoTags
		}
	}
}

// WithPhotoTagsDelaySeconds добавляет ссылки в профилях ботов
func WithPhotoTagsDelaySeconds(photoTagsDelaySeconds *uint) UpdateOption {
	return func(params *dbmodel.Task) {
		if photoTagsDelaySeconds != nil {
			params.PhotoTagsDelaySeconds = int32(*photoTagsDelaySeconds)
		}
	}
}

// WithPostsPerBot добавляет ссылки в профилях ботов
func WithPostsPerBot(postsPerBot *uint) UpdateOption {
	return func(params *dbmodel.Task) {
		if postsPerBot != nil {
			params.PostsPerBot = int(*postsPerBot)
		}
	}
}

// WithTargetsPerPost добавляет ссылки в профилях ботов
func WithTargetsPerPost(targetsPerPost *uint) UpdateOption {
	return func(params *dbmodel.Task) {
		if targetsPerPost != nil {
			params.TargetsPerPost = int(*targetsPerPost)
		}
	}
}

// WithPhotoTargetsPerPost добавляет ссылки в профилях ботов
func WithPhotoTargetsPerPost(targetsPerPost *uint) UpdateOption {
	return func(params *dbmodel.Task) {
		if targetsPerPost != nil {
			params.PhotoTargetsPerPost = int(*targetsPerPost)
		}
	}
}

// WithTargetsPerPost добавляет ссылки в профилях ботов
func WithPhotoTagsPostsPerBot(photoTagsPostsPerBot *uint) UpdateOption {
	return func(params *dbmodel.Task) {
		if photoTagsPostsPerBot != nil {
			params.PhotoTagsPostsPerBot = int(*photoTagsPostsPerBot)
		}
	}
}

// WithPerPostSleepSeconds добавляет ссылки в профилях ботов
func WithPerPostSleepSeconds(perPostSleepSeconds *uint) UpdateOption {
	return func(params *dbmodel.Task) {
		if perPostSleepSeconds != nil {
			params.PerPostSleepSeconds = int32(*perPostSleepSeconds)
		}
	}
}

func WithFixedTagUpdateOption(username *string) UpdateOption {
	return func(task *dbmodel.Task) {
		if username != nil && *username != "" {
			task.FixedTag = username
		} else {
			task.FixedTag = nil
		}
	}
}

func WithFixedPhotoTagUpdateOption(userID *int64) UpdateOption {
	return func(task *dbmodel.Task) {
		if userID != nil && *userID != 0 {
			task.FixedPhotoTag = userID
		} else {
			task.FixedPhotoTag = nil
		}
	}
}
