package domain

import (
	"bytes"
	"strconv"

	"github.com/inst-api/poster/internal/dbmodel"
)

type Targets []dbmodel.TargetUser

func (t Targets) ToProto(format int) []string {
	strings := make([]string, len(t))

	for i, target := range t {
		strings[i] = formatTarget(target, format)
	}

	return strings
}

func formatTarget(target dbmodel.TargetUser, format int) string {
	switch format {
	case 1:
		return strconv.FormatInt(target.UserID, 10)
	case 2:
		return target.Username
	case 3:
		b := bytes.Buffer{}

		b.WriteString(target.Username)
		b.WriteByte(',')
		b.WriteString(strconv.FormatInt(target.UserID, 10))
		return b.String()
	}

	return ""

}
