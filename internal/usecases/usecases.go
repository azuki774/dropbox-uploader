package usecases

import (
	"azuki774/dropbox-uploader/internal/model"
	"io"

	"go.uber.org/zap"
)

type Client interface {
	RenewClient(newToken string)
	UploadFile(path string, content io.Reader) (err error)
}

type NewTokenClient interface {
	Do() (resp model.RefreshResponse, err error)
}

type Usecases struct {
	Logger         *zap.Logger
	Client         Client
	NewTokenClient NewTokenClient
}

func (u *Usecases) GetNewAccessToken() (err error) {
	resp, err := u.NewTokenClient.Do()
	if err != nil {
		u.Logger.Error("failed to fetch new token", zap.Error(err))
		return err
	}

	u.Client.RenewClient(resp.AccessToken)
	u.Logger.Info("update new access token", zap.String("access_token", resp.AccessToken))
	return nil
}

func (u *Usecases) UploadFile(path string, file io.Reader) (err error) {
	err = u.Client.UploadFile(path, file)
	if err != nil {
		u.Logger.Error("failed to upload file", zap.Error(err))
		return err
	}
	return nil
}
