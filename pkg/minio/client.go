package minio

import (
	"context"
	"fmt"
	"io"
	"log"

	"gitlab.com/_spacemc_/web/users/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StoragePort interface {
	Upload(ctx context.Context, bucket string, fileName string, fileReader io.Reader, fileSize int64) (string, error)
}

type Client struct {
	Client *minio.Client
	Config *config.MinioConfig
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

	return &Client{Client: minioClient, Config: cfg}, nil
}

func (c *Client) Upload(ctx context.Context, bucket string, fileName string, fileReader io.Reader, fileSize int64) (string, error) {
	info, err := c.Client.PutObject(ctx, bucket, fileName, fileReader, fileSize, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"Name": fileName,
		},
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("%s/browser/%s/%s", c.Config.ClientEndpoint, bucket, info.Key)

	return url, nil

}

func (c *Client) Delete(ctx context.Context, bucket string, fileName string) error {
	err := c.Client.RemoveObject(ctx, bucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file %s from bucket %s: %w", fileName, bucket, err)
	}

	return nil
}
