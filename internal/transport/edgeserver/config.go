package edgeserver

import (
	"github.com/soltiHQ/control-plane/internal/transport/config"

	"github.com/rs/zerolog"
)

// Config represents the configuration for the edge server.
type Config struct {
	configHTTP config.HttpConfig
	configGRPC config.GrpcConfig

	addrHTTP string
	addrGRPC string

	logLevel zerolog.Level
}

// NewConfig creates a new configuration instance.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		configHTTP: config.NewHttpConfig(),
		configGRPC: config.NewGrpcConfig(),
		logLevel:   zerolog.InfoLevel,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
