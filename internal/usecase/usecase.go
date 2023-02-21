package usecase

import "go.uber.org/zap"

type Usecase struct {
	Logger     *zap.Logger
	Client     Client
	SrcRootDir string
	DstRootDir string
}

type Client interface {
}
