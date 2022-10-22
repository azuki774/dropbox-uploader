package model

import "time"

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type OAuth2Get struct {
	TokenName    string    `json:"token_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}
type OAuth2Update struct {
	TokenName    string    `json:"token_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredIn    int64     `json:"expired_in"`
	RefreshURL   string    `json:"refresh_url"`
	RequestAt    time.Time `json:"requested_at"`
}
