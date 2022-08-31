package uploader

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_checkSrcDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "file select",
			args:      args{path: "test/testfile1.file"},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "directory select",
			args:      args{path: "test/"},
			wantCount: 2,
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
