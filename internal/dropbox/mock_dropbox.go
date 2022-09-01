package dropbox

import "os"

type MockUploadClient struct {
	OK  bool
	Err error
}

func (m *MockUploadClient) Upload(srcFile string, dstdir string, content *os.File) (ok bool, err error) {
	return m.OK, m.Err
}
