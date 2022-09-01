package uploader

import "os"

type FileContent struct{}

type FileContentOperator interface {
	Open(path string) (content *os.File, err error)
	Close() (err error)
}

type osFileContent struct {
	content *os.File
}

func NewOsFileContent() *osFileContent {
	return &osFileContent{}
}

func (o *osFileContent) Open(path string) (content *os.File, err error) {
	o.content, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	return o.content, nil
}

func (o *osFileContent) Close() (err error) {
	return o.content.Close()
}
