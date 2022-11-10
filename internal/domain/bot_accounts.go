package domain

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/pb/parser"
	"github.com/inst-api/poster/pkg/logger"
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

func (b BotAccounts) ToGRPCProto(ctx context.Context) []*parser.Bot {
	protoBots := make([]*parser.Bot, 0, len(b))

	for _, botAccount := range b {
		userID, err := strconv.ParseInt(botAccount.Headers.AuthData.DsUserID, 10, 64)
		if err != nil {
			logger.Errorf(ctx, "failed to parse user id from '%s': %v", botAccount.Headers.DsUserID, err)
			continue
		}

		proxy := botAccount.ResProxy
		if proxy == nil {
			proxy = botAccount.WorkProxy
		}

		if proxy == nil {
			continue
		}

		protoBots = append(protoBots, &parser.Bot{
			Username:  botAccount.Username,
			UserId:    userID,
			SessionId: botAccount.Headers.AuthData.SessionID,
			Proxy: &parser.Proxy{
				Host:  proxy.Host,
				Port:  proxy.Port,
				Login: proxy.Login,
				Pass:  proxy.Pass,
			},
		})
	}

	return protoBots
}

// Ids возвращает список айдишников аккаунтов
func Ids[T IDer](models []T) []uuid.UUID {
	ids := make([]uuid.UUID, len(models))

	for i, account := range models {
		ids[i] = account.GetID()
	}

	return ids
}
