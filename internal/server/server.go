package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

// Server orchestrates the lifecycle of multiple runners.
type Server struct {
	cfg Config
	log zerolog.Logger

	runners []Runner
	stopped atomic.Bool
}

// New creates a Server. Cfg is normalized with defaults.
func New(cfg Config, log zerolog.Logger, runners ...Runner) *Server {
	rs := make([]Runner, len(runners))
	copy(rs, runners)

	return &Server{
		cfg:     cfg.withDefaults(),
		log:     log,
		runners: rs,
	}
}

// Run starts all runners and blocks until shutdown completes.
func (s *Server) Run(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	if len(s.runners) == 0 {
		return ErrNoRunners
	}
	trigger := make(chan error, 1)

	var wg sync.WaitGroup
	for _, r := range s.runners {
		wg.Add(1)
		go func(r Runner) {
			defer wg.Done()

			s.log.Info().Str("runner", r.Name()).Msg("starting")
			err := r.Start(ctx)

			if err != nil && !errors.Is(err, context.Canceled) {
				s.log.Error().Err(err).Str("runner", r.Name()).Msg("exited")
			} else {
				s.log.Info().Str("runner", r.Name()).Msg("exited")
			}

			select {
			case trigger <- err:
			default:
			}
		}(r)
	}

	var result error
	select {
	case <-ctx.Done():
		result = ctx.Err()
	case err := <-trigger:
		if err != nil {
			result = err
		} else {
			result = ErrRunnerExited
		}
	}

	if err := s.Shutdown(context.Background()); err != nil {
		s.log.Error().Err(err).Msg("shutdown completed with error")
	}

	wg.Wait()

	return result
}

// Shutdown gracefully stops all runners in reverse registration order.
// Idempotent: only the first call performs work.
func (s *Server) Shutdown(ctx context.Context) error {
	if ctx == nil {
		return ErrNilContext
	}
	if !s.stopped.CompareAndSwap(false, true) {
		return nil
	}
	return s.doShutdown(ctx)
}

func (s *Server) doShutdown(ctx context.Context) error {
	if s.cfg.ShutdownTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
		defer cancel()
	}

	s.log.Info().
		Dur("timeout", s.cfg.ShutdownTimeout).
		Int("runners", len(s.runners)).
		Msg("shutting down")

	var errs []error

	for i := len(s.runners) - 1; i >= 0; i-- {
		r := s.runners[i]

		s.log.Info().Str("runner", r.Name()).Msg("stopping")

		if err := r.Stop(ctx); err != nil {
			s.log.Error().Err(err).Str("runner", r.Name()).Msg("stop failed")
			errs = append(errs, fmt.Errorf("%s: %w", r.Name(), err))
		}
	}

	return errors.Join(errs...)
}
