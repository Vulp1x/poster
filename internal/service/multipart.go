package service

import (
	"context"
	"fmt"
	"mime/multipart"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/reader"
)

// TasksServiceUploadFileDecoderFunc implements the multipart decoder for
// service "auth_service" endpoint "upload file". The decoder must populate the
// argument p after encoding.
func TasksServiceUploadFileDecoderFunc(mr *multipart.Reader, p **tasksservice.UploadFilePayload) error {
	// Add multipart request decoder logic here
	// payload := &tasksservice.UploadFilePayload{
	// 	TaskID: (*p).TaskID,
	// 	Token:  (*p).Token,
	// }

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		part, err := mr.NextPart()
		if err != nil {
			return fmt.Errorf("failed to get next part: %v", err)
		}

		switch part.FormName() {
		case "bots":
			(*p).Bots, errs = reader.ParseUsersList(ctx, part)
			if errs != nil {

			}
		}

	}

	p = &payload

	fmt.Println(len(bots), len(errs))

	return nil
}

// TasksServiceUploadFileEncoderFunc implements the multipart encoder for
// service "upload_service" endpoint "upload file".
func TasksServiceUploadFileEncoderFunc(mw *multipart.Writer, p *tasksservice.UploadFilePayload) error {
	// Add multipart request encoder logic here
	return nil
}
