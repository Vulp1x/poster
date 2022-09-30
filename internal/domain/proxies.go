package domain

import (
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
)

type Proxies []Proxy

func (p Proxies) ToSaveParams(taskID uuid.UUID) []dbmodel.SaveProxiesParams {
	dbProxies := make([]dbmodel.SaveProxiesParams, 0, len(p))
	uniqueMap := make(map[string]bool, len(p))

	for _, proxy := range p {
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

func (p Proxies) ToStrings() []string {
	ret := make([]string, len(p))

	for i, proxy := range p {
		bytes, _ := json.Marshal(proxy)
		ret[i] = string(bytes)
	}

	return ret
}
