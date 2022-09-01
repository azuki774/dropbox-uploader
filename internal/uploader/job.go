package uploader

import (
	"azuki774/dropbox-uploader/internal/dropbox"
	"io/ioutil"
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
	Dryrun      bool
}

type UploadOperator struct {
	Logger          *zap.Logger
	DstDir          string // Dropbox directory
	BaseDir         string // source based directory
	OverwriteMode   dropbox.OverwriteMode
	AccessToken     string
	ContentOperator FileContentOperator
	UploadClient    dropbox.UploadClient
}

func normalizeArgs(opt *UploadOption) (err error) {
	// 	// src-dir : dir/ -> /***/dir
	// 	// dst-dir : dir/ -> dir
	opt.SrcDir, err = filepath.Abs(opt.SrcDir)
	if err != nil {
		opt.Logger.Error("failed to abs path", zap.Error(err))
		return err
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
	opt.Logger.Info("upload mode", zap.String("mode", string(o.OverwriteMode)))

	for _, fName := range fileList {
		opt.Logger.Debug("process file", zap.String("filename", fName))
		if opt.Dryrun {
			opt.Logger.Info("upload dry-run", zap.String("filename", fName))
			continue
		}

		if err = o.UploadFile(fName); err != nil {
			opt.Logger.Error("upload error", zap.Error(err))
			return err
		}
	}

	return nil
}

func NewUploadOperator(opt *UploadOption) *UploadOperator {
	Upop := UploadOperator{Logger: opt.Logger, BaseDir: filepath.Dir(opt.SrcDir), DstDir: opt.DstDir, AccessToken: opt.AccessToken}

	Upop.OverwriteMode = dropbox.ModeAdd
	if opt.Overwrite {
		Upop.OverwriteMode = dropbox.ModeOverWrite
	}

	if opt.Update {
		Upop.OverwriteMode = dropbox.ModeUpdate
	}

	Upop.ContentOperator = NewOsFileContent()
	Upop.UploadClient = dropbox.NewUploadClient(opt.Logger, opt.AccessToken, Upop.OverwriteMode)
	return &Upop
}

// UploadFile uploads srcFile(full-path) to dropbox
func (o *UploadOperator) UploadFile(abspath string) (err error) {
	content, err := o.ContentOperator.Open(abspath)
	if err != nil {
		return err
	}
	defer o.ContentOperator.Close()

	// abs path -> dropbox-relative path
	srcFile, err := filepath.Rel(o.BaseDir, abspath)
	if err != nil {
		return err
	}

	ok, err := o.UploadClient.Upload(srcFile, o.DstDir, content)

	if err != nil {
		// can not continue upload error
		o.Logger.Error("failed to upload error", zap.String("filename", abspath), zap.Error(err))
		return err
	}

	if !ok {
		// upload error, but continue
		o.Logger.Warn("failed to upload", zap.String("filename", abspath))
	}
	return nil
}

// checkSrcDir returns a file-fullpath list which should upload.
func checkSrcDir(abspath string) ([]string, error) {
	fList := []string{}
	fInfo, err := os.Stat(abspath)
	if err != nil {
		return []string{}, err
	}

	// return if file
	if !fInfo.IsDir() {
		return []string{abspath}, nil
	}

	// directory

	files, err := ioutil.ReadDir(abspath)
	if err != nil {
		return []string{}, err
	}

	for _, f := range files {
		if f.IsDir() {
			addFList, err := checkSrcDir(filepath.Join(abspath, f.Name()))
			if err != nil {
				return []string{}, err
			}

			fList = append(fList, addFList...)
			continue
		}
		fList = append(fList, filepath.Join(abspath, f.Name()))
	}

	return fList, err
}
