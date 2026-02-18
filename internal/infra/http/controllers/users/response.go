package users

import (
	"time"
)

type UsersGetByIdResponse struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Avatar       string    `json:"avatar"`
	Skin         string    `json:"skin"`
	Cloak        string    `json:"cloak"`
	RegisteredAt time.Time `json:"registered_at"`
	IsActive     bool      `json:"is_active"`
}
