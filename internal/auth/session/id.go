package session

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	id16Bytes = 16
)

// newID16 returns a new 16-byte random identifier encoded as lowercase hex.
func newID16() (string, error) {
	var b [id16Bytes]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
