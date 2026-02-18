package users

import (
	"context"
	"go-users/internal/domain/models"
	"go-users/internal/domain/ports"

	"go-users/internal/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type usersRepository struct {
	q *query.Queries
}

func New(pool *pgxpool.Pool) ports.UsersRepositoryPort {
	return &usersRepository{q: query.New(pool)}
}

func (r *usersRepository) CreateUser(ctx context.Context, user models.User) error {
	userParams := query.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Avatar:       user.Avatar,
		Skin:         user.Skin,
		Cloak:        user.Cloak,
		RegisteredAt: user.RegisteredAt,
		IsActive:     user.IsActive,
	}
	_, err := r.q.CreateUser(ctx, userParams)

	if err != nil {
		return err
	}

	return nil
}

func (r *usersRepository) GetById(ctx context.Context, userId string) (*models.User, error) {
	userIdUUID, err := uuid.Parse(userId)

	if err != nil {
		return nil, err
	}

	user, err := r.q.GetUserById(ctx, userIdUUID)

	userModel := models.User{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Skin:     user.Skin,
		Cloak:    user.Cloak,
	}

	return &userModel, nil
}
