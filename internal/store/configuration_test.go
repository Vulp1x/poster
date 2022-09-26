package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestWithDBTXFunc(t *testing.T) {
	dbtxFunc := func(ctx context.Context) dbmodel.DBTX {
		val := ctx.Value("test")
		return val.(dbmodel.DBTX)
	}

	type args struct {
		dbTXFunc func(context.Context) dbmodel.DBTX
	}

	tests := []struct {
		name string
		args func(t minimock.Tester) args

		inspectCfg func(t minimock.Tester, cfg Configuration)
	}{
		{
			name: "ok",
			args: func(t minimock.Tester) args {
				return args{dbTXFunc: dbtxFunc}
			},
			inspectCfg: func(t minimock.Tester, cfg Configuration) {
				val := dbmodel.DBTX(&pgxpool.Pool{})
				if dbtx := cfg.DBTXf(context.WithValue(context.Background(), "test", val)); dbtx != val {
					t.Errorf("expected nil database in context, got %v", dbtx)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)

			got1 := WithDBTXFunc(tArgs.dbTXFunc)
			var cfg Configuration
			got1(&cfg)

			tt.inspectCfg(t, cfg)
		})
	}
}

func TestWithTxFunc(t *testing.T) {
	var txFunc dbmodel.TxFunc = func(ctx context.Context) (dbmodel.Tx, error) {
		val := ctx.Value("test")
		return nil, val.(error)
	}
	type args struct {
		txFunc dbmodel.TxFunc
	}
	tests := []struct {
		name string
		args func(t minimock.Tester) args

		inspectCfg func(t minimock.Tester, cfg Configuration)
	}{
		{
			name: "ok",
			args: func(t minimock.Tester) args {
				return args{txFunc: txFunc}
			},
			inspectCfg: func(t minimock.Tester, cfg Configuration) {
				val := errors.New("no error, just some value")
				if _, testValue := cfg.TxFunc(context.WithValue(context.Background(), "test", val)); testValue != val {
					t.Errorf("expected %v in context, got %v", val, testValue)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)

			got1 := WithTxFunc(tArgs.txFunc)
			var cfg Configuration
			got1(&cfg)

			tt.inspectCfg(t, cfg)
		})
	}
}
