// Package cookie manages site cookies, used by the browser-facing UI flow.
package cookie

import "net/http"

// SetAuth writes the three auth cookies to the response.
func SetAuth(w http.ResponseWriter, r *http.Request, accessToken, refreshToken, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   r.TLS != nil,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
}

// DeleteAuth expires all three auth cookies.
func DeleteAuth(w http.ResponseWriter, r *http.Request) {
	secure := r.TLS != nil

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   secure,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
	})
}

// GetAccessToken returns the access token cookie.
func GetAccessToken(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("access_token")
}

// GetRefreshToken returns the refresh token cookie.
func GetRefreshToken(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("refresh_token")
}

// GetSessionID returns the session ID cookie.
func GetSessionID(r *http.Request) (*http.Cookie, error) {
	return r.Cookie("session_id")
}
