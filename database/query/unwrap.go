package query

import (
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

var ErrInvalid = errors.New("invalid query")

func Unwrap(builder squirrel.Sqlizer) (string, []interface{}, error) {
	query, args, err := builder.ToSql()

	return query, args, errors.Wrap(err, ErrInvalid.Error())
}
