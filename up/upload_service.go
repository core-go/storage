package upload

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
)

type UploadManager interface {
	Upload(ctx context.Context, id string, data Request) (*Attachment, error)
	Delete(ctx context.Context, id string, url string) (int64, error)
}

type StorageService interface {
	Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, id string) (bool, error)
}

type UploadService struct {
	repository       StorageRepository
	Service          StorageService
	Provider         string
	GeneralDirectory string
	Directory        string
	KeyFile          string
	SizesImage       []int
	SizesCover       []int
}

func NewUploadService(
	repository StorageRepository,
	service StorageService, provider string, generalDirectory string,
	keyFile string, directory string) UploadManager {

	return &UploadService{Service: service, Provider: provider, GeneralDirectory: generalDirectory,
		KeyFile: keyFile, Directory: directory, repository: repository}
}

func (u *UploadService) Upload(ctx context.Context, id string, req Request) (*Attachment, error) {
	e, err := u.repository.Load(ctx, id)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, nil
	}

	url, err := u.uploadFileOnServer(ctx, req.Filename, req.Type, req.Data)
	if err != nil {
		return nil, err
	}
	attachment := Attachment{
		OriginalFileName: req.OriginalFileName,
		FileName:         req.Filename,
		Type:             strings.Split(req.Type, "/")[0],
		Url:              url,
	}
	e.Attachment = attachment
	rows, err := u.repository.Update(ctx, id, e.Attachment)
	if err != nil {
		return nil, err
	}
	if rows > 0 {
		return &attachment, nil
	}
	return nil, nil
}

func (u *UploadService) Delete(ctx context.Context, id string, url string) (int64, error) {
	rs, err := u.repository.Load(ctx, id)
	if err != nil {
		return 0, err
	}
	if rs == nil {
		return -1, nil
	}
	attachments := rs.Attachment
	exist := false
	if rs.Attachment.Url == url {
		_, err2 := u.deleteFile(url, ctx)
		if err2 != nil {
			return 0, err2
		}
		exist = true
	}
	if exist == false {
		return -1, errors.New("no exist file " + url)
	}
	_, err2 := u.repository.Update(ctx, id, attachments)
	if err2 != nil {
		return 0, err2
	}
	return 1, nil
}

func (u *UploadService) uploadFileOnServer(ctx context.Context, fileName string, contentType string, data []byte) (rs string, errorRespone error) {
	directory := u.Directory
	rs, err2 := u.Service.Upload(ctx, directory, fileName, data, contentType)
	if err2 != nil {
		errorRespone = err2
		return
	}
	return
}

func (u *UploadService) deleteFile(url string, ctx context.Context) (bool, error) {
	arrOrigin := strings.Split(url, "/")
	delOriginUrl := arrOrigin[len(arrOrigin)-2] + "/" + arrOrigin[len(arrOrigin)-1]
	rs, err := u.Service.Delete(ctx, delOriginUrl)
	return rs, err
}

func getExt(file string) string {
	ext := filepath.Ext(file)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
		return ext
	}
	return ext
}

func removeExt(file string) string {
	return file[:len(file)-len(filepath.Ext(file))]
}
