package http

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"gitlab.com/_spacemc_/web/users/internal/domain/models"
)

type InfrastructureService interface {
	Start(ctx context.Context) error
	GracefulShutdown(ctx context.Context) error
}

type UsersService interface {
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	RegisterUser(ctx context.Context, username string, email string, password string) (string, error)
	UpdateUserStatus(ctx context.Context, userID string, active bool) (string, error)
}

type FilesService interface {
	GetUserSkin(ctx context.Context, userID string) (string, error)
	GetUserCloak(ctx context.Context, userID string) (string, error)

	UploadSkin(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	DeleteSkin(ctx context.Context, userID string) error

	UploadCloak(ctx context.Context, userID string, fileReader io.Reader, filename string, fileSize int64) (string, error)
	DeleteCloak(ctx context.Context, userID string) error
}

type HttpService struct {
	address string
	server  *http.Server
}

// @title           Go Microservice Users API
// @version         1.0
// @description     Users microservice

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Тип: Bearer {token}. Токен получается через Keycloak.
func NewServer(cfg *ServerConfig, usersService UsersService, filesService FilesService) InfrastructureService {
	router := NewRouter(Deps{
		UsersService: usersService,
		FilesService: filesService,
	})
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return &HttpService{
		address: address,
		server: &http.Server{
			Addr:         address,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
			Handler:      router,
		},
	}
}

func (s *HttpService) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.address, err)
	}

	return s.server.Serve(lis)
}

func (s *HttpService) GracefulShutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
