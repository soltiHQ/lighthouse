package session

import "time"

// TokenPair is returned to callers after login/refresh.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Config controls session and token lifetimes and refresh rotation behavior.
type Config struct {
	// AccessTTL is the lifetime of issued access tokens.
	AccessTTL time.Duration
	// RefreshTTL is the lifetime of refresh tokens stored in sessions.
	RefreshTTL time.Duration
	// Issuer is embedded into issued access token identities.
	Issuer string
	// Audience is embedded into issued access token identities.
	Audience string
	// RotateRefresh controls whether Refresh rotates refresh tokens.
	RotateRefresh bool
}
