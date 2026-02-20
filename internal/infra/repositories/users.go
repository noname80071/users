package users

import (
	"context"
	"go-users/internal/domain/models"
	"go-users/internal/domain/ports"

	"go-users/internal/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	q *query.Queries
}

func New(pool *pgxpool.Pool) ports.UsersRepositoryPort {
	return &UsersRepository{q: query.New(pool)}
}

func (r *UsersRepository) CreateUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	userParams := query.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Avatar:       user.Avatar,
		Skin:         pgtype.Text{String: user.Skin, Valid: true},
		Cloak:        pgtype.Text{String: user.Cloak, Valid: true},
		RegisteredAt: pgtype.Timestamptz{Time: user.RegisteredAt, Valid: true},
		IsActive:     user.IsActive,
	}
	id, err := r.q.CreateUser(ctx, userParams)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *UsersRepository) GetByID(ctx context.Context, userId string) (*models.User, error) {
	userIDUUID, err := uuid.Parse(userId)

	if err != nil {
		return nil, err
	}

	user, err := r.q.GetUserByID(ctx, userIDUUID)

	if err != nil {
		return nil, err
	}

	userModel := models.User{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Skin:     user.Skin.String,
		Cloak:    user.Cloak.String,
	}

	return &userModel, nil
}

func (r *UsersRepository) UploadSkin(ctx context.Context, userID, skinURL string) error {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	params := query.UpdateUserSkinParams{
		Skin: pgtype.Text{String: skinURL, Valid: true},
		ID:   userIDUUID,
	}

	err = r.q.UpdateUserSkin(ctx, params)

	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) UploadCloak(ctx context.Context, userID, cloakURL string) error {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	params := query.UpdateUserCloakParams{
		Cloak: pgtype.Text{String: cloakURL, Valid: true},
		ID:    userIDUUID,
	}

	err = r.q.UpdateUserCloak(ctx, params)

	if err != nil {
		return err
	}

	return nil
}
