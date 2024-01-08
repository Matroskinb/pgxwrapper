package transaction

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewBuilder(pool *pgxpool.Pool) Builder {
	return Builder{pool: pool}
}

type Builder struct {
	pool *pgxpool.Pool
}

func (t Builder) Build(ctx context.Context) (tx pgx.Tx, err error) {
	tx, exists := t.fromCtx(ctx)
	if !exists {
		tx, err = t.pool.Begin(ctx)
		if err != nil {
			return tx, err
		}

		ctx = t.toCtx(ctx, tx)
	}

	return tx, nil
}

func (t Builder) toCtx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, transactionCtxKey{}, tx)
}

func (t Builder) fromCtx(ctx context.Context) (pgx.Tx, bool) {
	val := ctx.Value(transactionCtxKey{})
	if val == nil {
		return nil, false
	}

	tx, casted := val.(pgx.Tx)

	return tx, casted
}

type transactionCtxKey struct{}
