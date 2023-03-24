package workers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/inst-api/poster/pkg/pgqueue"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const landingAccountPlaceholder = "@account"

type PostPhotoHandler struct {
	dbTxF         dbmodel.DBTXFunc
	cli           api.InstaProxyClient
	instagrapiCLi instagrapi.Client
	queue         *pgqueue.Queue
}

func (s *PostPhotoHandler) HandleTask(ctx context.Context, task pgqueue.Task) error {
	logger.Infof(ctx, "starting processing task %s", task.ExternalKey)

	externalKeyParts := strings.Split(task.ExternalKey, "::")
	if len(externalKeyParts) != 3 {
		return fmt.Errorf("%w: expected 3 parts in external key with separator '::' in '%s', got %d parts", pgqueue.ErrMustCancelTask, task.ExternalKey, len(externalKeyParts))
	}

	taskID, err := uuid.Parse(externalKeyParts[0])
	if err != nil {
		return fmt.Errorf("%w: failed to parse datasaet id from '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[0], err)
	}

	db := s.dbTxF(ctx)
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
		return fmt.Errorf("%w: failed to find bot username '%s': %v", pgqueue.ErrMustCancelTask, externalKeyParts[1], err)
	}

	ctx = logger.WithKV(ctx, "bot_account", bot.Username)

	if bot.Status != dbmodel.StartedBotStatus {
		return fmt.Errorf("%w: got bot status %d: expected %d (started)", pgqueue.ErrMustCancelTask, bot.Status, dbmodel.StartedTaskStatus)
	}

	logger.InfoKV(ctx, "got account for processing")

	targetsLimitForNextPost := postingTask.TargetsPerPost
	if postingTask.NeedPhotoTags && bot.PostsCount < postingTask.PhotoTagsPostsPerBot {
		targetsLimitForNextPost += postingTask.PhotoTargetsPerPost
	}

	targets, err := q.LockTargetsForTask(ctx, dbmodel.LockTargetsForTaskParams{
		TaskID: taskID,
		Limit:  int32(targetsLimitForNextPost),
	})
	if err != nil {
		return fmt.Errorf("failed to find targets for post: %v", err)
	}

	logger.Infof(ctx, "got %d targets for %d post", len(targets), bot.PostsCount+1)

	if len(postingTask.AccountProfileImages) > 0 {
		err = s.instagrapiCLi.EditProfile(
			ctx,
			"",
			bot.Headers.AuthData.SessionID,
			domain.RandomFromSlice(postingTask.AccountProfileImages),
		)
		if err != nil {
			logger.Errorf(ctx, "failed to edit profile: %v", err)
		}
	}

	if postingTask.FollowTargets {
		err = s.instagrapiCLi.FollowTargets(ctx, domain.BotWithTargets{BotAccount: domain.BotAccount(bot), Targets: targets})
		if err != nil {
			logger.Errorf(ctx, "failed to follow targets: %v", err)
		}
	}

	var cheapProxy *api.Proxy
	if bot.WorkProxy == nil {
		logger.Warnf(ctx, "bot has empty cheap proxy, so using residential for post upload")
		cheapProxy = &api.Proxy{
			Host:  bot.ResProxy.Host,
			Port:  bot.ResProxy.Port,
			Login: bot.ResProxy.Login,
			Pass:  bot.ResProxy.Pass,
		}
	} else {
		cheapProxy = &api.Proxy{
			Host:  bot.WorkProxy.Host,
			Port:  bot.WorkProxy.Port,
			Login: bot.WorkProxy.Login,
			Pass:  bot.WorkProxy.Pass,
		}
	}

	// err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.StartedBotStatus, ID: bot.ID})
	// if err != nil {
	// 	return fmt.Errorf("failed to set bot status to 'started': %v", err)
	//
	// }

	landingAccount := domain.RandomFromSlice(postingTask.LandingAccounts)
	var (
		postsDone int32
	)

	if bot.PostsCount > 0 {
		logger.Infof(ctx, "bot already has %d posts, adding new", bot.PostsCount)
		postsDone += int32(bot.PostsCount)
	}

	generator, err := images.NewRandomGammaGenerator(postingTask.Images)
	if err != nil {
		return fmt.Errorf("failed to create new random gamma generator: %v", err)
	}

	mediaTargets := preparePostCaption(postingTask, landingAccount, targets)

	resp, err := s.cli.PostPicture(ctx, &api.PostPictureRequest{
		Photo:                   generator.Next(ctx),
		UserId:                  bot.InstID,
		Caption:                 mediaTargets.Caption,
		UserTags:                mediaTargets.PhotoTargets,
		CheapProxy:              cheapProxy,
		UpdatePhotoDelaySeconds: postingTask.PhotoTagsDelaySeconds,
	})
	if err != nil {
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

		return fmt.Errorf("failed to make post: %v", err)
	}

	postsDone++

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q = dbmodel.New(tx)
	savePostedMediaParams := dbmodel.SavePostedMediaParams{
		Pk:     resp.MediaPk,
		Kind:   dbmodel.MediasKindPhoto,
		InstID: resp.MediaId,
		BotID:  bot.ID,
	}
	media, err := q.SavePostedMedia(ctx, savePostedMediaParams)
	if err != nil {
		return fmt.Errorf("failed to save created media with params %+v: %v", savePostedMediaParams, err)
	}

	markTargetsAsNotifiedParams := dbmodel.MarkTargetsAsNotifiedParams{
		MediaFk:         &media.Pk,
		InteractionType: dbmodel.TargetsInteractionPostDescription,
		TaskID:          postingTask.ID,
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

	mediasCount, err := q.CountBotMedias(ctx, bot.ID)
	if err != nil {
		return fmt.Errorf("failed to count bot medias: %v", err)
	}

	if err = q.SetBotPostsCount(ctx, dbmodel.SetBotPostsCountParams{PostsCount: int(mediasCount), ID: bot.ID}); err != nil {
		return fmt.Errorf("failed to set posts count: %v", err)
	}

	if int(mediasCount) < postingTask.PostsPerBot {
		logger.Infof(ctx, "adding one more task for posting, got %d/%d posts", mediasCount, postingTask.PostsPerBot)
		newTask := pgqueue.Task{
			Kind:        MakePhotoPostsTaskKind,
			Payload:     EmptyPayload,
			ExternalKey: fmt.Sprintf("%s::%s::%d", postingTask.ID.String(), bot.Username, mediasCount),
		}
		if err = s.queue.PushTaskTx(ctx, tx, newTask, pgqueue.WithDelay(time.Duration(postingTask.PerPostSleepSeconds)*time.Second)); err != nil {
			return fmt.Errorf("failed to push task (%#v) to queue: %v", newTask, err)
		}
	} else {
		logger.Infof(ctx, "already has %d/%d posts, don't add new posting task", mediasCount, postingTask.PostsPerBot)
		err = q.SetBotStatus(ctx, dbmodel.SetBotStatusParams{Status: dbmodel.DoneBotStatus, ID: bot.ID})
		if err != nil {
			return fmt.Errorf("failed to set bot status to %d: %v", dbmodel.EditingPostsDoneBotStatus, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

type postTargets struct {
	Caption            string
	DescriptionTargets []int64
	PhotoTargets       []int64
}

func preparePostCaption(task dbmodel.Task, landingAccount string, targets []dbmodel.TargetUser) postTargets {
	b := strings.Builder{}
	b.WriteString(strings.Replace(task.TextTemplate, landingAccountPlaceholder, "@"+landingAccount, 1))
	i := 0

	var descriptionTargetsUserIds = make([]int64, 0, task.TargetsPerPost)
	var photoTargetsUserIds = make([]int64, 0, task.PhotoTargetsPerPost)

	if task.FixedTag != nil {
		b.WriteByte(' ')
		b.WriteByte('@')
		b.WriteString(*task.FixedTag)
	}

	if task.FixedPhotoTag != nil {
		photoTargetsUserIds = append(photoTargetsUserIds, *task.FixedPhotoTag)
	}

	for _, target := range targets {
		b.WriteByte(' ')
		b.WriteByte('@')
		b.WriteString(target.Username)
		descriptionTargetsUserIds = append(descriptionTargetsUserIds, target.UserID)
		i++
		if i >= task.TargetsPerPost {
			break
		}
	}

	if task.NeedPhotoTags && i < len(targets) {
		for ; i < len(targets); i++ {
			if len(photoTargetsUserIds) == task.PhotoTargetsPerPost {
				break
			}

			photoTargetsUserIds = append(photoTargetsUserIds, targets[i].UserID)
		}
	}

	return postTargets{
		Caption:            b.String(),
		DescriptionTargets: descriptionTargetsUserIds,
		PhotoTargets:       photoTargetsUserIds,
	}
}
