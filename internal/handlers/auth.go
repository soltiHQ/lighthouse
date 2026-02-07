package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/soltiHQ/control-plane/internal/auth/auth/session"
	"github.com/soltiHQ/control-plane/internal/transport/http/response"
)

// Auth handles authentication endpoints.
type Auth struct {
	session *session.Service
	json    *response.JSONResponder
}

// NewAuth creates an auth handler.
func NewAuth(session *session.Service, json *response.JSONResponder) *Auth {
	return &Auth{
		session: session,
		json:    json,
	}
}

// Routes registers auth routes on the given mux.
// These routes are public — no Auth middleware required.
func (a *Auth) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/login", a.Login)
}

// loginRequest is the expected JSON body for login.
type loginRequest struct {
	Subject  string `json:"subject"`
	Password string `json:"password"`
}

// loginResponse is the JSON response on successful login.
type loginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	Subject      string   `json:"subject"`
	UserID       string   `json:"user_id"`
	Permissions  []string `json:"permissions"`
}

// Login authenticates by subject/password and returns a JWT token pair.
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.json.Error(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Subject == "" || req.Password == "" {
		a.json.Error(w, r, http.StatusBadRequest, "subject and password are required")
		return
	}

	pair, id, err := a.session.Login(r.Context(), req.Subject, req.Password)
	if err != nil {
		// session.Login returns auth.ErrInvalidCredentials, auth.ErrUnauthorized, etc.
		// Don't leak internals — always 401 for auth failures.
		a.json.Error(w, r, http.StatusUnauthorized, "invalid credentials")
		return
	}

	perms := make([]string, 0, len(id.Permissions))
	for _, p := range id.Permissions {
		perms = append(perms, string(p))
	}

	a.json.Respond(w, r, http.StatusOK, &response.View{
		Data: loginResponse{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
			ExpiresAt:    id.ExpiresAt.Unix(),
			Subject:      id.Subject,
			UserID:       id.UserID,
			Permissions:  perms,
		},
	})
}
