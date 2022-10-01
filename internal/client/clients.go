package client

import (
	"io"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type Client struct {
	dbxFiles *files.Client
}

func (c *Client) RenewClient(newToken string) {
	config := dropbox.Config{
		Token:    newToken,
		LogLevel: dropbox.LogInfo, // if needed, set the desired logging level. Default is off
	}
	dbxFilesImpl := files.New(config)
	c.dbxFiles = &dbxFilesImpl
}

func (c *Client) UploadFile(path string, content io.Reader) (err error) {
	d := *(c.dbxFiles)
	arg := files.NewUploadArg(path)
	_, err = d.Upload(arg, content)
	if err != nil {
		return err
	}
	return nil
}
