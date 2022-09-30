package domain

import (
	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
)

type BotAccounts []BotAccount

func (b BotAccounts) ToSaveParams(taskID uuid.UUID) []dbmodel.SaveBotAccountsParams {
	dbBots := make([]dbmodel.SaveBotAccountsParams, 0, len(b))
	uniqueMap := make(map[string]bool, len(b))

	for _, botAccount := range b {
		_, ok := uniqueMap[botAccount.Username]
		if ok {
			// бот с таким username уже есть, пропускаем его
			continue
		}

		dbBots = append(dbBots, dbmodel.SaveBotAccountsParams{
			TaskID:     taskID,
			Username:   botAccount.Username,
			Password:   botAccount.Password,
			UserAgent:  botAccount.UserAgent,
			DeviceData: botAccount.DeviceData,
			Session:    botAccount.Session,
			Headers:    botAccount.Headers,
			Status:     dbmodel.CreatedBotStatus,
		})

		uniqueMap[botAccount.Username] = true

	}

	return dbBots
}

// Ids возвращает список айдишников аккаунтов
func Ids[T IDer](models []T) []uuid.UUID {
	ids := make([]uuid.UUID, len(models))

	for i, account := range models {
		ids[i] = account.GetID()
	}

	return ids
}
