package users

import (
	"context"
	"errors"
	"time"

	usersErrors "gitlab.com/_spacemc_/web/users/errors"
	"gitlab.com/_spacemc_/web/users/internal/domain/models"
	"gitlab.com/_spacemc_/web/users/internal/domain/ports"
	usersRepo "gitlab.com/_spacemc_/web/users/internal/infra/repositories"

	"github.com/jackc/pgx/v5"
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersErrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UsersService) GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := s.repository.GetByEmail(ctx, userEmail)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersErrors.ErrUserNotFound
		}
	}

	return user, nil

}

func (s *UsersService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repository.GetByUsername(ctx, username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersErrors.ErrUserNotFound
		}
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

func (s *UsersService) UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error) {
	id, err := s.repository.UpdateUserStatus(ctx, userID, active)

	if err != nil {
		return "", usersErrors.ErrUserNotFound
	}

	return id, nil
}
