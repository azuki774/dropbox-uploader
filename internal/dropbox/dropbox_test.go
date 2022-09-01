package dropbox

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func Test_uploadClient_Upload(t *testing.T) {
	type fields struct {
		logger *zap.Logger
		token  string
		mode   OverwriteMode
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &uploadClient{
				logger: tt.fields.logger,
				token:  tt.fields.token,
				mode:   tt.fields.mode,
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
