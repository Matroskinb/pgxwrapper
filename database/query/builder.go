package query

import (
	"github.com/Masterminds/squirrel"
)

func NewBuilder() Builder {
	builder := Builder{}
	builder.boot()

	return builder
}

type Builder struct {
	squirrel.StatementBuilderType
}

func (q *Builder) boot() {
	q.StatementBuilderType = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
