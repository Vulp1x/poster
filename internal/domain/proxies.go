package domain

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
)

type Proxies []Proxy

func (p Proxies) ToSaveParams(taskID uuid.UUID, isCheap bool) []dbmodel.SaveProxiesParams {
	dbProxies := make([]dbmodel.SaveProxiesParams, 0, len(p))
	// uniqueMap := make(map[string]bool, len(p))

	proxyType := dbmodel.ResidentialProxyType
	if isCheap {
		proxyType = dbmodel.CheapProxyType
	}

	for _, proxy := range p {
		// uniqueKey := proxy.Host + strconv.FormatInt(int64(proxy.Port), 10)
		// _, ok := uniqueMap[uniqueKey]
		// if ok {
		// 	прокси с таким username уже есть, пропускаем его
		// continue
		// }

		dbProxies = append(dbProxies, dbmodel.SaveProxiesParams{
			TaskID: taskID,
			Host:   proxy.Host,
			Port:   proxy.Port,
			Login:  proxy.Login,
			Pass:   proxy.Pass,
			Type:   proxyType,
		})

		// uniqueMap[uniqueKey] = true
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

func ProxiesFromDB(dbProxies []dbmodel.Proxy) Proxies {
	proxies := make([]Proxy, len(dbProxies))
	for i, proxy := range dbProxies {
		proxies[i] = Proxy(proxy)
	}

	return proxies
}
