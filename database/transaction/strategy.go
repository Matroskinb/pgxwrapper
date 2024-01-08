package transaction

import (
	"context"

	"github.com/Matroskinb/pgxwrapper/database/exception"
	"github.com/Matroskinb/pgxwrapper/database/query"

	"github.com/jackc/pgx/v5"
)

func NewStrategy(tx pgx.Tx) *Strategy {
	return &Strategy{tx: tx}
}

type Strategy struct {
	tx pgx.Tx
}

func (strategy *Strategy) Handle(
	ctx context.Context,
	handler func(ctx context.Context, tx *Executor) error,
) (err error) {
	if err = handler(ctx, NewTransactionExecutor(strategy.tx, query.NewExecutor(strategy.tx))); err == nil {
		err = strategy.tx.Commit(ctx)

		return exception.New(err)
	}

	_ = strategy.tx.Rollback(ctx)

	return err
}
