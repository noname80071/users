package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password_hash"`
	Avatar       string             `json:"avatar"`
	Skin         pgtype.Text        `json:"skin"`
	Cloak        pgtype.Text        `json:"cloak"`
	RegisteredAt pgtype.Timestamptz `json:"registered_at"`
	IsActive     bool               `json:"is_active"`
}
