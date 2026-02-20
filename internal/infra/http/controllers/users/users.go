// Контроллер - принимает запрос - вызывает сервисы - возвращает ответ

package users

import (
	"errors"
	serviceErrors "go-users/errors"
	"go-users/internal/domain/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ports.UsersServicePort
}

func NewHandler(service ports.UsersServicePort) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.service.GetUserByID(c.Request.Context(), userID)

	if err != nil {
		switch {
		case errors.Is(err, serviceErrors.ErrUserNotFound):
			c.JSON(http.StatusNotFound, "User not found")
		}

		return

	}

	response := UsersGetByIdResponse{
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin,
		Cloak:        user.Cloak,
		RegisteredAt: user.RegisteredAt,
		IsActive:     user.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UserRegister(c *gin.Context) {

	var request UserRegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	id, err := h.service.RegisterUser(c.Request.Context(),
		request.Username,
		request.Email,
		request.Password,
	)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
