package policy

import (
	"github.com/soltiHQ/control-plane/domain/kind"
	"github.com/soltiHQ/control-plane/internal/auth/identity"
)

// Convenience aliases to keep policy builders readable.
const (
	usersEdit   = kind.UsersEdit
	usersDelete = kind.UsersDelete
)

// permSet builds a lookup set from an identity's permissions.
func permSet(id *identity.Identity) map[kind.Permission]struct{} {
	m := make(map[kind.Permission]struct{}, len(id.Permissions))
	for _, p := range id.Permissions {
		m[p] = struct{}{}
	}
	return m
}

// hasAny returns true when any of wants is present in set.
func hasAny(set map[kind.Permission]struct{}, wants ...kind.Permission) bool {
	for _, w := range wants {
		if _, ok := set[w]; ok {
			return true
		}
	}
	return false
}
