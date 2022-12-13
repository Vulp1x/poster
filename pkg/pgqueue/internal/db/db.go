package db

import (
	"github.com/inst-api/poster/pkg/pgqueue/pkg/executor"
)

type Queries struct {
	executor executor.Executor
}

func New(executor executor.Executor) *Queries {
	return &Queries{executor}
}
