package transaction

import (
	"context"

	"github.com/Matroskinb/pgxwrapper/database/query"

	"github.com/jackc/pgx/v5"
)

func NewTransactionExecutor(tx pgx.Tx, queryRunner *query.Executor) *Executor {
	return &Executor{Executor: queryRunner, tx: tx}
}

type Executor struct {
	*query.Executor

	tx pgx.Tx
}

func (executor *Executor) Commit(ctx context.Context) error {
	return executor.tx.Commit(ctx)
}

func (executor *Executor) Rollback(ctx context.Context) error {
	return executor.tx.Rollback(ctx)
}
