package users

import (
	"context"
	"errors"
	"time"

	usersErrors "gitlab.com/_spacemc_/web/users/errors"
	"gitlab.com/_spacemc_/web/users/internal/domain/models"
	validate "gitlab.com/_spacemc_/web/users/pkg/validate"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type UsersRepository interface {
	CreateUser(ctx context.Context, user models.User) (uuid.UUID, error)

	GetByID(ctx context.Context, userID string) (*models.User, error)
	GetByEmail(ctx context.Context, userEmail string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	GetSkin(ctx context.Context, userID string) (string, error)
	GetCloak(ctx context.Context, userID string) (string, error)

	UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error)
}

type UsersService struct {
	repository UsersRepository
}

func New(repo UsersRepository) *UsersService {
	return &UsersService{repository: repo}
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
		return nil, err
	}

	return user, nil

}

func (s *UsersService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.repository.GetByUsername(ctx, username)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usersErrors.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UsersService) RegisterUser(ctx context.Context, username string, email string, password string) (string, error) {

	passwordHash, err := HashPassword(password)

	if err != nil {
		return "", err
	}

	// Валидация
	if !validate.IsValidEmail(email) {
		return "", usersErrors.ErrEmailInvalid
	}

	if !validate.IsValidPassword(password) {
		return "", usersErrors.ErrPasswordWeak
	}

	if !validate.IsValidUsername(username) {
		return "", usersErrors.ErrUsernameInvalid
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
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return "", usersErrors.ErrUserAlreadyExists
		}
		return "", err
	}

	return id.String(), nil
}

func (s *UsersService) UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error) {
	id, err := s.repository.UpdateUserStatus(ctx, userID, active)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", usersErrors.ErrUserNotFound
		}

		return "", err
	}

	return id, nil
}
