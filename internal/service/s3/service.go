package s3

import (
	"context"
	"fmt"
	"io"
	"time"


	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	aconfig "asona/config"
)

// Service defines S3 operations.
type Service interface {
	Upload(ctx context.Context, key string, body io.Reader) (string, error)
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}

type service struct {
	client     *s3.Client
	bucketName string
	region     string
}

func New() Service {
	cfg := aconfig.GetConfig()

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSS3Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSS3AccessKey,
			cfg.AWSS3SecretKey,
			"")),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load S3 configuration: %v", err))
	}

	client := s3.NewFromConfig(awsCfg)

	return &service{
		client:     client,
		bucketName: cfg.AWSS3BucketName,
		region:     cfg.AWSS3Region,
	}
}

func (s *service) Upload(ctx context.Context, key string, body io.Reader) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
		Body:   body,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)
	return url, nil
}


func (s *service) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *service) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	return err
}

func (s *service) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return request.URL, nil
}
