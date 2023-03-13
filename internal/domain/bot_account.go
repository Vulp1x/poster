package domain

import (
	"bytes"
	"context"
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

var userAgents = []string{
	"Instagram 265.0.0.19.301 Android (24/7.0; 240dpi; 2069x1080; samsung; SM-G901F; kccat6; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (27/8.1.0; 560dpi; 1080x2233; TCT; ONE TOUCH 4015D; Yaris35_GSM; mt6572; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 640dpi; 1440x2780; samsung; SM-J415G; j4primelte; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.2; 320dpi; 2224x1080; infinix; infinixnote3; note3; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (24/7.0; 280dpi; 1080x2280; Huawei; 201HW; hwu9201L; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 356dpi; 1080x2006; SENSEIT; SENSEIT-L301; SENSEITL301; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1; 160dpi; 1080x2062; Samsung; SM-G3502T; cs02ve3gdtv; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (24/7.0; 540dpi; 1080x2031; HUAWEI/HONOR; DUA-L22; HWDUA-M; m; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (28/9; 544dpi; 1440x2891; Lenovo; LenovoA396; A396; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 272dpi; 1440x2711; HUAWEI/HONOR; LLD-L31; HWLLD-H; hi6; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (29/10; 420dpi; 1080x2210; samsung; sm-g930f; herolte; sams; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.1; 440dpi; 1440x2910; HMD Global/Nokia; Nokia 8.1; PNX_sprout; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1; 540dpi; 1440x2672; motorola; MotoE2; otus_ds; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (27/8.1.0; 213dpi; 2042x1080; Blackview; BV7000; BV7000; mt6735; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 640dpi; 720x1369; HUAWEI; HUAWEI TAG-L01; HWTAG-L6753; mt6735; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 320dpi; 720x1472; ZTE; Z730; ada; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 160dpi; 640x960; samsung; SM-N970U1; d1q; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 540dpi; 1080x1808; samsung; sm-g920f; zeroflte; samsungexynos7; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.2; 440dpi; 480x894; LGE; LG-SU760; cosmo_450-05; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (27/8.1.0; 272dpi; 1080x2092; TCT; ALCATELONETOUCH; d2aio; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1; 480dpi; 1080x2207; Samsung; SAMSUNG-SGH-I747Z; d2aio; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 213dpi; 640x1208; Lenovo; IdeaTabS2109A-F; S2109A; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 120dpi; 600x1162; Sony; E2043; E2043; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.1; 420dpi; 480x784; HTC; D626x; htc_a32ul; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 408dpi; 720x1394; Huawei; U8836D; hwu8836D; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.1; 320dpi; 1080x2231; Meizu; U10; U10; mt6755; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (29/10; 440dpi; 1080x2201; samsung; SAMSUNG-SM-G900A; klteatt; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.1; 408dpi; 1080x2221; samsung; sm-a520f; a5y17lte; samsungexyn; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.1; 356dpi; 720x1332; Samsung; SM-T116NU; goyave3gsea; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.2; 440dpi; 1080x2178; LGE; AS740; aloha; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (29/10; 540dpi; 1080x2295; motorola/verizon; DROID3; cdma_solana; mapphone_cdma; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (28/9; 544dpi; 1440x2768; samsung; SM-G925F; zerolte; samsungexyno; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (27/8.1.0; 480dpi; 1334x750; Huawei; U8500; msm7225; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (28/9; 480dpi; 1080x2228; Xiaomi/POCO; M2012K11AG; alioth; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (28/9; 544dpi; 720x1185; HTC; gtou; HTC_Desire_200; hi3660; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (27/8.1.0; 280dpi; 1152x1920; Motorola; sholes; Droid; hi3660; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (28/9; 544dpi; 1440x2891; Lenovo; LenovoA396; A396; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (28/9; 190dpi; 720x1453; WIKO; RAINBOWJAM; s5250; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (31/12; 440dpi; 2560x1492; Sony; E2115; E2115; mt6582; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 120dpi; 2041x1080; Lenovo; LenovoTV40S9; jazz; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 272dpi; 1378x720; Clementoni; Clempad_HR; Clempad_HR; h1; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (29/10; 356dpi; 480x899; Huawei; T8830Pro; hwT8830Pro; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 480dpi; 1080x2304; samsung; SM-M205FN; m20lte; exynos7885; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.2; 560dpi; 1080x2080; Acer; N3-2200; da2; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1; 120dpi; 1440x2698; LGE; LG-E975T; geehrc; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (25/7.1.2; 640dpi; 1920x1200; motorola; moto g power; sofia; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (26/8.0.0; 240dpi; 1080x2024; Vestel; Venus E2 Plus; Ada; qcom; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (27/8.1.0; 280dpi; 1152x1920; Motorola; sholes; Droid; hi3660; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 480dpi; 1080x2194; Xiaomi; Redmi Note 4X; nikel; mt6797; ru_RU; 436384447)",
	"Instagram 265.0.0.19.301 Android (30/11; 213dpi; 1080x2131; HUAWEI; HUAWEI VNS-L22; HWVNS-H; hi6250; ru_RU; 436384447)",
}

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

func ParseBotAccounts(ctx context.Context, bots []*tasksservice.BotAccountRecord) (BotAccounts, []*tasksservice.UploadError) {
	botAccounts := make([]BotAccount, len(bots))
	var errs []*tasksservice.UploadError
	var err error

	for i, botAccountRecord := range bots {
		err = botAccounts[i].parse(ctx, botAccountRecord.Record)
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
func (b *BotAccount) parse(ctx context.Context, fields []string) error {
	err := b.assignLoginAndPassword(fields[0])
	if err != nil {
		return err
	}

	err = b.assignUserAgent(ctx, fields[1])
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

func (b *BotAccount) assignUserAgent(ctx context.Context, input string) error {
	var err error

	fallbackUserAgent := RandomFromSlice(userAgents)
	if input == "" {
		input = RandomFromSlice(userAgents)
	}

	b.UserAgent = input
	b.DeviceData, err = headers.NewDeviceSettings(ctx, input, fallbackUserAgent)
	if err != nil {
		return fmt.Errorf("failed to parse device settings: %v", err)
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

	b.InstID, err = strconv.ParseInt(b.Headers.DsUserID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse instagram user id from '%s': %v", b.Headers.DsUserID, err)
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
