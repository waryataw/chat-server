package config

// AuthGRPCConfig Конфигурация для подключения к Auth сервису по GRPC
type AuthGRPCConfig struct {
	Host string `env:"GRPC_AUTH_HOST" envDefault:"localhost"`
	Port string `env:"GRPC_AUTH_PORT" envDefault:"50052"`
}
