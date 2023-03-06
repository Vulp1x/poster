package workers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EditPhotoHandler struct {
	dbTxF         dbmodel.DBTXFunc
	cli           api.InstaProxyClient
	instagrapiCLi instagrapi.Client
	queue         *pgqueue.Queue
}

func (h *EditPhotoHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) != 3 {
		return fmt.Errorf("%w: expected 3 parts in external key with separator '::' in '%s', got %d parts", pgqueue.ErrMustCancelTask, task.ExternalKey, len(externalKeyParts))
	}

	taskID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	mediaSequenceNumber, err := strconv.ParseInt(externalKeyParts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("%w: failed to parse media sequence number from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[2], err)
	}

	db := h.dbTxF(ctx)
	q := dbmodel.New(db)

	postingTask, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("%w: failed to find task with id '%s': %v", pgqueue.ErrMustCancelTask, taskID, err)
	}

	bot, err := q.FindTaskBotByUsername(ctx, dbmodel.FindTaskBotByUsernameParams{
		TaskID:   taskID,
		Username: externalKeyParts[1],
	})
	if err != nil {
		return fmt.Errorf("%w: failed to find bot with id '%s': %v", pgqueue.ErrMustCancelTask, taskID, err)
	}

	ctx = logger.WithKV(ctx, "bot_account", bot.Username)

	if bot.Status != dbmodel.EditingPostsBotStatus {
		return fmt.Errorf("%w: got bot status %d: expected %d (editing posts)", pgqueue.ErrMustCancelTask, bot.Status, dbmodel.EditingPostsBotStatus)
	}

	medias, err := q.GetBotMedias(ctx, bot.ID)
	if err != nil {
		return fmt.Errorf("failed to find bot's medias: %v", err)
	}

	logger.Infof(ctx, "found %d medias", len(medias))

	var mediaToEdit dbmodel.Media

	for _, media := range medias {
		if !media.IsEdited {
			logger.Infof(ctx, "media %d isn't edited, going to edit it's caption", media.Pk)
			mediaToEdit = media
			break
		}
	}

	if mediaToEdit.Pk == 0 {
		// все медиа уже обновлены
		logger.Info(ctx, "all medias are updated, setting bot status to 8 (editing posts done)")
		err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.EditingPostsDoneBotStatus, ID: bot.ID})
		if err != nil {
			return fmt.Errorf("failed to set bot status to %d: %v", dbmodel.EditingPostsDoneBotStatus, err)
		}
	}

	targetsLimitForNextPost := postingTask.TargetsPerPost
	if postingTask.NeedPhotoTags && bot.PostsCount < postingTask.PhotoTargetsPerPost {
		targetsLimitForNextPost += postingTask.PhotoTargetsPerPost
	}

	targets, err := q.LockTargetsForTask(ctx, dbmodel.LockTargetsForTaskParams{
		TaskID: taskID,
		Limit:  int32(targetsLimitForNextPost),
	})
	if err != nil {
		return fmt.Errorf("failed to find targets for post: %v", err)
	}

	logger.Infof(ctx, "got %d targets for editing post %d", len(targets), mediaToEdit.Pk)

	landingAccount := domain.RandomFromSlice(postingTask.LandingAccounts)
	mediaTargets := preparePostCaption(postingTask, landingAccount, targets)

	if _, err = h.cli.UpdatePicture(ctx, &api.UpdatePostRequest{
		UserId:   bot.InstID,
		Caption:  mediaTargets.Caption,
		UserTags: mediaTargets.PhotoTargets,
		MediaPk:  mediaToEdit.Pk,
	}); err != nil {
		logger.Errorf(ctx, "failed to edit media %d: %v", mediaToEdit.Pk, err)

		err2 := q.RollbackTargetsStatus(ctx, domain.Ids(targets))
		if err2 != nil {
			logger.Errorf(ctx, "failed to rollback targets status: %v", err2)
		}

		if status.Code(err) == codes.PermissionDenied {
			// бот заблокирован, нужно отметить это в базе
			logger.Warnf(ctx, "going to block bot, because of instaproxy error: %v", err)
			err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.BlockedBotStatus, ID: bot.ID})
			if err != nil {
				return fmt.Errorf("failed to block bot: %v", err)
			}

			return fmt.Errorf("%w: бот заблокирован", pgqueue.ErrMustCancelTask)
		}

		return fmt.Errorf("failed to edit media %d: %v", mediaToEdit.Pk, err)
	}

	return h.postProcessMediaEdit(ctx, bot, postingTask, mediaToEdit, mediaTargets, mediaSequenceNumber+1)
}

func (h *EditPhotoHandler) postProcessMediaEdit(
	ctx context.Context,
	bot dbmodel.BotAccount,
	task dbmodel.Task,
	editedMedia dbmodel.Media,
	mediaTargets postTargets,
	nextMediaSequenceNumber int64,
) error {
	tx, err := h.dbTxF(ctx).Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	if err = q.SetMediaIsEdited(ctx, editedMedia.Pk); err != nil {
		return fmt.Errorf("failed to set is_edited=true for media %d: %v", editedMedia.Pk, err)
	}

	markTargetsAsNotifiedParams := dbmodel.MarkTargetsAsNotifiedParams{
		MediaFk:         &editedMedia.Pk,
		InteractionType: dbmodel.TargetsInteractionPostDescription,
		TaskID:          task.ID,
		TargetIds:       mediaTargets.DescriptionTargets,
	}

	err = q.MarkTargetsAsNotified(ctx, markTargetsAsNotifiedParams)
	if err != nil {
		return fmt.Errorf("failed to set targets statuses to 'notified' for description tagged targets '%#v': %v", markTargetsAsNotifiedParams, err)
	}

	if len(mediaTargets.PhotoTargets) > 0 {
		markTargetsAsNotifiedParams.TargetIds = mediaTargets.PhotoTargets
		markTargetsAsNotifiedParams.InteractionType = dbmodel.TargetsInteractionPhotoTag
		err = q.MarkTargetsAsNotified(ctx, markTargetsAsNotifiedParams)
		if err != nil {
			return fmt.Errorf("failed to set targets statuses to 'notified' for photo tagged targets '%#v': %v", markTargetsAsNotifiedParams, err)
		}
	}

	newTask := pgqueue.Task{
		Kind:        EditMediaTaskKind,
		Payload:     EmptyPayload,
		ExternalKey: fmt.Sprintf("%s::%s::%d", task.ID.String(), bot.Username, nextMediaSequenceNumber),
	}
	if err = h.queue.PushTaskTx(ctx, tx, newTask, pgqueue.WithDelay(2*time.Duration(task.PerPostSleepSeconds)*time.Second)); err != nil {
		return fmt.Errorf("failed to push task (%#v) to queue: %v", newTask, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
