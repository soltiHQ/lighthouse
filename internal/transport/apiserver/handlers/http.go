package handlers

import (
	"net/http"

	"github.com/soltiHQ/control-plane/internal/backend"
	"github.com/soltiHQ/control-plane/internal/logctx"
	"github.com/soltiHQ/control-plane/internal/storage"
	"github.com/soltiHQ/control-plane/internal/transport/response"

	"github.com/rs/zerolog"
)

// Http implements the HTTP api service.
type Http struct {
	logger  zerolog.Logger
	storage storage.Storage
}

// NewHttp creates a new HTTP api handler.
func NewHttp(logger zerolog.Logger, storage storage.Storage) *Http {
	return &Http{
		logger: logger.With().
			Str("type", "http").
			Logger(),
		storage: storage,
	}
}

// AgentList handles HTTP agent list request.
func (h *Http) AgentList(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = logctx.From(ctx, h.logger)
	)
	if r.Method != http.MethodGet {
		logger.Warn().Str("method", r.Method).Msg("invalid method")
		if err := response.NotAllowed(ctx, w, "method not supported"); err != nil {
			logctx.Error(ctx, h.logger, err, "failed to write not-allowed response")
		}
		return
	}

	_, err := backend.AgentList(ctx, logger, h.storage)
	if err != nil {
		if err = response.InternalError(ctx, w, "internal error"); err != nil {
			logctx.Error(ctx, h.logger, err, "failed to write internal-error response")
		}
		return
	}
	if err = response.OK(ctx, w, "mock"); err != nil {
		logctx.Error(ctx, h.logger, err, "failed to write ok response")
	}
}
