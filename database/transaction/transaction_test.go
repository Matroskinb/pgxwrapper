package transaction_test

import (
	"context"
	"testing"

	"github.com/Matroskinb/pgxwrapper/database"
	"github.com/Matroskinb/pgxwrapper/database/transaction"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultipleQueriesWithWrapTx(t *testing.T) {
	db, err := database.NewForURL("postgresql://postgres:postgres@localhost:5432/testing", nil)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = db.ExecuteRaw(ctx, `drop table if exists testing_transactions`)
	require.NoError(t, err)

	_, err = db.ExecuteRaw(ctx, `
		create table testing_transactions (
			id varchar(10) not null constraint testing_transactions_pk primary key
		);
	`)
	require.NoError(t, err)

	err = db.WithTx(ctx, func(innerCtx context.Context, queryRunner *transaction.Executor) error {
		err = db.WithTx(innerCtx, func(innerCtx context.Context, queryRunner *transaction.Executor) error {
			_, err = queryRunner.Insert(innerCtx, db.QueryBuilder().Insert("testing_transactions").Columns("id").Values("first"))

			return err
		})

		require.NoError(t, err)

		err = db.WithTx(innerCtx, func(innerCtx context.Context, queryRunner *transaction.Executor) error {
			_, err = queryRunner.Insert(innerCtx, db.QueryBuilder().Insert("testing_transactions").Columns("id").Values("second"))

			return err
		})

		require.NoError(t, err)

		count, err := database.Get[int](ctx, db, db.QueryBuilder().Select("count(*)").From("testing_transactions"))
		require.NoError(t, err)

		assert.Equal(t, 2, *count)

		return nil
	})
	require.NoError(t, err)
}

func TestMultipleQueriesWithTx(t *testing.T) {
	db, err := database.NewForURL("postgresql://postgres:postgres@localhost:5432/testing", nil)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = db.ExecuteRaw(ctx, `drop table if exists testing_transactions`)
	require.NoError(t, err)

	_, err = db.ExecuteRaw(ctx, `
		create table testing_transactions (
			id varchar(10) not null constraint testing_transactions_pk primary key
		);
	`)
	require.NoError(t, err)

	err = db.WithTx(ctx, func(ctx context.Context, queryRunner *transaction.Executor) error {
		_, err = queryRunner.Insert(ctx, db.QueryBuilder().Insert("testing_transactions").Columns("id").Values("first"))

		return err
	})

	require.NoError(t, err)

	err = db.WithTx(ctx, func(ctx context.Context, queryRunner *transaction.Executor) error {
		_, err = queryRunner.Insert(ctx, db.QueryBuilder().Insert("testing_transactions").Columns("id").Values("second"))

		return err
	})

	require.NoError(t, err)

	var count int

	err = db.GetRaw(ctx, &count, "select count(*) from testing_transactions")
	require.NoError(t, err)

	assert.Equal(t, 2, count)
}
