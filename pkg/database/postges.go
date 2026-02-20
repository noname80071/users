// Подключение к БД

package database

import (
	"context"
	"fmt"

	"go-users/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxConnection interface {
	NewPgxConnection(ctx context.Context, cfg *config.DatabaseConfig) (*pgxpool.Pool, error)
}

type PoolAdapter struct {
	Pool *pgxpool.Pool
}

func NewPgxConnection(ctx context.Context, cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	// Строка подключения
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_min_conns=%d",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode, cfg.MaxConns, cfg.MinConns,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("Error parse configuration: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Error create pool: %w", err)
	}

	// Проверка соединения
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Error ping: %w", err)
	}

	return pool, nil
}
