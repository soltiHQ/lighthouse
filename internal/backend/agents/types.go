package agents

import "github.com/soltiHQ/control-plane/internal/storage"

const defaultListLimit = 30

// ListQuery describes a paginated agents listing request.
type ListQuery struct {
	// Filter is a storage-level filter. Backends validate that the filter
	// was constructed for that backend and return storage.ErrInvalidArgument otherwise.
	Filter storage.AgentFilter

	Cursor string
	Limit  int
}

// Page is a paginated agents listing result.
type Page struct {
	Items      []View
	NextCursor string
}

// View is a read-only projection of an agent suitable for transport layers.
type View struct {
	UptimeSeconds int64

	Metadata map[string]string
	Labels   map[string]string

	Name     string
	Endpoint string
	OS       string
	Arch     string
	Platform string
	ID       string
}

// PatchLabels updates control-plane owned labels for an agent.
//
// Semantics:
//   - ID is required.
//   - Labels replaces the entire label set (no merge).
type PatchLabels struct {
	Labels map[string]string
	ID     string
}
