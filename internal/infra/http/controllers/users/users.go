// Контроллер - принимает запрос - вызывает сервисы - возвращает ответ

package users

import (
	"errors"

	usersErrors "gitlab.com/_spacemc_/web/users/errors"
	"gitlab.com/_spacemc_/web/users/internal/domain/ports"

	"github.com/gin-gonic/gin"
	"gitlab.com/_spacemc_/web/gokit/ginx"
)

type Handler struct {
	service ports.UsersServicePort
}

func NewHandler(service ports.UsersServicePort) *Handler {
	return &Handler{service: service}
}

// GetUserByID godoc
// @Summary      Получить пользователя по ID
// @Description  Возвращает информацию о пользователе по его уникальному идентификатору
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "ID пользователя"
// @Success      200  {object} ginx.SuccessResponse[UsersGetByResponse] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Router       /users/{id} [get]
func (h *Handler) GetUserByID(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userID, err := p.GetPathString("id")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), *userID)

	if err != nil {
		switch {
		case errors.Is(err, usersErrors.ErrUserNotFound):
			ginx.WriteNotFound(c)
		}

		return

	}

	response := UsersGetByResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin,
		Cloak:        user.Cloak,
		RegisteredAt: user.RegisteredAt,
		IsActive:     user.IsActive,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// GetUserByEmail godoc
// @Summary      Получить пользователя по Email
// @Description  Возвращает информацию о пользователе по его Email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Email пользователя"
// @Success      200  {object} ginx.SuccessResponse[UsersGetByResponse] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Router       /users/email/{email} [get]
func (h *Handler) GetUserByEmail(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userEmail, err := p.GetPathString("email")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	user, err := h.service.GetUserByEmail(c.Request.Context(), *userEmail)

	if err != nil {
		switch {
		case errors.Is(err, usersErrors.ErrUserNotFound):
			ginx.WriteNotFound(c)
		}

		return

	}

	response := UsersGetByResponse{
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin,
		Cloak:        user.Cloak,
		RegisteredAt: user.RegisteredAt,
		IsActive:     user.IsActive,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// GetUserByUsername godoc
// @Summary      Получить пользователя по Username
// @Description  Возвращает информацию о пользователе по его Username
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Username пользователя"
// @Success      200  {object} ginx.SuccessResponse[UsersGetByResponse] "Успешный ответ"
// @Failure      404  {object}  ginx.ErrorResponse "Пользователя не существует"
// @Router       /users/username/{username} [get]
func (h *Handler) GetUserByUsername(c *gin.Context) {
	p := ginx.NewGinxParser(c)

	userUsername, err := p.GetPathString("username")

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	user, err := h.service.GetUserByUsername(c.Request.Context(), *userUsername)

	if err != nil {
		switch {
		case errors.Is(err, usersErrors.ErrUserNotFound):
			ginx.WriteNotFound(c)
		}

		return

	}

	response := UsersGetByResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin,
		Cloak:        user.Cloak,
		RegisteredAt: user.RegisteredAt,
		IsActive:     user.IsActive,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// UserRegister godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с указанными username, email и паролем
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body UserRegisterRequest true "Данные для регистрации"
// @Success      201  {object}  ginx.SuccessResponse[UserRegisterResponse] "Успешный ответ"
// @Failure      400  {object}  ginx.ErrorResponse  "Неверный формат запроса или ошибка валидации"
// @Failure      422  {object}  ginx.ErrorResponse  "Пользователь уже существует"
// @Router       /users/register [post]
func (h *Handler) UserRegister(c *gin.Context) {

	var request UserRegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
	}

	id, err := h.service.RegisterUser(c.Request.Context(),
		request.Username,
		request.Email,
		request.Password,
	)

	if err != nil {
		ginx.WriteErrorResponse(c, ginx.InternalServerError)
		return
	}

	response := UserRegisterResponse{
		ID: id,
	}

	ginx.WriteSuccessResponse(c, &response)
}

// UpdateUserStatus godoc
// @Summary      Обновление статуса пользователя
// @Description  Активирует или деактивирует пользователя по его ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id      path     string                   true "ID пользователя"
// @Param        request body     UpdateUserStatusRequest true "Новый статус пользователя"
// @Success      200     {object} ginx.SuccessResponse[UpdateUserStatusResponse] "Успешный запрос"
// @Failure      400     {object} ginx.ErrorResponse                  "Неверный формат запроса"
// @Failure      404     {object} ginx.ErrorResponse                  "Пользователь не найден"
// @Failure      422     {object} ginx.ErrorResponse                  "Ошибка бизнес-логики"
// @Router       /users/{id}/status [patch]
func (h *Handler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var request UpdateUserStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		ginx.WriteErrorResponse(c, ginx.BadRequest)
		return
	}

	id, err := h.service.UpdateUserStatus(c.Request.Context(), userID, request.Active)

	if err != nil {
		// switch {
		// case errors.Is(err, serviceErrors.ErrUserNotFound):
		// 	c.JSON(http.StatusNotFound, "User not found")
		// }
		return
	}

	response := UpdateUserStatusResponse{
		ID: id,
	}

	ginx.WriteSuccessResponse(c, &response)

}
