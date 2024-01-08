package database

import (
	"context"
	"time"

	"github.com/Matroskinb/pgxwrapper/database/query"
	"github.com/Matroskinb/pgxwrapper/database/transaction"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewForURL(url string, l *zap.Logger) (*Database, error) {
	cfg, err := NewConfig(url, l)
	if err != nil {
		return nil, err
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	return &Database{pool: pool, builder: query.NewBuilder()}, err
}

type Database struct {
	pool    *pgxpool.Pool
	builder query.Builder
}

func (r Database) QueryBuilder() query.Builder {
	return r.builder
}

func (r Database) Transaction(ctx context.Context) (pgx.Tx, error) {
	return transaction.NewBuilder(r.pool).Build(ctx)
}

func (r Database) WithTx(
	ctx context.Context,
	handler func(ctx context.Context, queryRunner *transaction.Executor) error,
) error {
	tx, err := transaction.NewBuilder(r.pool).Build(ctx)
	if err != nil {
		return err
	}

	return transaction.NewStrategy(tx).Handle(ctx, handler)
}

func (r Database) Update(ctx context.Context, builder squirrel.UpdateBuilder) (rowsAffected int64, err error) {
	rowsAffected, err = query.NewExecutor(r.pool).Update(ctx, builder)

	return rowsAffected, Error(err)
}

func (r Database) Delete(ctx context.Context, builder squirrel.DeleteBuilder) (rowsAffected int64, err error) {
	rowsAffected, err = query.NewExecutor(r.pool).Delete(ctx, builder)

	return rowsAffected, Error(err)
}

func (r Database) Insert(ctx context.Context, builder squirrel.InsertBuilder) (rowsAffected int64, err error) {
	rowsAffected, err = query.NewExecutor(r.pool).Insert(ctx, builder)

	return rowsAffected, Error(err)
}

func (r Database) SelectRaw(ctx context.Context, dst any, sql string, args ...any) error {
	err := query.NewExecutor(r.pool).SelectRaw(ctx, dst, sql, args...)

	return Error(err)
}

func (r Database) GetRaw(ctx context.Context, dst any, sql string, args ...any) error {
	err := query.NewExecutor(r.pool).GetRaw(ctx, dst, sql, args...)

	return Error(err)
}

func (r Database) ExecuteRaw(ctx context.Context, sql string, args ...any) (rowsAffected int64, err error) {
	rowsAffected, err = query.NewExecutor(r.pool).ExecuteRaw(ctx, sql, args...)

	return rowsAffected, Error(err)
}
