package requests

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/transport"
)

func PrepareUploadRequest(ctx context.Context, b domain.BotAccount, image []byte) (*http.Request, error) {
	// body := createHTTPBodyFromImageBytes(image)
	headers := b.ConstructHeaders(int64(len(image)))

	waterFallID := uuid.NewString()
	uploadID := time.Now().UnixMilli()
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
		bytes.NewReader(image),
	)
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return req, nil
}

func createHTTPBodyFromImageBytes(image []byte) *bytes.Buffer {
	enc := base64.StdEncoding
	const base64ImagePrefix = "data:image/jpeg;base64,"
	buf := make([]byte, len(base64ImagePrefix)+enc.EncodedLen(len(image)))

	for i, c := range base64ImagePrefix {
		buf[i] = byte(c)
	}

	enc.Encode(buf[len(base64ImagePrefix):], image)

	body := bytes.NewBuffer(buf)
	return body
}
