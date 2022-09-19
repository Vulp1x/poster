package headers

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

type DeviceSettings struct {
	AppVersion     string `json:"app_version"`
	AndroidVersion int    `json:"android_version"`
	AndroidRelease string `json:"android_release"`
	Dpi            string `json:"dpi"`
	Resolution     string `json:"resolution"`
	Manufacturer   string `json:"manufacturer"`
	Device         string `json:"device"`
	Model          string `json:"model"`
	Cpu            string `json:"cpu"`
	VersionCode    string `json:"version_code"`
}

type Session struct {
	DeviceID       string // for example: "android-0d735e1f4db26782"
	UUID           uuid.UUID
	PhoneID        uuid.UUID
	AdvertisingID  uuid.UUID
	FamilyDeviceID uuid.UUID
}

var userAgentRegexp = regexp.MustCompile(`(?m)Instagram (.+) Android \((\d+)/(.+); (\d+dpi); (\d+x\d+); (.+); (.+); (.+); (.+); (.+); (\d+)\)`)

// NewDeviceSettings возвращает информацию об устройстве основываясь на User-Agent
func NewDeviceSettings(userAgent string) (DeviceSettings, error) {
	matches := userAgentRegexp.FindStringSubmatch(userAgent)
	if len(matches) != 12 {
		return DeviceSettings{},
			fmt.Errorf("from user-agent '%s' got %d matches, expected %d", userAgent, len(matches), 12)
	}

	androidVersion, err := strconv.ParseInt(matches[2], 10, 32)
	if err != nil {
		return DeviceSettings{}, fmt.Errorf("failed to parse android version from '%s': %v", matches[2], err)
	}

	return DeviceSettings{
		AppVersion:     matches[1],
		AndroidVersion: int(androidVersion),
		AndroidRelease: matches[3],
		Dpi:            matches[4],
		Resolution:     matches[5],
		Manufacturer:   matches[6],
		Device:         matches[7],
		Model:          matches[8],
		Cpu:            matches[9],
		VersionCode:    matches[11],
	}, nil
}
