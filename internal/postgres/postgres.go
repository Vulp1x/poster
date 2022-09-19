package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	global  *pgxpool.Pool
	once    sync.Once
	onceErr error
)

// Load loads default database once and returns validation function.
func Load(ctx context.Context, cfg Configuration) func() error {
	v := func() error {
		if onceErr != nil {
			return onceErr
		} else if global == nil {
			return errors.New("missing db function")
		}

		return nil
	}

	once.Do(func() {
		pool, err := newPool(ctx, cfg)
		if err != nil {
			onceErr = fmt.Errorf("new db function: %w", err)

			return
		}
		global = pool
	})

	return v
}

// NewTxFunc sets default database function once, returns it and returns error.
func NewTxFunc(ctx context.Context, cfg Configuration) (dbmodel.TxFunc, error) {
	var dbtxf dbmodel.TxFunc
	err := TxFunc(ctx, &dbtxf, cfg)()
	return dbtxf, err
}

// TxFunc sets default database connection function and returns validation function.
func TxFunc(ctx context.Context, dbFunc *dbmodel.TxFunc, cfg Configuration) func() error {
	if dbFunc == nil {
		return func() error { return fmt.Errorf("missing argument %T", dbFunc) }
	}

	v := func() error {
		if onceErr != nil {
			return onceErr
		} else if *dbFunc == nil {
			return errors.New("missing db function")
		}
		return nil
	}

	if *dbFunc != nil {
		return v
	}

	v = Load(ctx, cfg)
	*dbFunc = func(ctx context.Context) (dbmodel.Tx, error) {
		return global.Begin(ctx)
	}

	return v
}

// NewDBTxFunc sets default database/transition function once, returns it and returns error.
func NewDBTxFunc(ctx context.Context, cfg Configuration) (dbmodel.DBTXFunc, error) {
	var dbtxf dbmodel.DBTXFunc
	err := DBTxFunc(ctx, &dbtxf, cfg)()
	return dbtxf, err
}

// DBTxFunc sets default database function and returns validation function.
func DBTxFunc(ctx context.Context, dbTxFunc *dbmodel.DBTXFunc, cfg Configuration) func() error {
	if dbTxFunc == nil {
		return func() error { return fmt.Errorf("missing argument %T", dbTxFunc) }
	}

	v := func() error {
		if onceErr != nil {
			return onceErr
		} else if *dbTxFunc == nil {
			return errors.New("missing db function")
		}
		return nil
	}

	if *dbTxFunc != nil {
		return v
	}

	v = Load(ctx, cfg)
	*dbTxFunc = func(ctx context.Context) dbmodel.DBTX { return global }

	return v
}

// newPool returns new db function.
func newPool(ctx context.Context, cfg Configuration) (*pgxpool.Pool, error) {
	fullDsn := cfg.buildFullDsn()
	config, err := pgxpool.ParseConfig(fullDsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from url %s: %w", fullDsn, err)
	}

	conn, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	return conn, nil
}
