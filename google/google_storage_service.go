package storage

import (
	"context"
	"path"

	"cloud.google.com/go/storage"
	. "github.com/common-go/storage"
)

const storageUrl = "https://storage.googleapis.com"

type GoogleStorageService struct {
	Client *storage.Client
	Config Config
	Bucket *storage.BucketHandle
}

func NewGoogleStorageServiceWithCredentialsFile(ctx context.Context, credentialsFile string, config Config) (*GoogleStorageService, error) {
	client, err := NewClient(ctx, credentialsFile)
	if err != nil {
		return nil, err
	}
	gs := NewGoogleStorageService(client, config)
	return gs, nil
}

func NewGoogleStorageService(client *storage.Client, config Config) *GoogleStorageService {
	return &GoogleStorageService{client,
		config,
		client.Bucket(config.BucketName)}
}

func (s GoogleStorageService) Upload(ctx context.Context, file File) (*StorageResult, error) {
	dir := file.Name
	if len(s.Config.SubDirectory) > 0 {
		dir = path.Join(s.Config.SubDirectory, file.Name)
	}
	object := s.Bucket.Object(dir)
	wc := object.NewWriter(ctx)
	wc.ContentType = file.ContentType // "image/png"

	if len(file.ContentType) > 0 {
		wc.ContentType = file.ContentType
	}
	n, err := wc.Write(file.Bytes)
	if err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	if s.Config.AllUsersAreReader != nil && *s.Config.AllUsersAreReader {
		if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return nil, err
		}
	}
	if s.Config.AllAuthenticatedUsersReader != nil && *s.Config.AllAuthenticatedUsersReader {
		if err := object.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
			return nil, err
		}
	}
	if s.Config.AllAuthenticatedUsersWriter != nil && *s.Config.AllAuthenticatedUsersWriter {
		if err := object.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleWriter); err != nil {
			return nil, err
		}
	}
	attrs, err := object.Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return &StorageResult{Status: int64(n), Name: file.Name, MediaLink: attrs.MediaLink, Link: getLinkPublic(s.Config.BucketName, dir)}, nil
}

func getLinkPublic(bucketName string, remoteFile string) string {
	return path.Join(storageUrl, bucketName, remoteFile)
}

func (s GoogleStorageService) Delete(ctx context.Context, fileName string) (bool, error) {
	dir := fileName
	if len(s.Config.SubDirectory) > 0 {
		dir = path.Join(s.Config.SubDirectory, fileName)
	}
	if err := s.Bucket.Object(dir).Delete(ctx); err != nil {
		return false, err
	}
	return true, nil
}
