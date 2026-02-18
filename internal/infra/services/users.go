package users

import (
	"context"
	"fmt"
	"go-users/internal/domain/models"
	"go-users/internal/domain/ports"
	usersRepo "go-users/internal/infra/repositories"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	repository ports.UsersRepositoryPort
}

func New(pool *pgxpool.Pool) ports.UsersServicePort {
	return &usersService{repository: usersRepo.New(pool)}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *usersService) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	user, err := s.repository.GetById(ctx, userId)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersService) RegisterUser(ctx context.Context, username string, email string, password string) error {

	passwordHash, err := HashPassword(password)

	if err != nil {
		return err
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Avatar:       "",
		Skin:         pgtype.Text{String: "", Valid: true},
		Cloak:        pgtype.Text{String: "", Valid: true},
		RegisteredAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		IsActive:     true,
	}

	fmt.Println(user)

	errS := s.repository.CreateUser(ctx, user)

	if errS != nil {
		return err
	}

	return nil
}
