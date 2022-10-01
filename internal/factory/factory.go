package factory

import (
	"azuki774/dropbox-uploader/internal/client"
	"azuki774/dropbox-uploader/internal/server"
	"azuki774/dropbox-uploader/internal/usecases"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func NewUsecases(logger *zap.Logger, client client.Client, nclient client.NewTokenClient) *usecases.Usecases {
	return &usecases.Usecases{Logger: logger, Client: &client, NewTokenClient: &nclient}
}

func NewClient() client.Client {
	return client.Client{}
}

func NewNewTokenClient() (c client.NewTokenClient, err error) {
	// Refresh Token
	reftoken, ok := os.LookupEnv("REFRESH_TOKEN")
	if !ok {
		return c, fmt.Errorf("failed to load REFRESH_TOKEN")
	}

	// App Key
	appKey, ok := os.LookupEnv("APP_KEY")
	if !ok {
		return c, fmt.Errorf("failed to load APP_KEY")
	}

	// App Secret
	appSecret, ok := os.LookupEnv("APP_SECRET")
	if !ok {
		return c, fmt.Errorf("failed to load APP_SECRET")
	}

	c = client.NewTokenClient{Client: &http.Client{}, RefreshToken: reftoken, AppKey: appKey, AppSecret: appSecret}
	return c, nil
}

func NewServer(l *zap.Logger, us *usecases.Usecases) server.Server {
	return server.Server{Host: "", Port: "8080", Logger: l, Usecase: us}
}
