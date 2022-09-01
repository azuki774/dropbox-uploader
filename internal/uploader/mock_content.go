package uploader

import "os"

type MockosFileContent struct {
	OpenErr error
}

func (o *MockosFileContent) Open(path string) (content *os.File, err error) {
	return nil, o.OpenErr
}

func (o *MockosFileContent) Close() (err error) {
	return nil
}
