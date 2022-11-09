package bots

import (
	"context"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
)

func NewStore(dbtxFunc dbmodel.DBTXFunc, txFunc dbmodel.TxFunc) *Store {

	return &Store{
		dbtxf: dbtxFunc,
		txf:   txFunc,
	}
}

type Store struct {
	dbtxf dbmodel.DBTXFunc
	txf   dbmodel.TxFunc
}

func (s Store) FindReadyBots(ctx context.Context) (domain.BotAccounts, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	bots, err := q.FindReadyBots(ctx)
	if err != nil {
		return nil, err
	}

	domainBots := make([]domain.BotAccount, len(bots))

	for i := range bots {
		domainBots[i] = domain.BotAccount(bots[i])
	}

	return domainBots, nil
}
