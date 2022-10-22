package client

import (
	"azuki774/dropbox-uploader/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type NewTokenClient struct {
	Client       *http.Client
	RefreshToken string
	AppKey       string
	AppSecret    string
}

func (n *NewTokenClient) Do() (resp model.RefreshResponse, err error) {
	endpoint := "https://api.dropbox.com/oauth2/token"
	reqbody := fmt.Sprintf("refresh_token=%s&grant_type=refresh_token", n.RefreshToken)
	reader := strings.NewReader(reqbody)

	req, err := http.NewRequest("POST", endpoint, reader)
	req.SetBasicAuth(n.AppKey, n.AppSecret)
	res, err := n.Client.Do(req)
	if err != nil {
		return model.RefreshResponse{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return model.RefreshResponse{}, fmt.Errorf("unexpected status code: %v", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return model.RefreshResponse{}, err
	}

	resp.RefreshToken = n.RefreshToken
	return resp, nil
}
