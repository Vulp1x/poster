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
	EditMediaTaskKind             = 2
	TransitToSimilarFoundTaskKind = 4
)

var EmptyPayload = []byte(`{"empty":true}`)

func NewQueuue(ctx context.Context, executor executor.Executor, txFunc dbmodel.DBTXFunc, conn *grpc.ClientConn) *pgqueue.Queue {
	queue := pgqueue.NewQueue(ctx, executor)

	// выкладывание новых постов
	queue.RegisterKind(MakePhotoPostsTaskKind, &PostPhotoHandler{dbTxF: txFunc, cli: api.NewInstaProxyClient(conn), queue: queue}, pgqueue.KindOptions{
		Name:                 "post-photo",
		WorkerCount:          pgqueue.NewConstProvider(int16(100)),
		MaxAttempts:          7,
		AttemptTimeout:       100 * time.Second,
		MaxTaskErrorMessages: 10,
		Delayer:              delayer.NewJitterDelayer(delayer.EqualJitter, 20*time.Second),
		TerminalTasksTTL:     pgqueue.NewConstProvider(1000 * time.Hour),
		Loop: pgqueue.LoopOptions{
			JanitorPeriod: pgqueue.NewConstProvider(15 * time.Hour),
			FetcherPeriod: pgqueue.NewConstProvider(3 * time.Second),
		},
	})

	// ищем похожих блогеров на начальных блогеров
	queue.RegisterKind(EditMediaTaskKind, &EditPhotoHandler{dbTxF: txFunc, cli: api.NewInstaProxyClient(conn), queue: queue}, pgqueue.KindOptions{
		Name:                 "edit-old-posts",
		WorkerCount:          pgqueue.NewConstProvider(int16(2)),
		MaxAttempts:          5,
		AttemptTimeout:       30 * time.Second,
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
