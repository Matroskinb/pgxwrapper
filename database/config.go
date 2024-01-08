package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
)

func NewConfig(url string, l *zap.Logger) (*pgxpool.Config, error) {
	pgxPoolCfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	tracer := tracelog.TraceLog{Logger: NewLogger(l)}
	pgxPoolCfg.ConnConfig.Tracer = &tracer

	return pgxPoolCfg, nil
}
