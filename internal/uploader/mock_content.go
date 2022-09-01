package uploader

import "os"

type MockosFileContent struct{}

func NewMockOsFileContent() *MockosFileContent {
	return &MockosFileContent{}
}

func (o *MockosFileContent) Open(path string) (content *os.File, err error) {
	return nil, nil
}

func (o *MockosFileContent) Close() (err error) {
	return nil
}
