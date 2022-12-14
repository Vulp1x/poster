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
	queue.RegisterKind(MakePhotoPostsTaskKind, &PostPhotoHandler{dbTxF: txFunc, cli: api.NewInstaProxyClient(conn)}, pgqueue.KindOptions{
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

	queue.RegisterKind()

	queue.Start()

	return queue
}
