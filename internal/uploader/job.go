package uploader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type UploadOption struct {
	Logger      *zap.Logger
	SrcDir      string // File or Directory
	DstDir      string // Dropbox directory
	OverWrite   bool
	AccessToken string
}

func Run(opt *UploadOption) (err error) {
	opt.Logger.Info("upload job start")
	_, err = checkSrcDir(opt.SrcDir)
	if err != nil {
		opt.Logger.Error("failed to upload file list", zap.Error(err))
		return err
	}
	return nil
}

// checkSrcDir returns a file list which should upload.
func checkSrcDir(path string) ([]string, error) {
	fInfo, err := os.Stat(path)
	if err != nil {
		return []string{}, err
	}

	if !fInfo.IsDir() {
		// file
		fName, err := filepath.Abs(path)
		if err != nil {
			return []string{}, err
		}
		return []string{fName}, nil
	}

	// directory
	fList := []string{}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	for _, f := range files {
		fName, err := filepath.Abs(f.Name())
		if err != nil {
			return []string{}, err
		}
		fList = append(fList, fName)
	}

	return fList, err
}
