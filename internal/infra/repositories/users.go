package users

import (
	"context"

	"gitlab.com/_spacemc_/web/users/internal/domain/models"
	"gitlab.com/_spacemc_/web/users/internal/domain/ports"

	"gitlab.com/_spacemc_/web/users/internal/query"

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

func (r *UsersRepository) GetByID(ctx context.Context, userID string) (*models.User, error) {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return nil, err
	}

	user, err := r.q.GetUserByID(ctx, userIDUUID)

	if err != nil {
		return nil, err
	}

	userModel := models.User{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin.String,
		Cloak:        user.Cloak.String,
		RegisteredAt: user.RegisteredAt.Time,
		IsActive:     user.IsActive,
	}

	return &userModel, nil
}

func (r *UsersRepository) GetByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	user, err := r.q.GetUserByEmail(ctx, userEmail)

	if err != nil {
		return nil, err
	}

	userModel := models.User{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin.String,
		Cloak:        user.Cloak.String,
		RegisteredAt: user.RegisteredAt.Time,
		IsActive:     user.IsActive,
	}

	return &userModel, nil
}

func (r *UsersRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := r.q.GetUserByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	userModel := models.User{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin.String,
		Cloak:        user.Cloak.String,
		RegisteredAt: user.RegisteredAt.Time,
		IsActive:     user.IsActive,
	}

	return &userModel, nil
}

func (r *UsersRepository) GetSkin(ctx context.Context, userID string) (string, error) {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return "", err
	}

	skinURL, err := r.q.GetUserSkin(ctx, userIDUUID)

	if err != nil {
		return "", err
	}

	return skinURL.String, nil
}

func (r *UsersRepository) GetCloak(ctx context.Context, userID string) (string, error) {

	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return "", err
	}

	cloakURL, err := r.q.GetUserCloak(ctx, userIDUUID)

	if err != nil {
		return "", err
	}

	return cloakURL.String, nil
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

func (r *UsersRepository) DeleteSkin(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	err = r.q.DeleteUserSkin(ctx, userIDUUID)

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

func (r *UsersRepository) DeleteCloak(ctx context.Context, userID string) error {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	err = r.q.DeleteUserCloak(ctx, userIDUUID)

	if err != nil {
		return err
	}

	return nil

}

func (r *UsersRepository) UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error) {
	userIDUUID, err := uuid.Parse(userID)

	if err != nil {
		return "", err
	}

	params := query.UpdateUserStatusParams{
		IsActive: active,
		ID:       userIDUUID,
	}

	id, err := r.q.UpdateUserStatus(ctx, params)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}
