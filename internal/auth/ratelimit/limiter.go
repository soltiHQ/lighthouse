package ratelimit

import (
	"sync"
	"time"
)

// Config defines rate limiter parameters.
type Config struct {
	MaxAttempts int           // e.g. 5
	BlockWindow time.Duration // e.g. 10 * time.Minute
}

type entry struct {
	failures     int
	blockedUntil time.Time
}

// Limiter tracks failed attempts by key and blocks after threshold.
type Limiter struct {
	mu      sync.Mutex
	entries map[string]*entry
	cfg     Config
}

// New creates a rate limiter.
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

// Blocked reports whether the key is currently blocked.
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
	// Block window expired — reset.
	if !e.blockedUntil.IsZero() {
		delete(l.entries, key)
	}
	return false
}

// RecordFailure increments failure count and blocks if threshold exceeded.
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

// Reset clears failure count for a key (call on successful login).
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}
