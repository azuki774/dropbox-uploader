package uploader

import (
	"azuki774/dropbox-uploader/internal/dropbox"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func normalizeArgs(opt *UploadOption) (err error) {
	// src-dir : dir/ -> /***/dir
	// dst-dir : dir/ -> dir
	opt.SrcDir, err = filepath.Abs(opt.SrcDir)
	if err != nil {
		opt.Logger.Error("failed to abs path", zap.Error(err))
		return err
	}

	if strings.HasSuffix(opt.DstDir, "/") {
		l := len(opt.DstDir)
		opt.DstDir = opt.DstDir[:(l - 1)]
	}

	return nil
}

func Run(opt *UploadOption) (err error) {
	opt.Logger.Info("upload job start")
	if err := normalizeArgs(opt); err != nil {
		return err
	}

	fileList, err := checkSrcDir(opt.SrcDir)
	if err != nil {
		opt.Logger.Error("failed to upload file list", zap.Error(err))
		return err
	}

	o := NewUploadOperator(opt)

	for _, fName := range fileList {
		opt.Logger.Debug("process file", zap.String("filename", fName))
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

// UploadFile uploads srcFile(full-path) to dropbox
func (o *UploadOperator) UploadFile(srcFile string) (err error) {
	content, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer content.Close()

	req, err := dropbox.CreateUploadRequest(o.Logger, content, o.AccessToken, o.OverwriteMode, srcFile, o.DstDir)
	if err != nil {
		return err
	}

	client := new(http.Client)
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// checkSrcDir returns a file-fullpath list which should upload.
func checkSrcDir(abspath string) ([]string, error) {
	fInfo, err := os.Stat(abspath)
	if err != nil {
		return []string{}, err
	}

	if !fInfo.IsDir() {
		// file
		fName, err := filepath.Abs(abspath)
		if err != nil {
			return []string{}, err
		}
		return []string{fName}, nil
	}

	// directory
	fList := []string{}
	files, err := ioutil.ReadDir(abspath)
	if err != nil {
		return []string{}, err
	}

	for _, f := range files {
		fName := abspath + "/" + f.Name()
		fList = append(fList, fName)
	}

	return fList, err
}
