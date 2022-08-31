package logger

import (
	"context"
	"go.uber.org/zap"
)

var logger *zap.Logger

func New(ctx context.Context) *zap.SugaredLogger {
	return logger.Sugar()
}

func Init() *zap.Logger {
	logger = zap.Must(zap.NewDevelopment())
	return logger
}
