package ratelimitkey

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"strings"
)

func LoginKey(r *http.Request, subject string) string {
	subject = strings.TrimSpace(strings.ToLower(subject))

	var (
		ip  = remoteIP(r)
		uah = shortHash("")
	)
	if r != nil {
		uah = shortHash(r.UserAgent())
	}
	if subject == "" {
		return "login::" + ip + ":" + uah
	}
	return "login:" + subject + ":" + ip + ":" + uah
}

func remoteIP(r *http.Request) string {
	if r == nil {
		return "unknown"
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		if parsed := net.ParseIP(host); parsed != nil {
			return parsed.String()
		}
	}
	if parsed := net.ParseIP(r.RemoteAddr); parsed != nil {
		return parsed.String()
	}
	return "unknown"
}

func shortHash(s string) string {
	if s == "" {
		return "none"
	}
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:8])
}
