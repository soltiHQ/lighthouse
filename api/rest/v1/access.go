package restv1

// LoginRequest is the request body for password-based authentication.
type LoginRequest struct {
	Subject  string `json:"subject"`
	Password string `json:"password"`
}

// LoginResponse is the token pair returned on successful login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	SessionID    string `json:"session_id"`
}

// LogoutRequest identifies the session to terminate.
type LogoutRequest struct {
	SessionID string `json:"session_id"`
}
