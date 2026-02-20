package users

import (
	"context"
	serviceErrors "go-users/errors"
	"go-users/internal/domain/models"
	"go-users/internal/domain/ports"
	usersRepo "go-users/internal/infra/repositories"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	repository ports.UsersRepositoryPort
}

func New(pool *pgxpool.Pool) ports.UsersServicePort {
	return &UsersService{repository: usersRepo.New(pool)}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *UsersService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repository.GetByID(ctx, userID)

	if err != nil {
		return nil, serviceErrors.ErrUserNotFound
	}

	return user, nil
}

func (s *UsersService) RegisterUser(ctx context.Context, username string, email string, password string) (string, error) {

	passwordHash, err := HashPassword(password)

	if err != nil {
		return "", err
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Avatar:       "",
		Skin:         "",
		Cloak:        "",
		RegisteredAt: time.Now(),
		IsActive:     true,
	}

	id, err := s.repository.CreateUser(ctx, user)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}
