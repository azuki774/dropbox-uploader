package dropbox

import (
	"azuki774/dropbox-uploader/internal/logger"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
)

var lg *zap.Logger

const (
	MockUploadEndPoint = "https://mock.content.dropboxapi.com/2/files/upload"
	MockUploadResponse = `{
		"client_modified": "2015-05-12T15:50:38Z",
		"content_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"file_lock_info": {
			"created": "2015-05-12T15:50:38Z",
			"is_lockholder": true,
			"lockholder_name": "Imaginary User"
		},
		"has_explicit_shared_members": false,
		"id": "id:a4ayc_80_OEAAAAAAAAAXw",
		"is_downloadable": true,
		"name": "Prime_Numbers.txt",
		"path_display": "/Homework/math/Prime_Numbers.txt",
		"path_lower": "/homework/math/prime_numbers.txt",
		"property_groups": [
			{
				"fields": [
					{
						"name": "Security Policy",
						"value": "Confidential"
					}
				],
				"template_id": "ptid:1a5n2i6d3OYEAAAAAAAAAYa"
			}
		],
		"rev": "a1c10ce0dd78",
		"server_modified": "2015-05-12T15:50:38Z",
		"sharing_info": {
			"modified_by": "dbid:AAH4f99T0taONIb-OurWxbNQ6ywGRopQngc",
			"parent_shared_folder_id": "84528192421",
			"read_only": true
		},
		"size": 7212
	}`

	MockUploadEndPointFailed        = "https://failed.mock.content.dropboxapi.com/2/files/upload"
	MockUploadResponseNotAuthorized = "not authorized"
)

var (
	MockUploadResponseReadCloser              io.ReadCloser
	MockUploadResponseNotAuthorizedReadCloser io.ReadCloser
)

func TestMain(m *testing.M) {
	lg, _ = logger.NewLogger()
	// http-mock setup
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", MockUploadEndPoint,
		httpmock.NewStringResponder(200, string(MockUploadResponse)),
	)
	httpmock.RegisterResponder("POST", MockUploadEndPointFailed,
		httpmock.NewStringResponder(401, "not authorized"),
	)

	MockUploadResponseReader := strings.NewReader(MockUploadResponse)
	MockUploadResponseReadCloser = ioutil.NopCloser(MockUploadResponseReader)

	MockUploadResponseNotAuthorizedReader := strings.NewReader(MockUploadResponseNotAuthorized)
	MockUploadResponseNotAuthorizedReadCloser = ioutil.NopCloser(MockUploadResponseNotAuthorizedReader)
	m.Run()
}
func Test_uploadClient_Upload(t *testing.T) {
	type fields struct {
		logger   *zap.Logger
		token    string
		mode     OverwriteMode
		endpoint string // url
	}
	type args struct {
		srcFile string
		dstdir  string
		content *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name:    "ok",
			fields:  fields{logger: lg, token: "aaaaaaa", mode: ModeAdd, endpoint: MockUploadEndPoint},
			args:    args{srcFile: "srcfile", dstdir: "/dir/", content: nil},
			wantOk:  true,
			wantErr: false,
		},
		{
			name:    "not authorized",
			fields:  fields{logger: lg, token: "aaaaaaab", mode: ModeAdd, endpoint: MockUploadEndPointFailed},
			args:    args{srcFile: "srcfile", dstdir: "/dir/", content: nil},
			wantOk:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uploadClient{
				logger:   tt.fields.logger,
				token:    tt.fields.token,
				mode:     tt.fields.mode,
				endpoint: tt.fields.endpoint,
			}
			gotOk, err := u.Upload(tt.args.srcFile, tt.args.dstdir, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("uploadClient.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("uploadClient.Upload() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_parseUploadResponse(t *testing.T) {
	type args struct {
		l *zap.Logger
		r *http.Response
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				l: lg,
				r: &http.Response{
					StatusCode: 200,
					Body:       MockUploadResponseReadCloser,
				},
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "not authorized",
			args: args{
				l: lg,
				r: &http.Response{
					StatusCode: 401,
					Body:       MockUploadResponseReadCloser,
				},
			},
			wantOk:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := parseUploadResponse(tt.args.l, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUploadResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseUploadResponse() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
