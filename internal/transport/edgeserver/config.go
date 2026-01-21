package edgeserver

import (
	"github.com/soltiHQ/control-plane/internal/transport/config"

	"github.com/rs/zerolog"
)

type Config struct {
	addrHTTP string
	addrGRPC string

	configHTTP config.HttpConfig
	configGRPC config.GrpcConfig

	logLevel zerolog.Level
}

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
