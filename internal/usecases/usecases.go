package usecases

import (
	"azuki774/dropbox-uploader/internal/model"

	"go.uber.org/zap"
)

type Client interface {
	RenewClient(newToken string)
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
	u.Logger.Info("get new access token", zap.String("access_token", resp.AccessToken))
	return nil
}
