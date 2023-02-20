package client

import (
	"azuki774/dropbox-uploader/internal/model"
	"testing"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var l *zap.Logger
var fetchNewRefreshTokenNormalJson httpmock.Responder

func init() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = JSTTimeEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	l, _ = config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))

	var err error
	fetchNewRefreshTokenNormalJson, err = httpmock.NewJsonResponder(200, httpmock.File("./test/fetchNewRefreshToken_normal.json"))
	if err != nil {
		panic(err)
	}
}

func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

var fetchNewAccessTokenNormal = model.ResponseAuthTokenRefreshToken{
	AccessToken: "superaccesstoken",
	TokenType:   "refresh_token",
	ExpiresIn:   14400,
}

func TestClient_fetchNewAccessToken(t *testing.T) {
	type fields struct {
		Logger       *zap.Logger
		RefreshToken string
		AppKey       string
		AppSecret    string
		filesClient  files.Client
		accessToken  string
	}
	tests := []struct {
		name            string
		fields          fields
		wantAccessToken string
		wantErr         bool
		setmock         func()
	}{
		{
			name: "ok",
			fields: fields{
				Logger:       l,
				RefreshToken: "superrefreshtoken",
				AppKey:       "appkey",
				AppSecret:    "appsecret",
				filesClient:  nil,
				accessToken:  "",
			},
			wantAccessToken: "superaccesstoken",
			wantErr:         false,
			setmock: func() {
				httpmock.RegisterResponder(
					"POST",
					"https://api.dropbox.com/oauth2/token",
					fetchNewRefreshTokenNormalJson,
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Logger:       tt.fields.Logger,
				RefreshToken: tt.fields.RefreshToken,
				AppKey:       tt.fields.AppKey,
				AppSecret:    tt.fields.AppSecret,
				filesClient:  tt.fields.filesClient,
				accessToken:  tt.fields.accessToken,
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			tt.setmock()

			gotAccessToken, err := c.fetchNewAccessToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.fetchNewRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAccessToken != tt.wantAccessToken {
				t.Errorf("Client.fetchNewRefreshToken() = %v, want %v", gotAccessToken, tt.wantAccessToken)
			}
		})
	}
}
