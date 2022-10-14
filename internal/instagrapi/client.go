package instagrapi

import (
	"bytes"
	"context"
	"fmt"
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

func NewClient() *Client {
	return &Client{cli: transport.InitHTTPClient(), saveResponseFunc: saveResponse}
}

// MakePost создает новый
func (c *Client) MakePost(ctx context.Context, bot domain.BotAccount, sessionID, caption string, image []byte) error {
	startedAt := time.Now()
	buf, contentType, err := prepareUploadImageBody(image, sessionID, bot.WorkProxy.PythonString(), caption)
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
