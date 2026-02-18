package http

import "time"

type ServerConfig struct {
	Host         string        `env:"HOST" env-default:"localhost"`
	Port         int           `env:"PORT" env-default:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" env-default:"15s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" env-default:"15s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}
