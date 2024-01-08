package query

import (
	"context"

	"github.com/Matroskinb/pgxwrapper/database/exception"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Runner interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewExecutor(executor Runner) *Executor {
	return &Executor{executor: executor}
}

type Executor struct {
	executor Runner
}

func (r *Executor) Update(ctx context.Context, builder squirrel.UpdateBuilder) (int64, error) {
	query, args, err := Unwrap(builder)
	if err != nil {
		return 0, err
	}

	res, err := r.ExecuteRaw(ctx, query, args...)

	return res, exception.New(err)
}

func (r *Executor) Delete(ctx context.Context, builder squirrel.DeleteBuilder) (int64, error) {
	query, args, err := Unwrap(builder)
	if err != nil {
		return 0, err
	}

	res, err := r.ExecuteRaw(ctx, query, args...)

	return res, exception.New(err)
}

func (r *Executor) Insert(ctx context.Context, builder squirrel.InsertBuilder) (int64, error) {
	query, args, err := Unwrap(builder)
	if err != nil {
		return 0, err
	}

	res, err := r.ExecuteRaw(ctx, query, args...)

	return res, exception.New(err)
}

func (r *Executor) SelectRaw(ctx context.Context, dst any, query string, args ...any) error {
	err := pgxscan.Select(ctx, r.executor, dst, query, args...)

	return exception.New(err)
}

func (r *Executor) GetRaw(ctx context.Context, dst any, query string, args ...any) error {
	err := pgxscan.Get(ctx, r.executor, dst, query, args...)

	return exception.New(err)
}

func (r *Executor) ExecuteRaw(ctx context.Context, query string, args ...any) (rowsAffected int64, err error) {
	tag, err := r.executor.Exec(ctx, query, args...)
	rowsAffected = tag.RowsAffected()

	return rowsAffected, exception.New(err)
}
