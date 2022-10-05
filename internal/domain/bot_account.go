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
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/headers"
)

type BotAccount dbmodel.BotAccount

func (b BotAccount) ProxyURL() *url.URL {
	var buf bytes.Buffer

	if b.ResProxy == nil {
		return nil
	}

	buf.WriteString(b.ResProxy.Host)
	buf.WriteByte(':')
	buf.WriteString(strconv.FormatInt(int64(b.ResProxy.Port), 10))

	return &url.URL{
		Scheme: "http",
		User:   url.UserPassword(b.ResProxy.Login, b.ResProxy.Pass),
		Host:   buf.String(),
	}
}

func ParseBotAccounts(bots []*tasksservice.BotAccountRecord) (BotAccounts, []*tasksservice.UploadError) {
	botAccounts := make([]BotAccount, len(bots))
	var errs []*tasksservice.UploadError
	var err error

	for i, botAccountRecord := range bots {
		err = botAccounts[i].parse(botAccountRecord.Record)
		if err != nil {
			errs = append(errs, &tasksservice.UploadError{
				Type:   tasksservice.BotAccountUploadErrorType,
				Line:   botAccountRecord.LineNumber,
				Input:  strings.Join(botAccountRecord.Record, "|"),
				Reason: err.Error(),
			})
		}
	}

	return botAccounts, errs
}

// parse заполняет информацию об аккаунте
func (b *BotAccount) parse(fields []string) error {
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
	headersMap := parseHeaders(input)

	auth, ok := headersMap["authorization"]
	if !ok {
		return fmt.Errorf("'Authorization' header is missing in %+v", headersMap)
	}
	authData, err := parseSessionToken(strings.TrimPrefix(auth, authHeaderPrefix))
	if err != nil {
		return err
	}

	authData.SessionID = strings.ReplaceAll(authData.SessionID, "%3A", ":")

	b.Headers = headers.Base{
		Mid:             headersMap["x-mid"],
		DsUserID:        headersMap["ig-u-ds-user-id"],
		Rur:             headersMap["ig-u-rur"],
		Authorization:   auth,
		WWWClaim:        headersMap["x-ig-www-claim"],
		AuthData:        authData,
		BlocksVersionID: buildBloksVersioningID(b.DeviceData),
	}

	return nil
}

// parseHeaders из строки вида key1=val1;key2=val2 делает мапу {key1: val1, ke2:val2}
func parseHeaders(input string) map[string]string {
	headerPairs := strings.Split(input, ";")
	m := make(map[string]string)
	for _, pair := range headerPairs {
		keyAndValue := strings.SplitN(pair, "=", 2)
		if len(keyAndValue) != 2 {
			continue
		}

		m[strings.ToLower(keyAndValue[0])] = keyAndValue[1]
	}

	return m
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

	if authData.SessionID == "" {
		return headers.AuthorizationData{}, fmt.Errorf("got empty session id from %s", authToken)
	}

	return authData, nil
}

// ConstructHeaders создает заголовки, которые используются для запросов
func (b BotAccount) ConstructHeaders(contentLen int64) map[string][]string {
	h := headers.Default()

	h["User-Agent"] = []string{b.UserAgent}
	h["X-Pigeon-Session-Id"] = []string{"UFS-" + uuid.NewString() + "-1"}
	h["X-IG-Device-ID"] = []string{b.Session.UUID.String()}
	h["X-IG-Family-Device-ID"] = []string{b.Session.FamilyDeviceID.String()}
	h["X-IG-Android-ID"] = []string{b.Session.DeviceID}
	h["X-MID"] = []string{b.Headers.Mid}
	h["X-Bloks-Version-Id"] = []string{b.Headers.BlocksVersionID}
	h["Content-Length"] = []string{strconv.FormatInt(contentLen, 10)}
	h["Authorization"] = []string{b.Headers.Authorization}

	return h
}

func (b BotAccount) PrepareDevice() string {
	return fmt.Sprintf(
		`{"manufacturer": "%s","model": "%s","android_version": %d, "android_release": "%s"}`,
		b.DeviceData.Manufacturer,
		b.DeviceData.Model,
		b.DeviceData.AndroidVersion,
		b.DeviceData.AndroidRelease,
	)
}

func buildBloksVersioningID(deviceSettings headers.DeviceSettings) string {
	deviceBytes := []byte(fmt.Sprintf("{\"app_version\": \"%s\", \"android_version\": %d, \"android_release\": \"%s\", \"dpi\": \"%s\", \"resolution\": \"%s\", \"manufacturer\": \"%s\", \"device\": \"%s\", \"model\": \"%s\", \"cpu\": \"%s\", \"version_code\": \"%s\"}",
		deviceSettings.AppVersion, deviceSettings.AndroidVersion, deviceSettings.AndroidRelease, deviceSettings.Dpi, deviceSettings.Resolution, deviceSettings.Manufacturer, deviceSettings.Device, deviceSettings.Model, deviceSettings.Cpu, deviceSettings.VersionCode))
	sum256 := sha256.Sum256(deviceBytes)
	return hex.EncodeToString(sum256[:])
}
