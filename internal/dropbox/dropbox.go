package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type OverwriteMode string

const (
	ModeAdd       = OverwriteMode("add")
	ModeUpdate    = OverwriteMode("update")
	ModeOverWrite = OverwriteMode("overwrite")

	// UploadEndPoint = "https://content.dropboxapi.com/2/files/upload"
	UploadEndPoint = "https://test.azk.home"
)

type UploadRequest struct {
	Autorename     bool   `json:"autorename"`
	Mode           string `json:"mode"`
	Mute           bool   `json:"mute"`
	Path           string `json:"path"`
	StrictConflict bool   `json:"strict_conflict"`
}

type UploadResponse struct {
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

func CreateUploadRequest(l *zap.Logger, token string, mode OverwriteMode, srcfile string, dstdir string) (*http.Request, error) {
	bodyBuf, err := createUploadBody(srcfile)
	if err != nil {
		return nil, err
	}

	apiArgs := UploadRequest{
		Autorename:     false,
		Mode:           string(mode),
		Mute:           false,
		Path:           dstdir,
		StrictConflict: false,
	}
	apiArgsBytes, err := json.Marshal(apiArgs)
	if err != nil {
		return nil, err
	}
	apiArgsString := string(apiArgsBytes)

	req, _ := http.NewRequest("POST", UploadEndPoint, bodyBuf)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Dropbox-API-Arg", apiArgsString)

	return req, nil
}

func createUploadBody(srcfile string) (bb *bytes.Buffer, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//キーとなる操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", srcfile)
	if err != nil {
		return nil, err
	}

	fh, err := os.Open(srcfile)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}

	bodyWriter.Close()
	return bodyBuf, nil
}
