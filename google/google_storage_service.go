package google

import (
	"cloud.google.com/go/storage"
	"context"
	st "github.com/core-go/storage"
	"path"
)

const storageUrl = "https://storage.googleapis.com"

type GoogleStorageService struct {
	Client *storage.Client
	Config st.Config
	Bucket *storage.BucketHandle
}
func NewGoogleStorageServiceWithCredentials(ctx context.Context, credentials []byte, config st.Config) (*GoogleStorageService, error) {
	client, err := NewClient(ctx, credentials)
	if err != nil {
		return nil, err
	}
	gs := NewGoogleStorageService(client, config)
	return gs, nil
}
func NewGoogleStorageServiceWithCredentialsFile(ctx context.Context, credentialsFile string, config st.Config) (*GoogleStorageService, error) {
	client, err := NewClientWithCredentialsFile(ctx, credentialsFile)
	if err != nil {
		return nil, err
	}
	gs := NewGoogleStorageService(client, config)
	return gs, nil
}

func NewGoogleStorageService(client *storage.Client, config st.Config) *GoogleStorageService {
	return &GoogleStorageService{client,
		config,
		client.Bucket(config.Bucket)}
}

func (s GoogleStorageService) Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error) {
	dir := filename
	if len(directory) > 0 {
		dir = path.Join(directory, filename)
	}
	object := s.Bucket.Object(dir)
	wc := object.NewWriter(ctx)
	wc.ContentType = contentType // "image/png"

	if len(contentType) > 0 {
		wc.ContentType = contentType
	}
	_, err := wc.Write(data)
	if err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	if s.Config.Public != nil && *s.Config.Public == true {
		if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return "", err
		}
	} else {
		if s.Config.Private == nil || *s.Config.Private == false {
			if err := object.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
				return "", err
			}
		} else if s.Config.AllAuthenticatedUsersWriter != nil && *s.Config.AllAuthenticatedUsersWriter {
			if err := object.ACL().Set(ctx, storage.AllAuthenticatedUsers, storage.RoleWriter); err != nil {
				return "", err
			}
		}
	}
	// attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}
	link :=  getLinkPublic(s.Config.Bucket, dir)
	return link, nil
}

func getLinkPublic(bucketName string, remoteFile string) string {
	return path.Join(storageUrl, bucketName, remoteFile)
}

func (s GoogleStorageService) Delete(ctx context.Context, directory string, fileName string) (bool, error) {
	dir := fileName
	if len(directory) > 0 {
		dir = path.Join(directory, fileName)
	}
	if err := s.Bucket.Object(dir).Delete(ctx); err != nil {
		return false, err
	}
	return true, nil
}
