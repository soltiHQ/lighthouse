package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/soltiHQ/control-plane/internal/auth"
	"github.com/soltiHQ/control-plane/internal/auth/svc"
	"github.com/soltiHQ/control-plane/internal/backend"
	"github.com/soltiHQ/control-plane/internal/transport/http/responder"
	"github.com/soltiHQ/control-plane/internal/transport/http/response"
)

// API represents a web API handler.
type API struct {
	logger  zerolog.Logger
	backend *backend.Login
	auth    *svc.Auth
}

// NewAPI creates a new API handler.
func NewAPI(logger zerolog.Logger, auth *svc.Auth, backend *backend.Login) *API {
	return &API{logger: logger, auth: auth, backend: backend}
}

type loginRequest struct {
	Subject  string `json:"subject"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}

// Login handles POST /api/v1/login (JSON).
func (x *API) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, r, response.RenderBlock)
		return
	}

	// key: оставь как есть пока; позже сделаешь subject+ip/ua
	key := req.Subject

	_, access, refresh, sessionID, err := x.backend.Do(r.Context(), req.Subject, req.Password, key)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrRateLimited):
			response.AuthRateLimit(w, r, response.RenderBlock)
		case errors.Is(err, auth.ErrInvalidRequest):
			response.BadRequest(w, r, response.RenderBlock)
		case errors.Is(err, auth.ErrInvalidCredentials), errors.Is(err, auth.ErrUnauthorized):
			response.Unauthorized(w, r, response.RenderBlock)
		default:
			x.logger.Warn().Err(err).Msg("api login failed")
			response.Unavailable(w, r, response.RenderBlock)
		}
		return
	}

	response.OK(w, r, response.RenderBlock, &responder.View{
		Data: loginResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			SessionID:    sessionID,
		},
	})
}
