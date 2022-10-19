package instagrapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/transport"
	"github.com/inst-api/poster/pkg/logger"
)

type Client struct {
	cli              *http.Client
	saveResponseFunc func(ctx context.Context, sessionID string, response *http.Response, d time.Duration) error
}

// CheckLandingAccounts проверяет все аккаунты, на которые ведем трафик, что они живы и у них в профиле есть ссылка
func (c *Client) CheckLandingAccounts(ctx context.Context, sessionID string, landingAccountUsernames []string) ([]string, error) {
	startedAt := time.Now()
	val := map[string][]string{
		"sessionid": {sessionID},
		"usernames": landingAccountUsernames,
	}

	resp, err := c.cli.PostForm("http://localhost:8000/check/landings", val)
	if err != nil {
		return nil, err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, time.Since(startedAt))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body bytes: %v", err)
	}

	var aliveLandings []string

	err = json.Unmarshal(respBytes, &aliveLandings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal alive landings: %v", err)
	}

	return aliveLandings, nil
}

func NewClient() *Client {
	return &Client{cli: transport.InitHTTPClient(), saveResponseFunc: saveResponse}
}

// MakePost создает новый
func (c *Client) MakePost(ctx context.Context, cheapProxy, sessionID, caption string, image []byte) error {
	startedAt := time.Now()
	buf, contentType, err := prepareUploadImageBody(image, sessionID, cheapProxy, caption)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/photo/upload", buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, time.Since(startedAt))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}

// InitBot создает бота в instagrapi-rest, чтобы потом отправлять от его лица запросы
func (c *Client) InitBot(ctx context.Context, bot domain.BotWithTargets) error {
	startedAt := time.Now()
	bodyBytes, err := prepareInitBody(bot)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/auth/add", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, bot.Headers.AuthData.SessionID, resp, time.Since(startedAt))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}
