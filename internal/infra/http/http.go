package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"go-users/internal/domain/ports"
)

type HttpService struct {
	address string
	server  *http.Server
}

func NewServer(cfg *ServerConfig, usersServicePort ports.UsersServicePort, filesServicePort ports.FilesServicePort) ports.InfrastructureService {
	router := NewRouter(Deps{
		UsersServicePort: usersServicePort,
		FilesServicePort: filesServicePort,
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
