package instagrapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/transport"
	"github.com/inst-api/poster/pkg/logger"
)

// ErrBotIsBlocked аккаунт заблокирован
var ErrBotIsBlocked = errors.New("bot account is blocked")

const landingAccountPlaceholder = "@account"

type Client struct {
	cli              *http.Client
	saveResponseFunc func(ctx context.Context, sessionID string, response *http.Response, opts ...SaveResponseOption) error
	host             string
}

func NewClient(host string) *Client {
	return &Client{cli: transport.InitHTTPClient(), saveResponseFunc: saveResponse, host: host}
}

// CheckLandingAccounts проверяет все аккаунты, на которые ведем трафик, что они живы и у них в профиле есть ссылка
func (c *Client) CheckLandingAccounts(ctx context.Context, sessionID string, landingAccountUsernames []string) ([]string, error) {
	startedAt := time.Now()
	val := map[string][]string{
		"sessionid": {sessionID},
		"usernames": landingAccountUsernames,
	}

	resp, err := c.cli.PostForm(c.host+"/user/check/landings", val)
	if err != nil {
		return nil, err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)), WithReuseResponseBody(true))
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

// MakePost создает новый
func (c *Client) MakePost(
	ctx context.Context,
	task domain.Task,
	landingAccount, sessionID, cheapProxy string,
	targets []dbmodel.TargetUser,
	postImage []byte,
) error {
	startedAt := time.Now()

	var userTags []UserTag
	var caption string

	if task.NeedPhotoTags {
		caption = c.preparePostCaption(task.TextTemplate, landingAccount, targets[:len(targets)/2])
		userTags = prepareUserTags(targets[len(targets)/2:])
	} else {
		caption = c.preparePostCaption(task.TextTemplate, landingAccount, targets)
	}

	buf, contentType, err := prepareUploadImageBody(postImage, sessionID, cheapProxy, caption, userTags)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.host+"/photo/upload", buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return ErrBotIsBlocked
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}

// MakeReelsPost выкладывает новый рилс
func (c *Client) MakeReelsPost(
	ctx context.Context,
	task domain.Task,
	landingAccount, sessionID, cheapProxy string,
	targets []dbmodel.TargetUser,
	postVideo []byte,
) error {
	startedAt := time.Now()

	if len(postVideo) == 0 {
		return fmt.Errorf("got empty video")
	}

	var caption string

	caption = c.preparePostCaption(task.TextTemplate, landingAccount, targets)

	buf, contentType, err := prepareUploadReelsBody(postVideo, sessionID, cheapProxy, caption)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.host+"/video/upload", buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return ErrBotIsBlocked
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}

// EditProfile редактирует профиль бота
func (c *Client) EditProfile(ctx context.Context, fullName, sessionID string, image []byte) error {
	startedAt := time.Now()
	buf, contentType, err := prepareEditProfileBody(image, sessionID, fullName)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.host+"/user/edit_profile", buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, sessionID, resp, WithElapsedTime(time.Since(startedAt)))
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

	req, err := http.NewRequestWithContext(ctx, "POST", c.host+"/auth/add", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, bot.Headers.AuthData.SessionID, resp, WithElapsedTime(time.Since(startedAt)))
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode == http.StatusBadRequest {
		return ErrBotIsBlocked
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}

// FollowTargets подписывается на все аккаунты из списка
func (c *Client) FollowTargets(ctx context.Context, bot domain.BotWithTargets) error {
	startedAt := time.Now()
	bodyBytes, err := prepareFollowTargetsBody(bot)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.host+"/auth/follow_targets", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	err = c.saveResponseFunc(ctx, bot.Headers.AuthData.SessionID, resp,
		WithElapsedTime(time.Since(startedAt)),
	)
	if err != nil {
		logger.Errorf(ctx, "failed to save response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got %d response code, expected 200", resp.StatusCode)
	}

	return nil
}

func (c *Client) preparePostCaption(template, landingAccount string, targetUsers []dbmodel.TargetUser) string {
	b := strings.Builder{}
	b.WriteString(strings.Replace(template, landingAccountPlaceholder, "@"+landingAccount, 1))

	for _, user := range targetUsers {
		b.WriteByte(' ')
		b.WriteByte('@')
		b.WriteString(user.Username)
	}

	return b.String()
}

func (c Client) postFormWithCtx(ctx context.Context, url string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.cli.Do(req)
}

func (c *Client) getWithCtx(ctx context.Context, url string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.cli.Do(req)
}
