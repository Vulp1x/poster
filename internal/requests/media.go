package requests

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/transport"
)

func PrepareConfigureMediaRequest(ctx context.Context, b domain.BotAccount, uploadID string, image []byte) (*http.Request, error) {
	// body := createHTTPBodyFromImageBytes(image)

	data := map[string]string{
		"camera_entry_point":         "1",
		"scene_capture_type":         "",
		"timezone_offset":            "3600",
		"usertags":                   "",
		"source_type":                "4",
		"_uid":                       "2182798860",
		"device_id":                  "android-6f16a782d9e31a11",
		"_uuid":                      "d1a182d7-8663-44b1-83d1-0bb71b798282",
		"creation_logger_session_id": "56fe941f-89a8-40b4-84ec-5f0aaaedf76b",
		"nav_chain":                  "MainFeedFragment:feed_timeline:1:cold_start::,QuickCaptureFragment:stories_precapture_camera:8:swipe::,QuickCaptureFragment:clips_precapture_camera:9:button::,QuickCaptureFragment:stories_precapture_camera:10:button::,GalleryPickerFragment:gallery_picker:11:button::,PhotoFilterFragment:photo_filter:12:button::,FollowersShareFragment:media_broadcast_share:13:next::",
		"caption":                    "Это+мои+друзья+@al_kharba_+@greg.postnikov+и+Серж+(Федя+фотографирует)",
		"upload_id":                  uploadID,
		"location":                   "",
		"device":                     b.PrepareDevice(),
		"edits":                      `{"crop_original_size": [	1620.0,	2160.0],"crop_center": [0.0,-0.0],"crop_zoom": 1.3333334}`,
		"extra":                      `{"source_width": 1620,"source_height": 2160}`,
	}

	body := generateSignature(data)

	headers := b.ConstructHeaders(int64(len(image)))

	waterFallID := uuid.NewString()
	uploadName := fmt.Sprintf("%d_0_%d", uploadID, rand.Intn(8999999999)+1000000000)
	//
	ruploadParams := []string{fmt.Sprintf(`{"retry_context":"{\"num_step_auto_retry\":0,\"num_reupload\":0,\"num_step_manual_retry\":0}","media_type":"1","xsharing_user_ids": "[]","upload_id": "%d","image_compression":"{\"lib_name\": \"moz\", \"lib_version\": \"3.1.m\", \"quality\": \"76\"}"}`, uploadID)}

	headers["Accept-Encoding"] = []string{"gzip"}
	headers["X-Instagram-Rupload-Params"] = ruploadParams
	headers["X_FB_PHOTO_WATERFALL_ID"] = []string{waterFallID}
	// headers["X-Entity-Type"] = []string{"image/webp"}
	headers["X-Entity-Type"] = []string{"image/jpeg"}
	headers["Offset"] = []string{"0"}
	headers["X-Entity-Name"] = []string{uploadName}
	headers["X-Entity-Length"] = []string{strconv.FormatInt(int64(len(image)), 10)}
	headers["Content-Type"] = []string{"application/octet-stream"}
	headers["Content-Length"] = []string{strconv.FormatInt(int64(len(image)), 10)}
	// headers["x-csrftoken"] = []string{b.Headers.AuthData.CSRFToken}

	req, err := http.NewRequestWithContext(
		transport.ContextWithProxy(ctx, b.ProxyURL()),
		"POST",
		fmt.Sprintf("https://%s/rupload_igphoto/%s", instHost, uploadName),
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return req, nil
}
