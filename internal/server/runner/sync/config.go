package sync

import "time"

const (
	defaultName         = "sync"
	defaultTickInterval = 10 * time.Second
	defaultPushTimeout  = 15 * time.Second
	defaultMaxRetries   = 5
)

// Config configures the sync runner.
type Config struct {
	Name         string
	TickInterval time.Duration
	PushTimeout  time.Duration
	MaxRetries   int
}

func (c Config) withDefaults() Config {
	if c.Name == "" {
		c.Name = defaultName
	}
	if c.TickInterval <= 0 {
		c.TickInterval = defaultTickInterval
	}
	if c.PushTimeout <= 0 {
		c.PushTimeout = defaultPushTimeout
	}
	if c.MaxRetries <= 0 {
		c.MaxRetries = defaultMaxRetries
	}
	return c
}
