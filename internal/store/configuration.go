package store

import (
	"errors"

	"github.com/inst-api/poster/internal/dbmodel"
)

// ErrTransactionFail shows big problems with database.
var ErrTransactionFail = errors.New("couldn't start new transaction")

type Configuration struct {
	DBTXf  dbmodel.DBTXFunc
	TxFunc dbmodel.TxFunc
}

type Option func(configuration *Configuration)

func WithDBTXFunc(dbTXFunc dbmodel.DBTXFunc) Option {
	return func(configuration *Configuration) {
		configuration.DBTXf = dbTXFunc
	}
}

// WithTxFunc adds function for getting database transaction to configuration
func WithTxFunc(txFunc dbmodel.TxFunc) Option {
	return func(configuration *Configuration) {
		configuration.TxFunc = txFunc
	}
}
