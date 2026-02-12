package sessions

import (
	"time"

	"github.com/soltiHQ/control-plane/domain/kind"
)

const defaultListLimit = 200

// ListByUserQuery describes listing sessions for a user.
// Storage contract is non-paginated, but we keep Limit to prevent footguns.
// Cursor intentionally omitted.
type ListByUserQuery struct {
	UserID string
	Limit  int
}

// Page is a list result.
type Page struct {
	Items []View
}

// View is a read-only projection of a session suitable for transport layers.
type View struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	RevokedAt time.Time

	AuthKind     kind.Auth
	ID           string
	UserID       string
	CredentialID string

	Revoked bool
}

// Delete describes deleting a single session by ID.
type Delete struct {
	ID string
}

// DeleteByUser describes deleting all sessions for a user.
type DeleteByUser struct {
	UserID string
}

// Revoke describes revoking a session by ID.
type Revoke struct {
	At time.Time
	ID string
}
