package uploader

import (
	"azuki774/dropbox-uploader/internal/dropbox"
	"azuki774/dropbox-uploader/internal/logger"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
)

var lg *zap.Logger

func TestMain(m *testing.M) {
	lg, _ = logger.NewLogger()
	// http-mock setup
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", dropbox.UploadEndPoint,
		httpmock.NewStringResponder(200, "mocked"),
	)

	m.Run()
}
func Test_checkSrcDir(t *testing.T) {
	type args struct {
		path string // use relative-path for test
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "file select",
			args:      args{path: "test/testfile1.txt"},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "directory select",
			args:      args{path: "test/"},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:    "not found",
			args:    args{path: "notfound/"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkSrcDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkSrcDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("checkSrcDir() = %v, want %v", got, tt.want)
			// }
			fmt.Printf("got = %v\n", got)
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("checkSrcDir() count = %v, want %v", got, tt.wantCount)
			}
		})
	}
}

func TestUploadOperator_UploadFile(t *testing.T) {
	type fields struct {
		Logger          *zap.Logger
		DstDir          string
		BaseDir         string
		OverwriteMode   dropbox.OverwriteMode
		AccessToken     string
		ContentOperator FileContentOperator
		UploadClient    dropbox.UploadClient
	}
	type args struct {
		abspath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Logger:          lg,
				DstDir:          "/dropbox/",
				BaseDir:         "/var/tmp/",
				OverwriteMode:   dropbox.ModeAdd,
				AccessToken:     "aaaaaaaaaa",
				ContentOperator: &MockosFileContent{},
				UploadClient:    &dropbox.MockUploadClient{OK: true, Err: nil},
			},
			args:    args{abspath: "/var/tmp/testfile"},
			wantErr: false,
		},
		{
			name: "upload failed no error",
			fields: fields{
				Logger:          lg,
				DstDir:          "/dropbox/",
				BaseDir:         "/var/tmp/",
				OverwriteMode:   dropbox.ModeAdd,
				AccessToken:     "aaaaaaaaaa",
				ContentOperator: &MockosFileContent{},
				UploadClient:    &dropbox.MockUploadClient{OK: false, Err: nil},
			},
			args:    args{abspath: "/var/tmp/testfile"},
			wantErr: false,
		},
		{
			name: "upload error",
			fields: fields{
				Logger:          lg,
				DstDir:          "/dropbox/",
				BaseDir:         "/var/tmp/",
				OverwriteMode:   dropbox.ModeAdd,
				AccessToken:     "aaaaaaaaaa",
				ContentOperator: &MockosFileContent{},
				UploadClient:    &dropbox.MockUploadClient{OK: false, Err: errors.New("error")},
			},
			args:    args{abspath: "/var/tmp/testfile"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &UploadOperator{
				Logger:          tt.fields.Logger,
				DstDir:          tt.fields.DstDir,
				BaseDir:         tt.fields.BaseDir,
				OverwriteMode:   tt.fields.OverwriteMode,
				AccessToken:     tt.fields.AccessToken,
				ContentOperator: tt.fields.ContentOperator,
				UploadClient:    tt.fields.UploadClient,
			}
			if err := o.UploadFile(tt.args.abspath); (err != nil) != tt.wantErr {
				t.Errorf("UploadOperator.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
