package bots

import (
	"context"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/instagrapi"
	"github.com/inst-api/poster/pkg/logger"
)

func NewStore(dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc, instagrapiHost string) *Store {

	return &Store{
		dbtxf: dbtxFunc,
		txf:   txFunc,
		cli:   instagrapi.NewClient(instagrapiHost),
	}
}

type Store struct {
	dbtxf dbmodel.DBTXFunc
	txf   dbmodel.TxFunc
	cli   *instagrapi.Client
}

func (s Store) FindReadyBots(ctx context.Context) (domain.BotAccounts, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	bots, err := q.FindReadyBots(ctx)
	if err != nil {
		return nil, err
	}

	domainBots := make([]domain.BotAccount, 0, len(bots))

	for i := range bots {
		domainBot := domain.BotAccount(bots[i])
		// mu := &sync.Mutex{}

		// logger.Infof(ctx, "start initing bot (%s) %d", bots[i].Username, i)
		// err = s.cli.InitBot(ctx, domain.BotWithTargets{BotAccount: domainBot})
		// if err != nil {
		// 	logger.Errorf(ctx, "failed to init bot: %v", err)
		// 	continue
		// }

		domainBots = append(domainBots, domainBot)

		// go func(m *sync.Mutex, account domain.BotAccount, bots []domain.BotAccount, i int) {
		// 	logger.Infof(ctx, "start initing bot (%s) %d", account.Username, i)
		// 	err = s.cli.InitBot(ctx, domain.BotWithTargets{BotAccount: account})
		// 	if err != nil {
		// 		logger.Errorf(ctx, "failed to init bot: %v", err)
		// 		return
		// 	}
		//
		// 	m.Lock()
		// 	bots = append(bots, account)
		// 	m.Unlock()
		//
		// }(mu, domainBot, domainBots, i)
	}

	logger.Infof(ctx, "successfully added %d/%d bots", len(domainBots), len(bots))

	return domainBots, nil
}
