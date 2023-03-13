package domain

import (
	"fmt"
	"strconv"
	"strings"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
)

type TargetUser dbmodel.TargetUser

func ParseTargetUsers(targetRecords []*tasksservice.TargetUserRecord, uploadErrs []*tasksservice.UploadError) []TargetUser {
	domainProxies := make([]TargetUser, 0, len(targetRecords))
	var err error
	for _, proxyRecord := range targetRecords {
		proxy := TargetUser{}
		err = proxy.parse(proxyRecord.Record)
		if err != nil {
			uploadErrs = append(uploadErrs, &tasksservice.UploadError{
				Type:   tasksservice.TargetsUploadErrorType,
				Line:   proxyRecord.LineNumber,
				Input:  strings.Join(proxyRecord.Record, "|"),
				Reason: err.Error(),
			})

			continue
		}

		domainProxies = append(domainProxies, proxy)
	}

	return domainProxies
}

func (p *TargetUser) parse(targetRecord []string) error {
	if targetRecord[0] == "" {
		return fmt.Errorf("missing username")
	}

	var userID int64
	var err error
	if len(targetRecord) == 2 {
		userID, err = strconv.ParseInt(targetRecord[1], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse user_id: %v", err)
		}

	}

	p.Username = targetRecord[0]
	p.UserID = userID

	return nil
}
