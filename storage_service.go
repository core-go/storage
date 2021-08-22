package storage

import "context"

type StorageService interface {
	Upload(ctx context.Context, directory string, filename string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, directory string, fileName string) (bool, error)
}
