package config

import (
	"time"

	"github.com/joho/godotenv"
)

// GRPCConfig Интерфейс конфигурации GRPC.
type GRPCConfig interface {
	Address() string
}

// PGConfig Интерфейс конфигурации для клиента Postgres.
type PGConfig interface {
	DSN() string
}

// RedisConfig Интерфейс конфигурации для Redis.
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// Load Configs.
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}
