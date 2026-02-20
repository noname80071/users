package minio

import (
	"context"
	"fmt"
	"go-users/config"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StoragePort interface {
	Upload(ctx context.Context, bucket string, fileName string, fileReader io.Reader, fileSize int64, contentType string) (string, error)
}

type Client struct {
	Client *minio.Client
}

func NewClient(ctx context.Context, cfg *config.MinioConfig) (*Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretAccessKey, ""),
		Secure: cfg.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	log.Printf("Клиент создан")

	return &Client{Client: minioClient}, nil
}

func (c *Client) Upload(ctx context.Context, bucket string, fileName string, fileReader io.Reader, fileSize int64, contentType string) (string, error) {

	exists, err := c.Client.BucketExists(ctx, bucket)

	if err != nil || !exists {
		err := c.Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Errorf("failed to create bucket")
			return "", err
		}
	}

	info, err := c.Client.PutObject(ctx, bucket, fileName, fileReader, fileSize, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"Name": fileName,
		},
		ContentType: "application/octet-stream",
	})
	if err != nil {
		fmt.Errorf("failed to upload file: %w", err)
		return "", err
	}

	url := fmt.Sprintf("/%s/%s", bucket, info.Key)

	return url, nil

}
