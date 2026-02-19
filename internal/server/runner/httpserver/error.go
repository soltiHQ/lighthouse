package httpserver

import "errors"

var (
	ErrNilHandler     = errors.New("httpserver: nil handler")
	ErrEmptyAddr      = errors.New("httpserver: empty addr")
	ErrAlreadyStarted = errors.New("httpserver: already started")
)
