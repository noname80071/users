package minio

import (
	"context"
	"fmt"
	"go-users/internal/domain/ports"
	usersRepo "go-users/internal/infra/repositories"
	"io"
	"path/filepath"
	"time"

	minioClient "go-users/pkg/minio"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FilesService struct {
	repository ports.UsersRepositoryPort
	storage    *minioClient.Client
}

func New(pool *pgxpool.Pool, minioClient *minioClient.Client) ports.FilesServicePort {
	return &FilesService{repository: usersRepo.New(pool), storage: minioClient}
}

func (s *FilesService) UploadSkin(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error) {
	// Проверка что пользователь существует!!!

	ext := filepath.Ext(filename)
	skinFilename := fmt.Sprintf("%s/skin%d%s", userID, time.Now().Unix(), ext)

	contentType := "application/octet-stream"

	skinURL, err := s.storage.Upload(ctx, "skins", skinFilename, fileReader, fileSize, contentType)
	if err != nil {
		return "", err
	}

	err = s.repository.UploadSkin(ctx, userID, skinURL)
	if err != nil {
		return "", err
	}

	return skinURL, nil
}

func (s *FilesService) UploadCloak(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error) {
	// Проверка что пользователь существует!!!

	ext := filepath.Ext(filename)
	cloakFilename := fmt.Sprintf("%s/cloak%d%s", userID, time.Now().Unix(), ext)

	contentType := "application/octet-stream"

	cloakURL, err := s.storage.Upload(ctx, "cloaks", cloakFilename, fileReader, fileSize, contentType)
	if err != nil {
		return "", err
	}

	err = s.repository.UploadCloak(ctx, userID, cloakURL)
	if err != nil {
		return "", err
	}

	return cloakURL, nil
}
