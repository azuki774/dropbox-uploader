package usecases

import (
	"azuki774/dropbox-uploader/internal/model"
	"io"
	"time"

	"go.uber.org/zap"
)

const (
	endpoint  = "https://api.dropbox.com/oauth2/token"
	tokenName = "dropbox"
)

type Client interface {
	RenewClient(newToken string)
	UploadFile(path string, content io.Reader) (err error)
}

type NewTokenClient interface {
	Do() (resp model.RefreshResponse, err error)
}

type TokenRepo interface {
	Notify(model.OAuth2Update) error
	Get() (model.OAuth2Get, error)
}

type Usecases struct {
	Logger         *zap.Logger
	Client         Client
	NewTokenClient NewTokenClient
	TokenRepo      TokenRepo
}

func (u *Usecases) GetNewAccessToken() (err error) {
	t := time.Now() // use requested_at
	resp, err := u.NewTokenClient.Do()
	if err != nil {
		u.Logger.Error("failed to fetch new token", zap.Error(err))
		return err
	}

	u.Client.RenewClient(resp.AccessToken)

	// Notify to Repo
	oa := model.OAuth2Update{
		TokenName:    tokenName,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiredIn:    int64(resp.ExpiresIn),
		RefreshURL:   endpoint,
		RequestAt:    t,
	}
	err = u.TokenRepo.Notify(oa)
	if err != nil {
		u.Logger.Error("failed to notify to repository", zap.Error(err))
		return err
	}

	u.Logger.Info("update new access token", zap.String("access_token", resp.AccessToken))
	return nil
}

func (u *Usecases) UploadFile(path string, file io.Reader) (err error) {
	err = u.Client.UploadFile(path, file)
	if err != nil {
		u.Logger.Error("failed to upload file", zap.Error(err))
		return err
	}
	u.Logger.Info("upload file", zap.String("path", path))
	return nil
}
