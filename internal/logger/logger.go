package logger

import (
	"go.uber.org/zap"
)

type UploadOptions struct {
	Logger      *zap.Logger
	SrcDir      string // File or Directory
	DstDir      string // Dropbox directory
	OverWrite   bool
	AccessToken string
}

func NewLogger() (Logger *zap.Logger, err error) {
	l, err := zap.NewProduction()
	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	return l, err
}
