package model

// ResponseAuthToken : Refer to https://www.dropbox.com/developers/documentation/http/documentation: When input grant_type=refresh_token:
type ResponseAuthTokenRefreshToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
