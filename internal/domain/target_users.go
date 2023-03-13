package domain

import (
	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
)

type TargetUsers []TargetUser

func (t TargetUsers) ToSaveParams(taskID uuid.UUID) []dbmodel.SaveTargetUsersParams {
	dbTargetUsers := make([]dbmodel.SaveTargetUsersParams, 0, len(t))
	uniqueMap := make(map[string]struct{}, len(t))

	for _, target := range t {
		_, ok := uniqueMap[target.Username]
		if ok {
			// target user с таким username уже есть, пропускаем его
			continue
		}

		dbTargetUsers = append(dbTargetUsers, dbmodel.SaveTargetUsersParams{
			TaskID:   taskID,
			Username: target.Username,
			UserID:   target.UserID,
			// Status:   dbmodel.UnusedTargetStatus,
		})

		uniqueMap[target.Username] = struct{}{}
	}

	return dbTargetUsers
}
