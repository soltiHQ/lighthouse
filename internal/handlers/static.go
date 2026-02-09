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

// NewStatic returns a new Static handler.
func NewStatic(logger zerolog.Logger) *Static {
	sub, err := fs.Sub(ui.Static, "static")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load static files")
	}
	return &Static{
		logger:  logger,
		handler: http.StripPrefix("/static/", http.FileServer(http.FS(sub))),
	}
}

func (s *Static) Routes(mux *http.ServeMux) {
	mux.Handle("/static/", s.handler)

	mux.HandleFunc("/static", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/", http.StatusMovedPermanently)
	})
}
