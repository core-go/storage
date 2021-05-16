package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
)

func NewClientWithCredentialsFile(ctx context.Context, credentialsFile string) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile)) //"resource/key.json"
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewClient(ctx context.Context, credentials []byte) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, err
	}
	return client, nil
}
func NewClientWithOptions(ctx context.Context, opt...option.ClientOption) (*storage.Client, error) {
	client, err := storage.NewClient(ctx, opt...) //"resource/key.json"
	if err != nil {
		return nil, err
	}
	return client, nil
}
