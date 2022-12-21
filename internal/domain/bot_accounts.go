package domain

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	api "github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
)

func BotsFromDBModels(dbBots []dbmodel.BotAccount) BotAccounts {
	domainBots := make([]BotAccount, len(dbBots))
	for i, bot := range dbBots {
		domainBots[i] = BotAccount(bot)
	}

	return domainBots
}

type BotAccounts []BotAccount

func (b BotAccounts) ToSaveParams(taskID uuid.UUID) []dbmodel.SaveBotAccountsParams {
	dbBots := make([]dbmodel.SaveBotAccountsParams, 0, len(b))
	uniqueMap := make(map[string]bool, len(b))

	var fileOrder = 1

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
			FileOrder:  int32(fileOrder),
			InstID:     botAccount.InstID,
		})

		uniqueMap[botAccount.Username] = true
		fileOrder++

	}

	return dbBots
}

func (b BotAccounts) ToGRPCProto(ctx context.Context) []*api.Bot {
	protoBots := make([]*api.Bot, 0, len(b))

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

		protoBots = append(protoBots, &api.Bot{
			Pk:        userID,
			Username:  botAccount.Username,
			Password:  botAccount.Password,
			SessionId: botAccount.Headers.AuthData.SessionID,
			Proxy: &api.Proxy{
				Host:  proxy.Host,
				Port:  proxy.Port,
				Login: proxy.Login,
				Pass:  proxy.Pass,
			},
			Settings: &api.BotSettings{
				UserAgent: botAccount.UserAgent,
				Bearer:    botAccount.Headers.Authorization,
				Headers: &api.BotSettings_Headers{
					Rur:            botAccount.Headers.Rur,
					Shbid:          "",
					Shbts:          "",
					Xmid:           botAccount.Headers.Mid,
					AndroidId:      botAccount.Session.DeviceID,
					DeviceId:       botAccount.Session.UUID.String(),
					PhoneId:        botAccount.Session.PhoneID.String(),
					AdvertisingId:  botAccount.Session.AdvertisingID.String(),
					FamilyDeviceId: botAccount.Session.FamilyDeviceID.String(),
				},
				Device: &api.BotSettings_DeviceSettings{
					AppVersion:     botAccount.DeviceData.AppVersion,
					AndroidVersion: int32(botAccount.DeviceData.AndroidVersion),
					AndroidRelease: botAccount.DeviceData.AndroidRelease,
					Dpi:            botAccount.DeviceData.Dpi,
					Resolution:     botAccount.DeviceData.Resolution,
					Manufacturer:   botAccount.DeviceData.Manufacturer,
					Device:         botAccount.DeviceData.Device,
					Model:          botAccount.DeviceData.Model,
					Cpu:            botAccount.DeviceData.Cpu,
					VersionCode:    botAccount.DeviceData.VersionCode,
				},
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
