package dropbox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

type OverwriteMode string

const (
	ModeAdd       = OverwriteMode("add")
	ModeUpdate    = OverwriteMode("update")
	ModeOverWrite = OverwriteMode("overwrite")

	UploadEndPoint = "https://content.dropboxapi.com/2/files/upload"
)

type UploadRequest struct {
	Autorename     bool   `json:"autorename"`
	Mode           string `json:"mode"`
	Mute           bool   `json:"mute"`
	Path           string `json:"path"`
	StrictConflict bool   `json:"strict_conflict"`
}

type UploadResponseOK struct {
	ClientModified time.Time `json:"client_modified"`
	ContentHash    string    `json:"content_hash"`
	FileLockInfo   struct {
		Created        time.Time `json:"created"`
		IsLockholder   bool      `json:"is_lockholder"`
		LockholderName string    `json:"lockholder_name"`
	} `json:"file_lock_info"`
	HasExplicitSharedMembers bool   `json:"has_explicit_shared_members"`
	ID                       string `json:"id"`
	IsDownloadable           bool   `json:"is_downloadable"`
	Name                     string `json:"name"`
	PathDisplay              string `json:"path_display"`
	PathLower                string `json:"path_lower"`
	PropertyGroups           []struct {
		Fields []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"fields"`
		TemplateID string `json:"template_id"`
	} `json:"property_groups"`
	Rev            string    `json:"rev"`
	ServerModified time.Time `json:"server_modified"`
	SharingInfo    struct {
		ModifiedBy           string `json:"modified_by"`
		ParentSharedFolderID string `json:"parent_shared_folder_id"`
		ReadOnly             bool   `json:"read_only"`
	} `json:"sharing_info"`
	Size int `json:"size"`
}

// CreateUploadRequest creates http.Request for uploading dropbox.
// dst-dir must not include '/' at the end of URL.
func CreateUploadRequest(l *zap.Logger, content *os.File, token string, mode OverwriteMode, srcFile string, dstdir string) (*http.Request, error) {
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

func ParseUploadResponse(l *zap.Logger, r *http.Response) (ok bool, err error) {
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

	l.Info("upload sucess", zap.String("path_display", res.PathDisplay), zap.Time("server_modified", res.ServerModified))

	return true, nil
}
