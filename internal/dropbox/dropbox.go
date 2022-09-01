package dropbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type OverwriteMode string

const (
	ModeAdd       = OverwriteMode("add")
	ModeUpdate    = OverwriteMode("update")
	ModeOverWrite = OverwriteMode("overwrite")

	UploadEndPoint = "https://content.dropboxapi.com/2/files/upload"
)

type UploadClient interface {
	Upload(srcFile string, dstdir string, content *os.File) (ok bool, err error)
}

type uploadClient struct {
	logger *zap.Logger
	token  string
	mode   OverwriteMode
}

func NewUploadClient(l *zap.Logger, token string, mode OverwriteMode) *uploadClient {
	return &uploadClient{
		logger: l,
		token:  token,
		mode:   mode,
	}
}

// Upload returns ok, err
// ok .. upload status
// err .. can continue uploading
func (u *uploadClient) Upload(srcFile string, dstdir string, content *os.File) (ok bool, err error) {
	req, err := createUploadRequest(u.logger, content, u.token, u.mode, srcFile, dstdir)
	if err != nil {
		return false, err
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	ok, err = parseUploadResponse(u.logger, res)
	if err != nil {
		u.logger.Error("failed to parse response", zap.Error(err))
		return false, err
	}

	return ok, nil
}

// CreateUploadRequest creates http.Request for uploading dropbox.
// dst-dir must not include '/' at the end of URL.
func createUploadRequest(l *zap.Logger, content *os.File, token string, mode OverwriteMode, srcFile string, dstdir string) (*http.Request, error) {
	apiArgs := UploadRequest{
		Autorename:     false,
		Mode:           string(mode),
		Mute:           false,
		Path:           filepath.Join(dstdir, srcFile),
		StrictConflict: false,
	}
	apiArgsBytes, err := json.Marshal(apiArgs)
	if err != nil {
		return nil, err
	}
	apiArgsString := string(apiArgsBytes)

	req, _ := http.NewRequest("POST", UploadEndPoint, content)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", apiArgsString)

	return req, nil
}

func parseUploadResponse(l *zap.Logger, r *http.Response) (ok bool, err error) {
	if r.StatusCode != 200 {
		l.Info("upload failed", zap.Int("status_code", r.StatusCode))
		return false, nil
	}

	// 200 ok
	res := &UploadResponseOK{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return false, err
	}

	l.Info("upload success", zap.String("path_display", res.PathDisplay), zap.Time("server_modified", res.ServerModified))

	return true, nil
}
