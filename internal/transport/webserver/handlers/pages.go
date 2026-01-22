package handlers

import (
	"net/http"

	"github.com/rs/zerolog"
)

type Pages struct {
	logger   zerolog.Logger
	renderer Renderer
}

func NewPages(logger zerolog.Logger, renderer Renderer) *Pages {
	return &Pages{
		logger:   logger.With().Str("handler", "pages").Logger(),
		renderer: renderer,
	}
}

func (p *Pages) Home(w http.ResponseWriter, r *http.Request) {
	type ViewData struct {
		Title string
	}

	if err := p.renderer.Render(w, "home.html", ViewData{
		Title: "Control Plane",
	}); err != nil {
		p.logger.Error().Err(err).Msg("failed to render home page")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
