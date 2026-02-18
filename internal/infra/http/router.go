// Маршруты для запросов

package http

import (
	"go-users/internal/domain/ports"
	"go-users/internal/infra/http/controllers/users"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	UsersServicePort ports.UsersServicePort
}

func NewRouter(deps Deps) *gin.Engine {
	router := gin.Default()

	usersHandler := users.NewHandler(deps.UsersServicePort)

	api := router.Group("/api/v1")

	{
		api.GET("users/:id", usersHandler.GetUserById)
		api.POST("users/", usersHandler.UserRegister)
	}

	return router
}
