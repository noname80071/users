// Загрузка окружения - загрузка конфигурации - подключение к БД - создание сервисов - запуск сервера

package app

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/_spacemc_/web/users/config"
	"gitlab.com/_spacemc_/web/users/internal/domain/ports"
	httpInfra "gitlab.com/_spacemc_/web/users/internal/infra/http"
	filessvc "gitlab.com/_spacemc_/web/users/internal/infra/services/files"
	userssvc "gitlab.com/_spacemc_/web/users/internal/infra/services/users"
	"gitlab.com/_spacemc_/web/users/pkg/database"
	"gitlab.com/_spacemc_/web/users/pkg/minio"
)

type App struct {
	config     *config.Config
	enviroment string
	server     ports.InfrastructureService
	postgres   *database.PoolAdapter
}

func NewApp(ctx context.Context) (*App, error) {
	enviroment, ok := os.LookupEnv("ENV")

	if !ok {
		enviroment = "default"
	}

	configImpl, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}

	pgxConnection, err := database.NewPgxConnection(ctx, &configImpl.Database)

	poolAdapter := &database.PoolAdapter{Pool: pgxConnection}

	// minio
	minioClient, err := minio.NewClient(ctx, &configImpl.Minio)

	// services
	usersServicePort := userssvc.New(pgxConnection)
	filesServicePort := filessvc.New(pgxConnection, minioClient)

	httpServer := httpInfra.NewServer(&configImpl.Server, usersServicePort, filesServicePort)
	app := &App{
		config:     configImpl,
		enviroment: enviroment,
		server:     httpServer,
		postgres:   poolAdapter,
	}

	return app, nil
}

func (app *App) Shutdown() []error {
	var errors []error

	downs := []func(context.Context) error{
		func(ctx context.Context) error {
			if app.server != nil {
				if err := app.server.GracefulShutdown(ctx); err != nil {
					return fmt.Errorf("http server graceful shutdown: %w", err)
				}
			}

			return nil
		},
		func(ctx context.Context) error {
			//app.postgres.Close()

			return nil
		},
	}

	for _, task := range downs {
		if err := task(context.Background()); err != nil {
			errors = append(errors, fmt.Errorf("GracefulShutdownError: %w", err))
		}
	}

	return errors
}

func (app *App) Start(ctx context.Context) error {
	errChan := make(chan error, 2)

	go func() {
		if err := app.server.Start(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start HTTP server: %w", err)
		}
	}()

	return <-errChan
}
