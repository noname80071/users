package ports

import (
	"context"
	"go-users/internal/domain/models"
)

type UsersServicePort interface {
	GetUserById(ctx context.Context, id string) (*models.User, error)
	RegisterUser(ctx context.Context, username string, email string, password string) error
}

type UsersRepositoryPort interface {
	CreateUser(ctx context.Context, user models.User) error
	GetById(ctx context.Context, id string) (*models.User, error)
}

type InfrastructureService interface {
	Start(ctx context.Context) error
	GracefulShutdown(ctx context.Context) error
}
