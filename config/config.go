// Загрузка конфига

package config

import (
	"go-users/internal/infra/http"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server http.ServerConfig `env-prefix:"SERVER_"`
	//Nats       nats.NatsConfig             `env-prefix:"NATS_"`
	Database DatabaseConfig `env-prefix:"DB_"`
	//Logging    logger.LoggingConfig        `env-prefix:"LOG_"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	DBName   string `env:"NAME"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	SSLMode  string `env:"SSL_MODE"`
	MaxConns int    `env:"MAX_CONNS"`
	MinConns int    `env:"MIN_CONNS"`
}

func Load() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("../../.env", &cfg)
	if err != nil {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}
