package jwt

import jwtlib "github.com/golang-jwt/jwt/v5"

// Claims is a typed set of JWT claims used by Control Plane.
type Claims struct {
	jwtlib.RegisteredClaims

	UserID      string   `json:"uid,omitempty"`
	Permissions []string `json:"perms,omitempty"`
}
