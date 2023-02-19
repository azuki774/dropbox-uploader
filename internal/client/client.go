package client

import (
	"azuki774/dropbox-uploader/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"go.uber.org/zap"
)

type Client struct {
	Logger zap.Logger
	// SrcRootDir   string -> usecase?
	// DstRootDir   string
	RefreshToken string
	AppKey       string
	AppSecret    string

	filesClient files.Client
	accessToken string
}

// fetchNewRefreshToken は 新しい accessTokenを取得する
func (c *Client) fetchNewRefreshToken() (accessToken string, err error) {
	endpoint := "https://api.dropbox.com/oauth2/token"
	reqbody := fmt.Sprintf("refresh_token=%s&grant_type=refresh_token", c.RefreshToken)
	reader := strings.NewReader(reqbody)

	req, err := http.NewRequest("POST", endpoint, reader)
	req.SetBasicAuth(c.AppKey, c.AppSecret)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %v", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var resp model.ResponseAuthTokenRefreshToken
	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

// newFilesClient は dropbox のfiles API 用の SDK Client を作成する
func (c *Client) newFilesClient() files.Client {
	config := dropbox.Config{
		Token:    c.accessToken,
		LogLevel: dropbox.LogInfo, // if needed, set the desired logging level. Default is off
	}
	fileClient := files.New(config)
	return fileClient
}

// UploadFile はファイルをアップロードする
func (c *Client) UploadFile(path string, content io.Reader) (err error) {
	if c.accessToken == "" { // AccessToken が未設定ならば、RefreshToken を使った先に取得する
		c.Logger.Info("try to fetch new access token")
		c.accessToken, err = c.fetchNewRefreshToken()
		if err != nil {
			c.Logger.Error("failed to fetch new access token")
			return err
		}
		c.Logger.Info("fetch new access token sucessfully")
	}

	if c.filesClient == nil {
		c.Logger.Info("create files API client")
		c.filesClient = c.newFilesClient()
	}

	arg := files.NewUploadArg(path)
	_, err = c.filesClient.Upload(arg, content)
	if err != nil {
		return err
	}
	return nil
}
