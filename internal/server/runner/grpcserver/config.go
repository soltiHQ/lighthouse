package grpcserver

const (
	defaultName    = "grpc"
	defaultNetwork = "tcp"
)

// Config controls the gRPC server runner behavior.
type Config struct {
	Name    string
	Addr    string
	Network string
}

func (c Config) withDefaults() Config {
	if c.Name == "" {
		c.Name = defaultName
	}
	if c.Network == "" {
		c.Network = defaultNetwork
	}
	return c
}
