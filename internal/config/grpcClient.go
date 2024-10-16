package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	grpcClientHostEnvName = "GRPC_AUTH_HOST"
	grpcClientPortEnvName = "GRPC_AUTH_PORT"
)

// GRPCClientConfig GRPCConfig GRPC config
type GRPCClientConfig interface {
	Address() string
}

type grpcClientConfig struct {
	host string
	port string
}

// NewGRPCClientConfig NewGRPCConfig GRPC config constructor
func NewGRPCClientConfig() (GRPCClientConfig, error) {
	host := os.Getenv(grpcClientHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcClientPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcClientConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcClientConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
