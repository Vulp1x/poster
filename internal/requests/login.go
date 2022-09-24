package requests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/inst-api/poster/internal/crypto"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/transport"
)

func PrepareContactPointPrefillRequest(b domain.BotAccount) *http.Request {
	body := generateSignature(map[string]string{
		"phone_id": b.Session.PhoneID.String(),
		"usage":    "prefill",
	})

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   instHost,
			Path:   "/api/v1/BotAccounts/contact_point_prefill/",
		},
		Header: b.Header(int64(body.Len())),
		Body:   io.NopCloser(body),
	}

	return req.WithContext(transport.ContextWithProxy(context.Background(), b.ProxyURL()))
}

func PrepareSyncLauncherRequest(b domain.BotAccount, login bool) *http.Request {
	data := map[string]string{
		"id":                      b.Session.UUID.String(),
		"server_config_retrieval": "1",
	}

	if !login {
		data["_uid"] = b.Headers.AuthData.DsUserID
		data["_uuid"] = b.Session.UUID.String()
		data["_csrftoken"] = b.Headers.AuthData.CSRFToken
	}

	body := generateSignature(data)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   instHost,
			Path:   "/api/v1/launcher/sync/",
		},
		Header: b.Header(int64(body.Len())),
		Body:   io.NopCloser(body),
	}

	return req.WithContext(transport.ContextWithProxy(context.Background(), b.ProxyURL()))
}

func PrepareLoginRequest(b domain.BotAccount) (*http.Request, error) {
	encPassword, err := crypto.EncryptPassword([]byte(b.Password), "", "")
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"jazoest":             generateJazoest(b.Session.PhoneID),
		"country_codes":       fmt.Sprintf("[{\"country_code\":\"%d\",\"source\":[\"default\"]}]", defaultCountryCode),
		"phone_id":            b.Session.PhoneID.String(),
		"enc_password":        encPassword,
		"username":            b.Username,
		"adid":                b.Session.AdvertisingID.String(),
		"guid":                b.Session.UUID.String(),
		"device_id":           b.Session.DeviceID,
		"google_tokens":       "[]",
		"login_attempt_count": "0",
	}

	body := generateSignature(data)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "https",
			Host:   instHost,
			Path:   "/api/v1/BotAccounts/contact_point_prefill/",
		},
		Header: b.Header(int64(body.Len())),
		Body:   io.NopCloser(body),
	}

	return req.WithContext(transport.ContextWithProxy(context.Background(), b.ProxyURL())), nil
}
