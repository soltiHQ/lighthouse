package server

import "errors"

var (
	// ErrRunnerExited indicates a runner exited unexpectedly.
	ErrRunnerExited = errors.New("server: runner exited unexpectedly")
	// ErrNoRunners indicates the server has no runners to start.
	ErrNoRunners = errors.New("server: no runners configured")
	// ErrNilContext indicates a nil context was provided.
	ErrNilContext = errors.New("server: nil context")
)
