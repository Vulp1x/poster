package domain

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/headers"
)

type BotAccount struct {
	dbmodel.BotAccount
}

func (b BotAccount) ProxyURL() *url.URL {
	var buf bytes.Buffer

	if b.ResProxy == nil {
		return nil
	}

	buf.WriteString(b.ResProxy.Host)
	buf.WriteByte(':')
	buf.WriteString(b.ResProxy.Port)

	return &url.URL{
		Scheme: "http",
		User:   url.UserPassword(b.ResProxy.Login, b.ResProxy.Pass),
		Host:   buf.String(),
	}
}

// Parse заполняет информацию об аккаунте
func (b *BotAccount) Parse(fields []string) error {
	err := b.assignLoginAndPassword(fields[0])
	if err != nil {
		return err
	}

	err = b.assignUserAgent(fields[1])
	if err != nil {
		return err
	}

	err = b.assignSessionData(fields[2])
	if err != nil {
		return err
	}

	err = b.assignHeaders(fields[3])
	if err != nil {
		return err
	}

	return nil
}

func (b *BotAccount) assignLoginAndPassword(input string) error {
	loginWithPass := strings.Split(input, ":")
	if len(loginWithPass) != 2 {
		return fmt.Errorf("failed to parse {login}:{pass} from '%s'", input)
	}

	b.Username = loginWithPass[0]
	b.Password = loginWithPass[1]

	return nil
}

func (b *BotAccount) assignUserAgent(input string) error {
	var err error

	b.UserAgent = input
	b.DeviceData, err = headers.NewDeviceSettings(input)
	if err != nil {
		return fmt.Errorf("failed to parse device settings from user agent: %v", err)
	}

	return nil
}

func (b *BotAccount) assignSessionData(input string) error {
	ids := strings.Split(input, ";")
	if len(ids) != 4 {
		return fmt.Errorf("expected 4 ids in '%s', got %d", input, len(ids))
	}

	sessionUUID, err := uuid.Parse(ids[1])
	if err != nil {
		return fmt.Errorf("failed to parse uuid from '%s': %v", ids[1], err)
	}

	phoneID, err := uuid.Parse(ids[2])
	if err != nil {
		return fmt.Errorf("failed to parse phone uuid from '%s': %v", ids[2], err)
	}

	advertisingID, err := uuid.Parse(ids[3])
	if err != nil {
		return fmt.Errorf("failed to parse advertising uuid from '%s': %v", ids[3], err)
	}

	b.Session = headers.Session{
		DeviceID:       ids[0],
		UUID:           sessionUUID,
		PhoneID:        phoneID,
		AdvertisingID:  advertisingID,
		FamilyDeviceID: uuid.New(),
	}

	return nil
}

var headersRegexp = regexp.MustCompile(`(?m)X-MID=(.*);IG-U-DS-USER-ID=(.+);IG-U-RUR=(.+);Authorization=(.+);X-IG-WWW-Claim=(.+)`)

const authHeaderPrefix = "Bearer IGT:2:"

func (b *BotAccount) assignHeaders(input string) error {
	matches := headersRegexp.FindStringSubmatch(input)
	if len(matches) != 6 {
		return fmt.Errorf("from headers '%s' got %d matches, expected %d", input, len(matches), 6)
	}

	authData, err := parseSessionToken(strings.TrimPrefix(matches[4], authHeaderPrefix))
	if err != nil {
		return err
	}

	authData.SessionID = strings.ReplaceAll(authData.SessionID, "%3A", ":")

	b.Headers = headers.Base{
		Mid:             matches[1],
		DsUserID:        matches[2],
		Rur:             matches[3],
		Authorization:   matches[4],
		WWWClaim:        matches[5],
		AuthData:        authData,
		BlocksVersionID: buildBloksVersioningID(b.DeviceData),
	}

	return nil
}

func parseSessionToken(authToken string) (headers.AuthorizationData, error) {
	tokenBytes, err := base64.StdEncoding.DecodeString(authToken)
	if err != nil {
		return headers.AuthorizationData{}, fmt.Errorf("failed to decode base64 token from '%s': %v", authToken, err)
	}

	var authData headers.AuthorizationData
	err = json.Unmarshal(tokenBytes, &authData)
	if err != nil {
		return headers.AuthorizationData{}, fmt.Errorf("failed to unmarshal auth data from '%s': %v", string(tokenBytes), err)
	}

	return authData, nil
}

func (b BotAccount) Header(contentLen int64) map[string][]string {
	h := headers.Default()

	h["User-Agent"] = []string{b.UserAgent}
	h["X-Pigeon-Session-Id"] = []string{"UFS-" + uuid.NewString() + "-1"}
	h["X-IG-Device-ID"] = []string{b.Session.UUID.String()}
	h["X-IG-Family-Device-ID"] = []string{b.Session.FamilyDeviceID.String()}
	h["X-IG-Android-ID"] = []string{b.Session.DeviceID}
	h["X-MID"] = []string{b.Headers.Mid}
	h["X-Bloks-Version-Id"] = []string{b.Headers.BlocksVersionID}
	h["Content-Length"] = []string{strconv.FormatInt(contentLen, 10)}

	return h
}

func buildBloksVersioningID(deviceSettings headers.DeviceSettings) string {
	deviceBytes := []byte(fmt.Sprintf("{\"app_version\": \"%s\", \"android_version\": %d, \"android_release\": \"%s\", \"dpi\": \"%s\", \"resolution\": \"%s\", \"manufacturer\": \"%s\", \"device\": \"%s\", \"model\": \"%s\", \"cpu\": \"%s\", \"version_code\": \"%s\"}",
		deviceSettings.AppVersion, deviceSettings.AndroidVersion, deviceSettings.AndroidRelease, deviceSettings.Dpi, deviceSettings.Resolution, deviceSettings.Manufacturer, deviceSettings.Device, deviceSettings.Model, deviceSettings.Cpu, deviceSettings.VersionCode))
	sum256 := sha256.Sum256(deviceBytes)
	return hex.EncodeToString(sum256[:])
}
