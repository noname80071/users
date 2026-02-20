// Маршруты для запросов

package http

import (
	"go-users/internal/domain/ports"
	"go-users/internal/infra/http/controllers/files"
	"go-users/internal/infra/http/controllers/users"

	"github.com/gin-gonic/gin"
)

type Deps struct {
	UsersServicePort ports.UsersServicePort
	FilesServicePort ports.FilesServicePort
}

func NewRouter(deps Deps) *gin.Engine {
	router := gin.Default()

	usersHandler := users.NewHandler(deps.UsersServicePort)
	filesHandler := files.NewHandler(deps.FilesServicePort)

	api := router.Group("/api/v1")

	{
		api.GET("users/:id", usersHandler.GetUserByID)

		api.POST("users/", usersHandler.UserRegister)
		api.POST("users/:id/skin", filesHandler.UploadSkin)
		api.POST("users/:id/cloak", filesHandler.UploadCloak)
	}

	return router
}
