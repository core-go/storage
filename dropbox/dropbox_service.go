package dropbox

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

type DropboxService struct {
	DropboxToken	string
}

func NewDropboxService(token string) (*DropboxService, error) {
	return &DropboxService{DropboxToken: token}, nil
}

func (d DropboxService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	file := bytes.NewReader(data)

	// create new client to access dropbox cloud with token generated in dropbox console
	config := dropbox.Config{
		Token: d.DropboxToken,
	}
	client := files.New(config)

	// create new upload info
	filepath := fmt.Sprintf("/%s/%s", directory, filename)
	arg := files.NewCommitInfo(filepath)

	// upload file
	_, err := client.Upload(arg, file)
	if err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("uploaded files '%s' to dropbox successfully, file location: https://www.dropbox.com/home/Apps/golang-upload/%s", filename, directory)
	return msg, err
}

func (d DropboxService)  Delete(ctx context.Context, directory string, fileName string) (bool, error) {
	return false, nil
}