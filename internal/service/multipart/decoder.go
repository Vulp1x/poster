package multipart

import (
	"context"
	"fmt"
	"mime/multipart"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/pkg/logger"
)

// TasksServiceUploadFileDecoderFunc implements the multipart decoder for
// service "auth_service" endpoint "upload file". The decoder must populate the
// argument p after encoding.
func TasksServiceUploadFileDecoderFunc(mr *multipart.Reader, p **tasksservice.UploadFilesPayload) error {
	// Add multipart request decoder logic here
	payload := &tasksservice.UploadFilesPayload{Filenames: &tasksservice.TaskFileNames{}}

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		part, err := mr.NextPart()
		if err != nil {
			return fmt.Errorf("failed to get next part: %v", err)
		}

		switch part.FormName() {
		case "bots":
			payload.Bots, err = readUsersList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read users list: %v", err)
			}
			payload.Filenames.BotsFilename = part.FileName()
		case "proxies":
			payload.Proxies, err = readProxiesList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read proxies list: %v", err)
			}
			payload.Filenames.ProxiesFilename = part.FileName()
		case "target_users":
			payload.Targets, err = readTargetsList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read targets list: %v", err)
			}
			payload.Filenames.TargetsFilename = part.FileName()
		}
	}

	logger.Infof(ctx, "read %d bots, %d proxies, %d targets", len(payload.Bots), len(payload.Proxies), len(payload.Targets))

	*p = payload

	return nil
}
