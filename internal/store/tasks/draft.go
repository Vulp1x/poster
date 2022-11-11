package tasks

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
)

// CreateDraftTask создает драфт задачи со
func (s *Store) CreateDraftTask(
	ctx context.Context,
	userID uuid.UUID,
	title, textTemplate string,
	accounts []string,
	images [][]byte,
	taskType tasksservice.TaskType,
	opts ...DraftOption,
) (uuid.UUID, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	dbTaskType := dbmodel.PhotoTaskType

	switch taskType {
	case tasksservice.TaskType(dbmodel.ReelsTaskType):
		dbTaskType = dbmodel.ReelsTaskType
	case tasksservice.TaskType(dbmodel.PhotoTaskType):
		dbTaskType = dbmodel.PhotoTaskType

	default:
		return [16]byte{}, fmt.Errorf(
			"%w: expected %d or %d, got: %d",
			ErrUnexpectedTaskType, dbmodel.PhotoTaskType, dbmodel.ReelsTaskType, taskType,
		)

	}

	mandatoryParams := dbmodel.CreateDraftTaskParams{
		ManagerID:       userID,
		TextTemplate:    textTemplate,
		LandingAccounts: accounts,
		Images:          images,
		Title:           title,
		Type:            dbTaskType,
	}

	for _, opt := range opts {
		opt(&mandatoryParams)
	}

	taskID, err := q.CreateDraftTask(ctx, mandatoryParams)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create task draft: %w", err)
	}

	return taskID, nil
}

// DraftOption позволяет добавить опциональные поля для создания драфтовой задачи
type DraftOption func(params *dbmodel.CreateDraftTaskParams)

// CreateDraftWithAccountNames добавляет имена ботов
func CreateDraftWithAccountNames(botNames []string) DraftOption {
	return func(params *dbmodel.CreateDraftTaskParams) {
		if len(botNames) != 0 {
			params.AccountNames = botNames
		}
	}
}

// CreateDraftWithAccountLastNames добавляет фамилии ботов
func CreateDraftWithAccountLastNames(botLastNames []string) DraftOption {
	return func(params *dbmodel.CreateDraftTaskParams) {
		if len(botLastNames) != 0 {
			params.AccountLastNames = botLastNames
		}
	}
}

// CreateDraftWithAccountURLs добавляет ссылки в профилях ботов
func CreateDraftWithAccountURLs(botURLs []string) DraftOption {
	return func(params *dbmodel.CreateDraftTaskParams) {
		if len(botURLs) != 0 {
			params.AccountUrls = botURLs
		}
	}
}

// CreateDraftWithAccountProfileImages добавляет фотографии профилей для ботов
func CreateDraftWithAccountProfileImages(profileImages [][]byte) DraftOption {
	return func(params *dbmodel.CreateDraftTaskParams) {
		if len(profileImages) != 0 {
			params.AccountProfileImages = profileImages
		}
	}
}
