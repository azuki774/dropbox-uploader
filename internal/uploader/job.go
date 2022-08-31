package uploader

import (
	"azuki774/dropbox-uploader/internal/dropbox"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type UploadOption struct {
	Logger      *zap.Logger
	SrcDir      string // File or Directory
	DstDir      string // Dropbox directory
	Overwrite   bool
	Update      bool
	AccessToken string
}

type UploadOperator struct {
	Logger        *zap.Logger
	DstDir        string // Dropbox directory
	OverwriteMode dropbox.OverwriteMode
	AccessToken   string
}

func Run(opt *UploadOption) (err error) {
	opt.Logger.Info("upload job start")
	fileList, err := checkSrcDir(opt.SrcDir)
	if err != nil {
		opt.Logger.Error("failed to upload file list", zap.Error(err))
		return err
	}

	for _, fName := range fileList {
		o := NewUploadOperator(opt)
		err := o.UploadFile(fName)
		if err != nil {
			opt.Logger.Error("failed to upload file", zap.Error(err))
			return err
		}
	}

	return nil
}

func NewUploadOperator(opt *UploadOption) *UploadOperator {
	Upop := UploadOperator{Logger: opt.Logger, DstDir: opt.DstDir, AccessToken: opt.AccessToken}

	Upop.OverwriteMode = dropbox.ModeAdd
	if opt.Overwrite {
		Upop.OverwriteMode = dropbox.ModeOverWrite
	}

	if opt.Update {
		Upop.OverwriteMode = dropbox.ModeUpdate
	}

	return &Upop
}

func (o *UploadOperator) UploadFile(sourceFile string) (err error) {
	req, err := dropbox.CreateUploadRequest(o.Logger, o.AccessToken, o.OverwriteMode, sourceFile, o.DstDir)
	if err != nil {
		return err
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

// checkSrcDir returns a file-fullpath list which should upload.
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
