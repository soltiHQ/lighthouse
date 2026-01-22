package handlers

import (
	"html/template"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/soltiHQ/control-plane/internal/logctx"
)

// Pages represent the pages handler.
type Pages struct {
	logger    zerolog.Logger
	templates *template.Template
}

// NewPages creates a new pages handler.
func NewPages(logger zerolog.Logger, tmpl *template.Template) *Pages {
	return &Pages{
		logger:    logger.With().Str("type", "pages").Logger(),
		templates: tmpl,
	}
}

func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = logctx.From(ctx, p.logger)
	)

	data := map[string]any{
		"Title": "Home Page",
	}

	if err := p.templates.ExecuteTemplate(w, "base", data); err != nil {
		logger.Error().Err(err).Msg("failed to render home template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
