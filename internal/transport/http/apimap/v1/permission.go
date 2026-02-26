package apimapv1

import "github.com/soltiHQ/control-plane/domain/kind"

// Permission converts a domain permission to its string representation.
func Permission(p kind.Permission) string {
	return string(p)
}
