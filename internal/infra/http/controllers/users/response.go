package users

import (
	"time"
)

type UsersGetByResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	Skin         string    `json:"skin"`
	Cloak        string    `json:"cloak"`
	RegisteredAt time.Time `json:"registered_at"`
	IsActive     bool      `json:"is_active"`
}

type UserRegisterResponse struct {
	ID string `json:"id"`
}

type UpdateUserStatusResponse struct {
	ID string `json:"id"`
}
