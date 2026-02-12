package users

import (
	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/storage"
)

const defaultListLimit = 60

// ListQuery describes a paginated users listing request.
type ListQuery struct {
	// Filter is a storage-level filter. Backends validate that the filter
	// was constructed for that backend and return storage.ErrInvalidArgument otherwise.
	Filter storage.UserFilter

	Cursor string
	Limit  int
}

// Page is a paginated users listing result.
type Page struct {
	Items      []View
	NextCursor string
}

// View is a read-only projection of a user suitable for transport layers.
type View struct {
	Permissions []kind.Permission
	RoleIDs     []string

	Subject string
	Email   string
	Name    string
	ID      string

	Disabled bool
}
