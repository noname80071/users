// Загрузка конфига

package config

import (
	"go-users/internal/infra/http"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   http.ServerConfig `env-prefix:"SERVER_"`
	Database DatabaseConfig    `env-prefix:"DB_"`
	Minio    MinioConfig       `env-prefix:"MINIO_"`
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

type MinioConfig struct {
	Endpoint        string `env:"ENDPOINT"`
	AccessKey       string `env:"ACCESS_KEY"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	SSLMode         bool   `env:"SSL_MODE"`
	Region          string `env:"REGION"`
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
