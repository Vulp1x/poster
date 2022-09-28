package domain

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
)

type Proxies []Proxy

func (b Proxies) ToSaveParams(taskID uuid.UUID) []dbmodel.SaveProxiesParams {
	dbProxies := make([]dbmodel.SaveProxiesParams, 0, len(b))
	uniqueMap := make(map[string]bool, len(b))

	for _, proxy := range b {
		uniqueKey := proxy.Host + strconv.FormatInt(int64(proxy.Port), 10)
		_, ok := uniqueMap[uniqueKey]
		if ok {
			// прокси с таким username уже есть, пропускаем его
			continue
		}
		dbProxies = append(dbProxies, dbmodel.SaveProxiesParams{
			TaskID: taskID,
			Host:   proxy.Host,
			Port:   proxy.Port,
			Login:  proxy.Login,
			Pass:   proxy.Pass,
			Type:   2,
		})

		uniqueMap[uniqueKey] = true
	}

	return dbProxies
}
