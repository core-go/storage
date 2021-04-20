package storage

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, credentialsFile string) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile)) //"resource/key.json"
	if err != nil {
		return nil, err
	}
	return client, nil
}
