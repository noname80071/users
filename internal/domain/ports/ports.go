package ports

import (
	"context"
	"go-users/internal/domain/models"
	"io"

	"github.com/google/uuid"
)

type UsersServicePort interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	RegisterUser(ctx context.Context, username string, email string, password string) (string, error)
}

type UsersRepositoryPort interface {
	CreateUser(ctx context.Context, user models.User) (uuid.UUID, error)
	GetByID(ctx context.Context, id string) (*models.User, error)

	UploadSkin(ctx context.Context, userID, skinURL string) error
	UploadCloak(ctx context.Context, userID, skinURL string) error
}

type InfrastructureService interface {
	Start(ctx context.Context) error
	GracefulShutdown(ctx context.Context) error
}

type FilesServicePort interface {
	UploadSkin(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	UploadCloak(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
}
