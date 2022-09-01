package logger

import (
	"go.uber.org/zap"
)

func NewLogger() (Logger *zap.Logger, err error) {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, err := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	return l, err
}
