package config

// ChatServerGRPCConfig Конфигурация текущего GRPC сервера
type ChatServerGRPCConfig struct {
	Host string `env:"GRPC_HOST" envDefault:"localhost"`
	Port string `env:"GRPC_PORT" envDefault:"50052"`
}
