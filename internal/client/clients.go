package client

import (
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
