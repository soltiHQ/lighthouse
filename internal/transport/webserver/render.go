package webserver

import (
	"html/template"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/soltiHQ/control-plane/ui"
)

type renderer struct {
	tmpl   *template.Template
	logger zerolog.Logger
}

func newRenderer(logger zerolog.Logger) (*renderer, error) {
	t, err := template.ParseFS(ui.Templates, "templates/**/*.html")
	if err != nil {
		return nil, err
	}
	return &renderer{
		tmpl:   t,
		logger: logger.With().Str("component", "renderer").Logger(),
	}, nil
}

func (r *renderer) Render(w http.ResponseWriter, name string, data any) error {
	if err := r.tmpl.ExecuteTemplate(w, name, data); err != nil {
		r.logger.Error().Err(err).
			Str("template", name).
			Msg("render failed")
		return err
	}
	return nil
}
