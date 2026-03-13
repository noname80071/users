package ports

import (
	"context"
	"io"

	"gitlab.com/_spacemc_/web/users/internal/domain/models"

	"github.com/google/uuid"
)

type UsersServicePort interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	RegisterUser(ctx context.Context, username string, email string, password string) (string, error)
	UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error)
}

type UsersRepositoryPort interface {
	CreateUser(ctx context.Context, user models.User) (uuid.UUID, error)

	GetByID(ctx context.Context, userID string) (*models.User, error)
	GetByEmail(ctx context.Context, userEmail string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	GetSkin(ctx context.Context, userID string) (string, error)
	GetCloak(ctx context.Context, userID string) (string, error)

	UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error)

	UploadSkin(ctx context.Context, userID, skinURL string) error
	DeleteSkin(ctx context.Context, userID string) error

	UploadCloak(ctx context.Context, userID, skinURL string) error
	DeleteCloak(ctx context.Context, userID string) error
}

type InfrastructureService interface {
	Start(ctx context.Context) error
	GracefulShutdown(ctx context.Context) error
}

type FilesServicePort interface {
	GetUserSkin(ctx context.Context, userID string) (string, error)
	GetUserCloak(ctx context.Context, userID string) (string, error)

	UploadSkin(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	DeleteSkin(ctx context.Context, userID string) error

	UploadCloak(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	DeleteCloak(ctx context.Context, userID string) error
}
