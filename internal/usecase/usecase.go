package usecase

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type Usecase struct {
	Logger        *zap.Logger
	Client        Client
	RemoteRootDir string
}

type Client interface {
	UploadFile(ctx context.Context, srcFile string, remoteFile string) (err error)
}

// UploadDirectory は targetDir内のファイルを dropbox へ転送する
func (u *Usecase) UploadDirectory(ctx context.Context, targetDir string) (err error) {
	// targetDir 内のファイルリストを取得する
	_, err = getTargetFileNames(targetDir)

	// TODO: 1ファイルずつ dropbox に upload する

	return nil
}

// getTargetFileNames は targetDir 内のファイルリストを取得する
func getTargetFileNames(targetDir string) (targetFiles []string, err error) {
	// targetDir の存在チェック
	_, err = os.Stat(targetDir)
	if err != nil {
		return []string{}, err
	}

	err = filepath.Walk(targetDir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			targetFiles = append(targetFiles, p)
		}
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return targetFiles, nil
}
