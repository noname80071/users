// Контроллер - принимает запрос - вызывает сервисы - возвращает ответ

package users

import (
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

func (h *Handler) GetUserById(c *gin.Context) {
	userId := c.Param("id")

	user, err := h.service.GetUserById(c.Request.Context(), userId)

	if err != nil {
		return
	}

	response := UsersGetByIdResponse{
		Username:     user.Username,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Skin:         user.Skin.String,
		Cloak:        user.Cloak.String,
		RegisteredAt: user.RegisteredAt.Time,
		IsActive:     user.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UserRegister(c *gin.Context) {

	var request UserRegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	err := h.service.RegisterUser(c.Request.Context(),
		request.Username,
		request.Email,
		request.Password,
	)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, "")
}
