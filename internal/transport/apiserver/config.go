package apiserver

import (
	"github.com/soltiHQ/control-plane/auth/authenticator"
	"github.com/soltiHQ/control-plane/internal/transport/config"

	"github.com/rs/zerolog"
)

// Config represents the configuration for the api server.
type Config struct {
	configHTTP config.HttpConfig

	addrHTTP string

	logLevel zerolog.Level
	authn    authenticator.Authenticator
}

// NewConfig creates a new configuration instance.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		configHTTP: config.NewHttpConfig(),
		logLevel:   zerolog.InfoLevel,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}
