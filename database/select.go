package database

import (
	"context"

	"github.com/Matroskinb/pgxwrapper/database/query"

	"github.com/Masterminds/squirrel"
)

// Select selecting data from the database.
// There is a restriction on using generic, that's why it's a method for now.
// Takes data from the database and unpacks it into the model.
func Select[T any](ctx context.Context, database *Database, builder squirrel.SelectBuilder) ([]T, error) {
	var result []T

	sql, args, err := query.Unwrap(builder)
	if err != nil {
		return result, err
	}

	err = database.SelectRaw(ctx, &result, sql, args...)

	return result, err
}

// Get fetching a row from the database.
// There is a restriction on using generic, that's why it's a method for now.
// Takes data from the database and unpacks it into the model.
func Get[T any](ctx context.Context, database *Database, builder squirrel.SelectBuilder) (*T, error) {
	result := new(T)

	sql, args, err := query.Unwrap(builder)
	if err != nil {
		return result, err
	}

	err = database.GetRaw(ctx, result, sql, args...)

	return result, err
}
