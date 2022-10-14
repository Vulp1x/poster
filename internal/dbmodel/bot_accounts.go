package dbmodel

import (
	"encoding/json"
	"fmt"

	"github.com/inst-api/poster/internal/headers"

	"github.com/google/uuid"
)

type botStatus int16

const (
	CreatedBotStatus        botStatus = 1
	ProxieAssignedBotStatus botStatus = 2
	StartedBotStatus        botStatus = 3
	DoneBotStatus           botStatus = 4 // DoneBotStatus проставляется после того как все посты выложены
	FailBotStatus           botStatus = 5 // FailBotStatus
)

type RequestHeaders struct {
	Device         headers.DeviceSettings
	AndroidID      string // for example: "android-0d735e1f4db26782"
	DeviceID       uuid.UUID
	FamilyDeviceID uuid.UUID
	PhoneID        uuid.UUID
	AdvertisingID  uuid.UUID

	BloksVersioningID string
	UserAgent         string                    `json:"user_agent"`
	Mid               string                    `json:"mid"`
	IgURur            string                    `json:"ig_u_rur"`
	IgWwwClaim        string                    `json:"ig_www_claim"`
	AuthData          headers.AuthorizationData `json:"authorization_data"`
}

// String is an implementation for driver.Valuer.
func (ev RequestHeaders) String() string {
	data, _ := json.Marshal(ev) // nolint:errchkjson
	return string(data)
}

// Scan is an implementation for sql.Scanner.
func (ev *RequestHeaders) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed for %T = %v", value, value)
	}

	return json.Unmarshal(data, ev)
}
