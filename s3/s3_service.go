package s3

import (
	"bytes"
	"context"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/common-go/storage"
)

type Config struct {
	Region          string `mapstructure:"region" json:"region,omitempty" gorm:"column:region" bson:"region,omitempty" dynamodbav:"region,omitempty" firestore:"region,omitempty"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"accessKeyID,omitempty" gorm:"column:accessKeyID" bson:"accessKeyID,omitempty" dynamodbav:"accessKeyID,omitempty" firestore:"accessKeyID,omitempty"`
	SecretAccessKey string `mapstructure:"secret_access_key" json:"secretAccessKey,omitempty" gorm:"column:secretaccesskey" bson:"secretAccessKey,omitempty" dynamodbav:"secretAccessKey,omitempty" firestore:"secretAccessKey,omitempty"`
	Token           string `mapstructure:"token" json:"token,omitempty" gorm:"column:token" bson:"token,omitempty" dynamodbav:"token,omitempty" firestore:"token,omitempty"`
}

type S3Service struct {
	session *session.Session
	Config  storage.Config
}

func NewSession(config Config) (*session.Session, error) {
	c := &aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, config.Token),
	}
	return session.NewSession(c)
}
func NewS3ServiceWithConfig(c Config, config storage.Config) (*S3Service, error) {
	session, err := NewSession(c)
	if err != nil {
		return nil, err
	}
	service := &S3Service{session: session, Config: config}
	return service, nil
}

func NewS3Service(session *session.Session, config storage.Config) *S3Service {
	service := &S3Service{session: session, Config: config}
	return service
}

func (s S3Service) Upload(ctx context.Context, objectFile storage.File) (*storage.StorageResult, error) {
	dir := objectFile.Name
	if len(s.Config.SubDirectory) > 0 {
		dir = path.Join(s.Config.SubDirectory, objectFile.Name)
	}
	uploader := s3manager.NewUploader(s.session)
	i := &s3manager.UploadInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(dir),
		Body:   bytes.NewReader(objectFile.Bytes),
	}
	if s.Config.AllUsersAreReader != nil && *s.Config.AllUsersAreReader {
		i.ACL = aws.String("public-read")
	} else if s.Config.AllAuthenticatedUsersReader != nil && *s.Config.AllAuthenticatedUsersReader {
		i.ACL = aws.String("authenticated-read")
	}
	up, err := uploader.Upload(i)
	if err != nil {
		return nil, err
	}
	return &storage.StorageResult{Status: 1, Name: objectFile.Name, MediaLink: up.Location, Link: up.Location}, nil
}

func (s S3Service) Delete(ctx context.Context, fileName string) (bool, error) {
	dir := fileName
	if len(s.Config.SubDirectory) > 0 {
		dir = path.Join(s.Config.SubDirectory, fileName)
	}
	out, err := s3.New(s.session).DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{Bucket: &s.Config.BucketName, Key: &dir})
	if err != nil {
		return false, err
	}
	return out != nil, nil
}
