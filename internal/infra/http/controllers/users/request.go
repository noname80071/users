package users

import (
	"github.com/google/uuid"
)

type UsersGetByIdRequest struct {
	ID uuid.UUID `json:"id"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
