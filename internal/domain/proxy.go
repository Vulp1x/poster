package domain

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
)

type Proxy dbmodel.Proxy

func ParseProxies(proxyRecords []*tasksservice.ProxyRecord, uploadErrs []*tasksservice.UploadError) Proxies {
	domainProxies := make([]Proxy, 0, len(proxyRecords))
	var err error
	for _, proxyRecord := range proxyRecords {
		proxy := Proxy{}
		err = proxy.parse(proxyRecord.Record)
		if err != nil {
			uploadErrs = append(uploadErrs, &tasksservice.UploadError{
				Type:   tasksservice.TargetsUploadErrorType,
				Line:   proxyRecord.LineNumber,
				Input:  strings.Join(proxyRecord.Record, ":"),
				Reason: err.Error(),
			})

			continue
		}

		domainProxies = append(domainProxies, proxy)
	}

	return domainProxies
}

func (p *Proxy) parse(proxyRecord []string) error {
	ip := net.ParseIP(proxyRecord[0])
	if ip == nil {
		return fmt.Errorf("failed to parse ip")
	}

	port, err := strconv.ParseInt(proxyRecord[1], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse port: %v", err)
	}

	p.Host = ip.String()
	p.Port = int32(port)
	p.Login = proxyRecord[2]
	p.Pass = proxyRecord[3]

	return nil
}
