package serviceErrors

import (
	"errors"
)

var (
	ErrUserAlreadyExists = errors.New("User with this email or username already exists")
	ErrUserNotFound      = errors.New("User not found")

	ErrUsernameInvalid = errors.New("Username must be a valid username")

	ErrEmailInvalid = errors.New("Email must be a valid email address")
	ErrEmailTaken   = errors.New("Email is already registered")

	ErrPasswordWeak = errors.New("Password is too weak")

	FailedToGetFile  = errors.New("Failed to get file")
	FailedToOpenFile = errors.New("Failed to open file")

	FailedToUploadFile = errors.New("Failed to upload file")
)
