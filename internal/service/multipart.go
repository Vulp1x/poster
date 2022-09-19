package service

import (
	"mime/multipart"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
)

// TasksServiceUploadFileDecoderFunc implements the multipart decoder for
// service "auth_service" endpoint "upload file". The decoder must populate the
// argument p after encoding.
func TasksServiceUploadFileDecoderFunc(mr *multipart.Reader, p **tasksservice.UploadFilePayload) error {
	// Add multipart request decoder logic here
	return nil
}

// TasksServiceUploadFileEncoderFunc implements the multipart encoder for
// service "upload_service" endpoint "upload file".
func TasksServiceUploadFileEncoderFunc(mw *multipart.Writer, p *tasksservice.UploadFilePayload) error {
	// Add multipart request encoder logic here
	return nil
}
