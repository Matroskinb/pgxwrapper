package database

import (
	"context"

	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(l *zap.Logger) Logger {
	return Logger{logger: l}
}

type Logger struct {
	logger *zap.Logger
}

func (l Logger) Log(_ context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	zapLevel, err := zapcore.ParseLevel(level.String())
	if err != nil {
		zapLevel = zapcore.DebugLevel
	}

	l.logger.Log(zapLevel, msg, zap.Any("data", data))
}
