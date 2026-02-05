package handlers

import (
	"io/fs"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/soltiHQ/control-plane/ui"
)

type Static struct {
	logger  zerolog.Logger
	handler http.Handler
}

func NewStatic(logger zerolog.Logger) *Static {
	sub, err := fs.Sub(ui.Static, "static")
	if err != nil {
		panic(err)
	}
	return &Static{
		logger:  logger.With().Str("handler", "static").Logger(),
		handler: http.StripPrefix("/static/", http.FileServer(http.FS(sub))),
	}
}

func (s *Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}
