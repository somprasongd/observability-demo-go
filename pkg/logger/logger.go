package logger

import (
	"context"
	"demo/pkg/ctxkey"

	"go.uber.org/zap"
)

var baseLogger *zap.Logger

func Init() {
	l, _ := zap.NewProduction()
	baseLogger = l.With(zap.String("app_name", "demo-app"))
}

func Default() *zap.Logger {
	return baseLogger
}

// FromContext extracts logger from context
func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ctxkey.Logger{}).(*zap.Logger)
	if !ok {
		return baseLogger
	}
	return logger
}
