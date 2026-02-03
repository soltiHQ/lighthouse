package credentials

import (
	"errors"

	"github.com/soltiHQ/control-plane/domain"
	
	"golang.org/x/crypto/bcrypt"
)

const (
	// PasswordHashKey is the key used to store bcrypt password hash in credential data.
	PasswordHashKey = "hash"
	// DefaultBcryptCost is the computational cost for bcrypt hashing (2^12 iterations).
	DefaultBcryptCost = 12
)

var (
	// ErrWrongCredentialType is returned when attempting password operations on non-password credentials.
	ErrWrongCredentialType = errors.New("credentials: wrong credential type")
	// ErrMissingPasswordHash is returned when a password credential lacks the hash field.
	ErrMissingPasswordHash = errors.New("credentials: missing password hash")
	// ErrPasswordMismatch is returned when the provided password doesn't match the stored hash.
	ErrPasswordMismatch = errors.New("credentials: password mismatch")
)

// NewPasswordCredential creates a new password credential with a bcrypt-hashed password.
func NewPasswordCredential(id, userID, plainPassword string) (*domain.CredentialModel, error) {
	cred, err := domain.NewCredentialModel(id, userID, domain.CredentialTypePassword)
	if err != nil {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), DefaultBcryptCost)
	if err != nil {
		return nil, err
	}
	cred.SetData(PasswordHashKey, string(hash))
	return cred, nil
}

// VerifyPassword checks if the provided plaintext password matches the stored hash.
func VerifyPassword(cred *domain.CredentialModel, plainPassword string) error {
	if cred.Type() != domain.CredentialTypePassword {
		return ErrWrongCredentialType
	}
	hash, ok := cred.GetData(PasswordHashKey)
	if !ok {
		return ErrMissingPasswordHash
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordMismatch
		}
		return err
	}
	return nil
}
