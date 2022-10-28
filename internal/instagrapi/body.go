package instagrapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/textproto"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/headers"
)

type addNewBotBody struct {
	SessionID      string                 `json:"session_id"`
	Uuids          UUIDs                  `json:"uuids"`
	DeviceSettings headers.DeviceSettings `json:"device_settings"`
	UserAgent      string                 `json:"user_agent"`
	Proxy          string                 `json:"proxy"`
	TargetUserIDs  []int                  `json:"target_user_ids"` // на них бот будет подписываться перед запуском
}

type UUIDs struct {
	AndroidID     string    `json:"android_id"`
	PhoneID       uuid.UUID `json:"phone_id"`
	UUID          uuid.UUID `json:"uuid"`
	AdvertisingID uuid.UUID `json:"advertising_id"`
}

func prepareInitBody(botAccount domain.BotWithTargets) ([]byte, error) {
	targetIDs := make([]int, len(botAccount.Targets))
	for i, targetUser := range botAccount.Targets {
		targetIDs[i] = int(targetUser.UserID)
	}

	body := addNewBotBody{
		SessionID: botAccount.Headers.AuthData.SessionID,
		Uuids: UUIDs{
			AndroidID:     botAccount.Session.DeviceID,
			PhoneID:       botAccount.Session.PhoneID,
			UUID:          botAccount.Session.UUID,
			AdvertisingID: botAccount.Session.AdvertisingID,
		},
		DeviceSettings: botAccount.DeviceData,
		UserAgent:      botAccount.UserAgent,
		Proxy:          botAccount.ResProxy.PythonString(),
		TargetUserIDs:  targetIDs,
	}

	return json.Marshal(body)
}

func prepareFollowTargetsBody(botAccount domain.BotWithTargets) ([]byte, error) {
	targetIDs := make([]int, len(botAccount.Targets))
	for i, targetUser := range botAccount.Targets {
		targetIDs[i] = int(targetUser.UserID)
	}

	body := struct {
		SessionID string `json:"session_id"`
		TargetIDs []int  `json:"target_user_ids"`
	}{
		SessionID: botAccount.Headers.AuthData.SessionID,
		TargetIDs: targetIDs,
	}

	return json.Marshal(body)
}

func prepareUploadImageBody(image []byte, sessionID, cheapProxy, caption string, tags []UserTag) (*bytes.Buffer, string, error) {
	buf := bytes.NewBuffer(nil)
	mpWriter := multipart.NewWriter(buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="IMG_%04d.jpeg"`, rand.Intn(10000)))
	h.Set("Content-Type", "image/jpeg")

	fileForm, err := mpWriter.CreatePart(h)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create new multipart file part: %v", err)
	}

	_, err = fileForm.Write(image)
	if err != nil {
		return nil, "", fmt.Errorf("failed to write image to file: %v", err)
	}

	sessionWriter, err := mpWriter.CreateFormField("sessionid")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form field sessionid: %v", err)
	}

	_, err = sessionWriter.Write([]byte(sessionID))
	if err != nil {
		return nil, "", fmt.Errorf("failed to write session id part: %v", err)
	}

	captionWriter, err := mpWriter.CreateFormField("caption")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form field caption: %v", err)
	}

	_, err = captionWriter.Write([]byte(caption))
	if err != nil {
		return nil, "", fmt.Errorf("failed to write caption id part: %v", err)
	}

	cheapProxyWriter, err := mpWriter.CreateFormField("cheap_proxy")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form field caption: %v", err)
	}

	_, err = cheapProxyWriter.Write([]byte(cheapProxy))
	if err != nil {
		return nil, "", fmt.Errorf("failed to write caption id part: %v", err)
	}

	if len(tags) != 0 {
		userTagsWriter, err := mpWriter.CreateFormField("usertags")
		if err != nil {
			return nil, "", fmt.Errorf("failed to create form field caption: %v", err)
		}

		userTagsBytes, err := json.Marshal(tags)
		if err != nil {
			return nil, "", fmt.Errorf("failed to marshal user tags: %v", err)
		}

		_, err = userTagsWriter.Write(userTagsBytes)
		if err != nil {
			return nil, "", fmt.Errorf("failed to write caption id part: %v", err)
		}
	}

	err = mpWriter.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close multi-part writer: %v", err)
	}

	return buf, mpWriter.FormDataContentType(), nil
}

func prepareEditProfileBody(image []byte, sessionID, fullName string) (*bytes.Buffer, string, error) {
	buf := bytes.NewBuffer(nil)
	mpWriter := multipart.NewWriter(buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="IMG_%04d.jpeg"`, rand.Intn(10000)))
	h.Set("Content-Type", "image/jpeg")

	fileForm, err := mpWriter.CreatePart(h)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create new multipart file part: %v", err)
	}

	_, err = fileForm.Write(image)
	if err != nil {
		return nil, "", fmt.Errorf("failed to write image to file: %v", err)
	}

	sessionWriter, err := mpWriter.CreateFormField("sessionid")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form field sessionid: %v", err)
	}

	_, err = sessionWriter.Write([]byte(sessionID))
	if err != nil {
		return nil, "", fmt.Errorf("failed to write session id part: %v", err)
	}

	captionWriter, err := mpWriter.CreateFormField("full_name")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form field fullName: %v", err)
	}

	_, err = captionWriter.Write([]byte(fullName))
	if err != nil {
		return nil, "", fmt.Errorf("failed to write fullName id part: %v", err)
	}

	err = mpWriter.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close multi-part writer: %v", err)
	}

	return buf, mpWriter.FormDataContentType(), nil
}
