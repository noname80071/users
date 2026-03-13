// Маршруты для запросов

package http

import (
	"gitlab.com/_spacemc_/web/users/internal/domain/ports"
	"gitlab.com/_spacemc_/web/users/internal/infra/http/controllers/files"
	"gitlab.com/_spacemc_/web/users/internal/infra/http/controllers/users"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gitlab.com/_spacemc_/web/gokit/adapters/keycloak"
	"gitlab.com/_spacemc_/web/gokit/ginx/middlewares"
	_ "gitlab.com/_spacemc_/web/users/docs"
)

type Deps struct {
	keycloakConfig   *keycloak.Config
	Logger           *zap.Logger
	UsersServicePort ports.UsersServicePort
	FilesServicePort ports.FilesServicePort
}

func NewRouter(deps Deps) *gin.Engine {
	router := gin.Default()

	kc := keycloak.NewAdapter(deps.keycloakConfig, deps.Logger)

	usersHandler := users.NewHandler(deps.UsersServicePort)
	filesHandler := files.NewHandler(deps.FilesServicePort)

	api := router.Group("/api/v1")
	{
		api.POST("users/", usersHandler.UserRegister)
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		protected := api.Group("/")
		protected.Use(middlewares.KeycloakAuthenticated(kc))
		{
			api.GET("users/:id", usersHandler.GetUserByID)
			api.GET("users/email/:email", usersHandler.GetUserByEmail)
			api.GET("users/username/:username", usersHandler.GetUserByUsername)
			api.GET("users/:id/skin", filesHandler.GetSkin)
			api.GET("users/:id/cloak", filesHandler.GetCloak)

			api.POST("users/:id/skin", filesHandler.UploadSkin)
			api.POST("users/:id/cloak", filesHandler.UploadCloak)

			api.PATCH("users/:id/skin", filesHandler.DeleteSkin)
			api.PATCH("users/:id/cloak", filesHandler.DeleteCloak)
			api.PATCH("users/:id", usersHandler.UpdateUserStatus) // Активация/деактивация пользователя
		}
	}

	return router
}
