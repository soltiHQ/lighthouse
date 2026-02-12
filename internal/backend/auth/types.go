package auth

// LoginRequest is an authentication request.
type LoginRequest struct {
	Subject  string
	Password string
	RateKey  string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
}

// LogoutRequest revokes a session.
type LogoutRequest struct {
	SessionID string
}
