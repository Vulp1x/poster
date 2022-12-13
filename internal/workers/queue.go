package workers

import (
	"context"
	"time"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/pgqueue"
	"github.com/inst-api/poster/pkg/pgqueue/pkg/delayer"
	"github.com/inst-api/poster/pkg/pgqueue/pkg/executor"
	"google.golang.org/grpc"
)

const (
	MakePhotoPostsTaskKind        = 1
	ParseBloggersMediaTaskKind    = 2
	ParseUsersFromMediaTaskKind   = 3
	TransitToSimilarFoundTaskKind = 4
)

var EmptyPayload = []byte(`{"empty":true}`)

func NewQueuue(ctx context.Context, executor executor.Executor, txFunc dbmodel.DBTXFunc, conn *grpc.ClientConn) *pgqueue.Queue {
	queue := pgqueue.NewQueue(ctx, executor)

	// ищем похожих блогеров на начальных блогеров
	queue.RegisterKind(MakePhotoPostsTaskKind, &SimilarBloggersHandler{dbTxF: txFunc, cli: api.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "similar-bloggers",
		WorkerCount:          pgqueue.NewConstProvider(int16(5)),
		MaxAttempts:          10,
		AttemptTimeout:       40 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 20*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(5 * time.Second),
		},
	})

	// ищем посты для дальнейшего парсинга
	queue.RegisterKind(ParseBloggersMediaTaskKind, &ParseMediasHandler{dbTxF: txFunc, cli: api.NewInstaProxyClient(conn), queue: queue}, pgqueue.KindOptions{
		Name:                 "find-medias",
		WorkerCount:          pgqueue.NewConstProvider(int16(10)),
		MaxAttempts:          10,
		AttemptTimeout:       30 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 15*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(5 * time.Second),
		},
	})

	// парсим комментаторов из конкретного поста у блоггера
	queue.RegisterKind(ParseUsersFromMediaTaskKind, &ParseUsersFromMediaHandler{dbTxF: txFunc, cli: api.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
		Name:                 "parse-targets",
		WorkerCount:          pgqueue.NewConstProvider(int16(40)),
		MaxAttempts:          10,
		AttemptTimeout:       30 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 15*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(3 * time.Second),
		},
	})

	// переводим датасет в статус парсинг закончен после того, как все блогеры распаршены
	queue.RegisterKind(TransitToSimilarFoundTaskKind, &TransitToSimilarFoundHandler{dbTxF: txFunc}, pgqueue.KindOptions{
		Name:                 "similar-found",
		WorkerCount:          pgqueue.NewConstProvider(int16(5)),
		MaxAttempts:          10,
		AttemptTimeout:       20 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 20*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(10 * time.Second),
		},
	})

	queue.Start()

	return queue
}
