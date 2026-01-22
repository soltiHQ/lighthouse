package handlers

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/soltiHQ/control-plane/ui"
)

type Static struct {
	logger zerolog.Logger
	fs     http.FileSystem
}

func NewStatic(logger zerolog.Logger) *Static {
	staticFS, err := fs.Sub(ui.Static, "static")
	if err != nil {
		panic("failed to create static sub-fs: " + err.Error())
	}
	return &Static{
		logger: logger.With().Str("type", "static").Logger(),
		fs:     http.FS(staticFS),
	}
}

func (s *Static) Serve() http.Handler {
	fileServer := http.FileServer(s.fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := strings.TrimPrefix(r.URL.Path, "/static/")
		r.URL.Path = upath

		if strings.HasSuffix(upath, "/") {
			http.NotFound(w, r)
			return
		}

		ext := path.Ext(upath)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}

		w.Header().Set("Cache-Control", "public, max-age=31536000")

		fileServer.ServeHTTP(w, r)
	})
}
