package credentials

import "github.com/soltiHQ/control-plane/domain/kind"

const defaultListLimit = 1000 // но мы не пагиним, это просто sanity-limit на будущее

// ListByUserQuery describes listing credentials for a user.
type ListByUserQuery struct {
	UserID string
}

// Page is a non-paginated list result for now (storage contract is non-paginated).
type Page struct {
	Items []View
}

// View is a read-only projection of a credential suitable for transport layers.
type View struct {
	Auth kind.Auth

	ID     string
	UserID string
}

// Delete describes a credential deletion request.
type Delete struct {
	ID string
}

// SetPassword sets/replaces password verification material for a user.
//
// Semantics:
//   - If CredentialID is empty, the service will try to find an existing password credential via
//     GetCredentialByUserAndAuth(userID, kind.Password). If missing -> ErrInvalidArgument.
//   - Verifier is replaced atomically from the application's POV:
//     DeleteVerifierByCredential(credID) then UpsertVerifier(newVerifier)
//     (DeleteVerifierByCredential is idempotent by contract).
type SetPassword struct {
	Cost int

	Password     string
	VerifierID   string
	CredentialID string
	UserID       string
}
