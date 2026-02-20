package models

import (
	"time"
)

type User struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Avatar       string    `json:"avatar"`
	Skin         string    `json:"skin"`
	Cloak        string    `json:"cloak"`
	RegisteredAt time.Time `json:"registered_at"`
	IsActive     bool      `json:"is_active"`
}
