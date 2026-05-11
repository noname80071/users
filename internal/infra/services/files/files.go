package files

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	usersErrors "gitlab.com/_spacemc_/web/users/errors"

	avatars "gitlab.com/_spacemc_/web/users/pkg/avatars"
	minioClient "gitlab.com/_spacemc_/web/users/pkg/minio"

	"github.com/jackc/pgx/v5"
)

type UsersRepository interface {
	GetSkin(ctx context.Context, userID string) (string, error)
	GetCloak(ctx context.Context, userID string) (string, error)

	UploadSkin(ctx context.Context, userID, skinURL string) error
	DeleteSkin(ctx context.Context, userID string) error

	UploadCloak(ctx context.Context, userID, skinURL string) error
	DeleteCloak(ctx context.Context, userID string) error

	UploadAvatar(ctx context.Context, userID, avatarURL string) error
}

type FilesService struct {
	repository UsersRepository
	storage    *minioClient.Client
}

func New(repo UsersRepository, minioClient *minioClient.Client) *FilesService {
	return &FilesService{repository: repo, storage: minioClient}
}

func (s *FilesService) GetUserSkin(ctx context.Context, userID string) (string, error) {
	skinURL, err := s.repository.GetSkin(ctx, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", usersErrors.ErrUserNotFound
		}
		return "", err
	}

	return skinURL, nil
}

func (s *FilesService) GetUserCloak(ctx context.Context, userID string) (string, error) {
	cloakURL, err := s.repository.GetCloak(ctx, userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", usersErrors.ErrUserNotFound
		}
		return "", err
	}

	return cloakURL, nil
}

func (s *FilesService) UploadSkin(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error) {

	// Чтение файла в память
	imageData, err := io.ReadAll(fileReader)
	if err != nil {
		return "", err
	}

	avatarBytes, err := avatars.CropAvatar(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to crop avatar: %w", err)
	}

	avatarReader := bytes.NewReader(avatarBytes)
	skinReader := bytes.NewReader(imageData)

	avatarFilename := fmt.Sprintf("avatars/avatar_%s.png", userID)
	skinFilename := fmt.Sprintf("skins/skin_%s.png", userID)

	avatarURL, err := s.storage.Upload(ctx, "users", avatarFilename, avatarReader, avatarReader.Size())
	if err != nil {
		return "", err
	}

	skinURL, err := s.storage.Upload(ctx, "users", skinFilename, skinReader, fileSize)
	if err != nil {
		return "", err
	}

	err = s.repository.UploadAvatar(ctx, userID, avatarURL)
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

	cloakFilename := fmt.Sprintf("cloaks/cloak_%s.png", userID)

	err := s.storage.Delete(ctx, "users", cloakFilename)
	if err != nil {
		return "", err
	}

	cloakURL, err := s.storage.Upload(ctx, "cloaks", cloakFilename, fileReader, fileSize)
	if err != nil {
		return "", usersErrors.FailedToUploadFile
	}

	err = s.repository.UploadCloak(ctx, userID, cloakURL)
	if err != nil {
		return "", err
	}

	return cloakURL, nil
}

func (s *FilesService) DeleteSkin(ctx context.Context, userID string) error {
	filename := fmt.Sprintf("skins/skin_%s.png", userID)

	err := s.storage.Delete(ctx, "users", filename)

	if err != nil {
		return err
	}

	err = s.repository.DeleteSkin(ctx, userID)

	if err != nil {
		return err
	}

	return nil

}

func (s *FilesService) DeleteCloak(ctx context.Context, userID string) error {
	filename := fmt.Sprintf("cloaks/cloak_%s.png", userID)

	err := s.storage.Delete(ctx, "users", filename)

	if err != nil {
		return err
	}

	err = s.repository.DeleteCloak(ctx, userID)

	if err != nil {
		return err
	}

	return nil
}
