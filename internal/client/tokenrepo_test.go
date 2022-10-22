package client

import (
	"azuki774/dropbox-uploader/internal/model"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

var (
	testEndpoint         = "http://example.com:8080/oauth2/"
	testEndpointError    = "http://error.example.com:8080/oauth2/"
	testEndpointGet      = "http://example.com:8080/oauth2/dropbox"
	testEndpointGetError = "http://error.example.com:8080/oauth2/dropbox"
)

func TestMain(m *testing.M) {
	testTokenRepoGet := model.OAuth2Get{
		TokenName:    "tokenname",
		AccessToken:  "accesstoken",
		RefreshToken: "refreshtoken",
		ExpiredAt:    time.Date(2020, 1, 23, 6, 0, 0, 0, time.Local),
	}
	// testTokenRepoGet1Reader := strings.NewReader(testTokenRepoGet1)
	// testTokenRepoGet1ReaderCloser = ioutil.NopCloser(testTokenRepoGet1Reader)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", testEndpointGet,
		func(req *http.Request) (*http.Response, error) {
			res, err := httpmock.NewJsonResponse(200, testTokenRepoGet)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return res, nil
		},
	)
	httpmock.RegisterResponder("GET", testEndpointGetError,
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(500, "error")
			return res, nil
		},
	)
	httpmock.RegisterResponder("PUT", testEndpoint,
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(200, "OK")
			return res, nil
		},
	)
	httpmock.RegisterResponder("PUT", testEndpointError,
		func(req *http.Request) (*http.Response, error) {
			res := httpmock.NewStringResponse(500, "error")
			return res, nil
		},
	)

	m.Run()
}

func TestTokenRepo_Get(t *testing.T) {
	type fields struct {
		Host string
		Port string
	}
	tests := []struct {
		name    string
		fields  fields
		want    model.OAuth2Get
		wantErr bool
	}{
		{
			name: "GET",
			fields: fields{
				Host: "example.com",
				Port: "8080",
			},
			want: model.OAuth2Get{
				TokenName:    "tokenname",
				AccessToken:  "accesstoken",
				RefreshToken: "refreshtoken",
				ExpiredAt:    time.Date(2020, 1, 23, 6, 0, 0, 0, time.Local),
			},
			wantErr: false,
		},
		{
			name: "Error (500)",
			fields: fields{
				Host: "error.example.com",
				Port: "8080",
			},
			want:    model.OAuth2Get{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TokenRepo{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			got, err := tr.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenRepo_Notify(t *testing.T) {
	type fields struct {
		Host string
		Port string
	}
	type args struct {
		oa model.OAuth2Update
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				Host: "example.com",
				Port: "8080",
			},
			args: args{
				oa: model.OAuth2Update{
					TokenName:    "tokenname",
					AccessToken:  "accesstoken",
					RefreshToken: "refreshtoken",
					ExpiredIn:    100,
					RefreshURL:   "http://refresh.example.com/",
					RequestAt:    time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				Host: "error.example.com",
				Port: "8080",
			},
			args: args{
				oa: model.OAuth2Update{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TokenRepo{
				Host: tt.fields.Host,
				Port: tt.fields.Port,
			}
			if err := tr.Notify(tt.args.oa); (err != nil) != tt.wantErr {
				t.Errorf("TokenRepo.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
