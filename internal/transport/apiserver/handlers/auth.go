package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/soltiHQ/control-plane/auth/authenticator"
	"github.com/soltiHQ/control-plane/internal/logctx"
	"github.com/soltiHQ/control-plane/internal/transport/response"
)

type tokenRequest struct {
	Subject  string `json:"subject"`
	Password string `json:"password"`
}

type tokenIdentityResponse struct {
	Subject     string   `json:"subject"`
	UserID      string   `json:"user_id"`
	Permissions []string `json:"permissions,omitempty"`
}

type tokenResponse struct {
	AccessToken string                 `json:"access_token"`
	TokenType   string                 `json:"token_type"`
	ExpiresIn   int64                  `json:"expires_in"` // seconds
	Identity    *tokenIdentityResponse `json:"identity,omitempty"`
}

// Token handles POST /v1/auth/token.
// It validates credentials and returns a signed JWT access token.
func (h *Http) Token(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logctx.From(ctx, h.logger)

	if r.Method != http.MethodPost {
		_ = response.NotAllowed(ctx, w, "method not supported")
		return
	}

	if h.authn == nil {
		_ = response.NotFound(ctx, w, "not found")
		return
	}

	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.BadRequest(ctx, w, "invalid JSON payload")
		return
	}

	token, id, err := h.authn.Authenticate(ctx, &authenticator.Request{
		Subject:  req.Subject,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, authenticator.ErrInvalidCredentials):
			_ = response.Unauthorized(ctx, w, "invalid credentials")
			logger.Warn().Err(err).Msg("authentication failed")
			return

		case errors.Is(err, authenticator.ErrUnauthorized):
			_ = response.Forbidden(ctx, w, "forbidden")
			logger.Warn().Err(err).Msg("authentication forbidden")
			return

		default:
			logger.Error().Err(err).Msg("authentication internal error")
			_ = response.FromError(ctx, w, err)
			return
		}
	}
	if id != nil {
		id.RawToken = ""
	}

	var expiresIn int64
	if id != nil && !id.ExpiresAt.IsZero() {
		d := time.Until(id.ExpiresAt)
		if d < 0 {
			d = 0
		}
		expiresIn = int64(d.Seconds())
	}
	var outID *tokenIdentityResponse
	if id != nil {
		outID = &tokenIdentityResponse{
			Subject:     id.Subject,
			UserID:      id.UserID,
			Permissions: append([]string(nil), id.Permissions...),
		}
	}
	_ = response.OK(ctx, w, tokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		Identity:    outID,
	})
}
