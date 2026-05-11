package http

import (
	"context"
	"io"

	"gitlab.com/_spacemc_/web/users/internal/domain/models"
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

type usersService interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	RegisterUser(ctx context.Context, username string, email string, password string) (string, error)
	UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error)
}

type filesService interface {
	GetUserSkin(ctx context.Context, userID string) (string, error)
	GetUserCloak(ctx context.Context, userID string) (string, error)

	UploadSkin(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	DeleteSkin(ctx context.Context, userID string) error

	UploadCloak(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	DeleteCloak(ctx context.Context, userID string) error
}

type Deps struct {
	keycloakConfig *keycloak.Config
	Logger         *zap.Logger
	UsersService   usersService
	FilesService   filesService
}

func NewRouter(deps Deps) *gin.Engine {
	router := gin.Default()

	kc := keycloak.NewAdapter(deps.keycloakConfig, deps.Logger)

	usersHandler := users.NewHandler(deps.UsersService)
	filesHandler := files.NewHandler(deps.FilesService)

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
