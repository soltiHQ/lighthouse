package ratelimit

import (
	"sync"
	"time"

	"github.com/soltiHQ/control-plane/internal/auth"
)

// Config defines rate limiting parameters.
//
// MaxAttempts specifies how many failed attempts are allowed
// before the key becomes temporarily blocked.
//
// BlockWindow specifies the duration of the block once
// the threshold is reached.
//
// Zero or negative values are replaced with safe defaults in New().
type Config struct {
	MaxAttempts int
	BlockWindow time.Duration
}

type entry struct {
	failures     int
	blockedUntil time.Time
}

// Limiter provides an in-memory, concurrency-safe rate limiter
// for authentication attempts identified by an arbitrary string key.
//
// The limiter is process-local and does not provide distributed guarantees.
// State is not persisted and is lost on restart.
type Limiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	cfg     Config
}

// New creates a new Limiter with normalized configuration.
//
// If MaxAttempts <= 0, a default value is used.
// If BlockWindow <= 0, a default value is used.
func New(cfg Config) *Limiter {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 5
	}
	if cfg.BlockWindow <= 0 {
		cfg.BlockWindow = 10 * time.Minute
	}
	return &Limiter{
		entries: make(map[string]*entry),
		cfg:     cfg,
	}
}

// Check returns auth.ErrRateLimited if the key is currently blocked.
func (l *Limiter) Check(key string, now time.Time) error {
	if l.Blocked(key, now) {
		return auth.ErrRateLimited
	}
	return nil
}

// Blocked reports whether the key is currently blocked at time now.
//
// If a previous block has expired, the internal state for the key
// is cleared and false is returned.
//
// Safe for concurrent use.
func (l *Limiter) Blocked(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok {
		return false
	}
	if !e.blockedUntil.IsZero() && now.Before(e.blockedUntil) {
		return true
	}
	if !e.blockedUntil.IsZero() {
		delete(l.entries, key)
	}
	return false
}

// RecordFailure records a failed attempt for the key.
//
// Once failures reach or exceed MaxAttempts, the key becomes blocked
// until now + BlockWindow.
//
// Safe for concurrent use.
func (l *Limiter) RecordFailure(key string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok {
		e = &entry{}
		l.entries[key] = e
	}
	e.failures++
	if e.failures >= l.cfg.MaxAttempts {
		e.blockedUntil = now.Add(l.cfg.BlockWindow)
	}
}

// Reset removes all rate-limit state associated with the key.
//
// Safe for concurrent use.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}
