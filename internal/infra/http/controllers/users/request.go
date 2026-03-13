package users

type UserRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserStatusRequest struct {
	Active bool `json:"is_active"`
}
