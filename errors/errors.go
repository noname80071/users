package serviceErrors

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("User with this email or username already exists")
	ErrUserNotFound      = errors.New("User not found")
	ErrInvalidEmail      = errors.New("Invalid email format")

	FailedToGetFile  = errors.New("Failed to get file")
	FailedToOpenFile = errors.New("Failed to open file")

	FailedToUploadFile = errors.New("Failed to upload file")
)
