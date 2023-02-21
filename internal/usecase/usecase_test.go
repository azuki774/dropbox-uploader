package usecase

import (
	"reflect"
	"testing"
)

func Test_getTargetFileNames(t *testing.T) {
	type args struct {
		targetDir string
	}
	tests := []struct {
		name            string
		args            args
		wantTargetFiles []string
		wantErr         bool
	}{
		{
			name: "ok (directory)",
			args: args{
				targetDir: "../client/test/",
			},
			wantTargetFiles: []string{"../client/test/cacao.png", "../client/test/fetchNewRefreshToken_normal.json", "../client/test/uploader_normal.json"},
			wantErr:         false,
		},
		{
			name: "ok (file)",
			args: args{
				targetDir: "../client/test/cacao.png",
			},
			wantTargetFiles: []string{"../client/test/cacao.png"},
			wantErr:         false,
		},
		{
			name: "not found",
			args: args{
				targetDir: "../client/test/not_found.png",
			},
			wantTargetFiles: []string{},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTargetFiles, err := getTargetFileNames(tt.args.targetDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTargetFileNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTargetFiles, tt.wantTargetFiles) {
				t.Errorf("getTargetFileNames() = %v, want %v", gotTargetFiles, tt.wantTargetFiles)
			}
		})
	}
}
