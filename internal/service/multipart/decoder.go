package multipart

import (
	"context"
	"fmt"
	"mime/multipart"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/pkg/logger"
)

const (
	botsPartName               = "bots"
	residentialProxiesPartName = "res_proxies"
	cheapProxiesPartName       = "cheap_proxies"
	targetUsersPartName        = "target_users"
)

// TasksServiceUploadFileDecoderFunc implements the multipart decoder for
// service "auth_service" endpoint "upload file". The decoder must populate the
// argument p after encoding.
func TasksServiceUploadFileDecoderFunc(mr *multipart.Reader, p **tasksservice.UploadFilesPayload) error {
	// Add multipart request decoder logic here
	payload := &tasksservice.UploadFilesPayload{Filenames: &tasksservice.TaskFileNames{}}

	ctx := context.Background()

	for i := 0; i < 4; i++ {
		part, err := mr.NextPart()
		if err != nil {
			return fmt.Errorf("failed to get next part: %v", err)
		}

		switch part.FormName() {
		case botsPartName:
			payload.Bots, err = readUsersList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read users list: %v", err)
			}
			payload.Filenames.BotsFilename = part.FileName()
		case residentialProxiesPartName:
			payload.ResidentialProxies, err = readProxiesList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read proxies list: %v", err)
			}
			payload.Filenames.ResidentialProxiesFilename = part.FileName()
		case cheapProxiesPartName:
			payload.CheapProxies, err = readProxiesList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read proxies list: %v", err)
			}
			payload.Filenames.CheapProxiesFilename = part.FileName()
		case targetUsersPartName:
			payload.Targets, err = readTargetsList(ctx, part)
			if err != nil {
				return fmt.Errorf("failed to read targets list: %v", err)
			}
			payload.Filenames.TargetsFilename = part.FileName()
		default:
			return fmt.Errorf("unknown part '%s' expected one of: [%s, %s, %s, %s]", part.FormName(),
				botsPartName, residentialProxiesPartName, cheapProxiesPartName, targetUsersPartName)
		}
	}

	logger.Infof(ctx, "read %d bots, %d residential proxies, %d cheap proxies, %d targets",
		len(payload.Bots), len(payload.ResidentialProxies), len(payload.CheapProxies), len(payload.Targets),
	)

	*p = payload

	return nil
}
