package inmemory

import (
	"github.com/soltiHQ/control-plane/domain"
	"github.com/soltiHQ/control-plane/internal/storage"
)

// Compile-time checks that filters implement their respective marker interfaces.
var (
	_ storage.AgentFilter = (*Filter)(nil)
	_ storage.UserFilter  = (*UserFilter)(nil)
)

// Filter provides predicate-based filtering for in-memory agent queries.
//
// Filters are composed by chaining builder methods, each adding a predicate
// that must be satisfied for an agent to match. All predicates are ANDed together.
type Filter struct {
	predicates []func(*domain.AgentModel) bool
}

// NewFilter creates a new empty filter that matches all agents.
func NewFilter() *Filter {
	return &Filter{
		predicates: make([]func(*domain.AgentModel) bool, 0),
	}
}

// ByPlatform adds a predicate matching agents on the specified platform.
func (f *Filter) ByPlatform(platform string) *Filter {
	f.predicates = append(f.predicates, func(a *domain.AgentModel) bool {
		return a.Platform() == platform
	})
	return f
}

// ByLabel adds a predicate matching agents with a specific label key-value pair.
func (f *Filter) ByLabel(key, value string) *Filter {
	f.predicates = append(f.predicates, func(a *domain.AgentModel) bool {
		v, ok := a.Label(key)
		return ok && v == value
	})
	return f
}

// ByOS adds a predicate matching agents running the specified operating system.
func (f *Filter) ByOS(os string) *Filter {
	f.predicates = append(f.predicates, func(a *domain.AgentModel) bool {
		return a.OS() == os
	})
	return f
}

// ByArch adds a predicate matching agents with the specified architecture.
func (f *Filter) ByArch(arch string) *Filter {
	f.predicates = append(f.predicates, func(a *domain.AgentModel) bool {
		return a.Arch() == arch
	})
	return f
}

// Matches evaluate whether an agent satisfies all predicates in this filter.
//
// Returns true if all predicates pass, false if any predicate fails.
// Empty filters (no predicates) match all agents.
func (f *Filter) Matches(a *domain.AgentModel) bool {
	for _, pred := range f.predicates {
		if !pred(a) {
			return false
		}
	}
	return true
}

// IsAgentFilter implements the storage.AgentFilter marker interface.
func (f *Filter) IsAgentFilter() {}

// UserFilter provides predicate-based filtering for in-memory user queries.
//
// Filters are composed by chaining builder methods, each adding a predicate
// that must be satisfied for a user to match. All predicates are ANDed together.
type UserFilter struct {
	predicates []func(*domain.UserModel) bool
}

// NewUserFilter creates a new empty filter that matches all users.
func NewUserFilter() *UserFilter {
	return &UserFilter{
		predicates: make([]func(*domain.UserModel) bool, 0),
	}
}

// ByEmail adds a predicate matching users with the specified email.
func (f *UserFilter) ByEmail(email string) *UserFilter {
	f.predicates = append(f.predicates, func(u *domain.UserModel) bool {
		return u.Email() == email
	})
	return f
}

// ByDisabled adds a predicate matching users based on their disabled status.
func (f *UserFilter) ByDisabled(disabled bool) *UserFilter {
	f.predicates = append(f.predicates, func(u *domain.UserModel) bool {
		return u.Disabled() == disabled
	})
	return f
}

// ByRole adds a predicate matching users who have the specified role.
func (f *UserFilter) ByRole(role domain.Role) *UserFilter {
	f.predicates = append(f.predicates, func(u *domain.UserModel) bool {
		for _, r := range u.Roles() {
			if r == role {
				return true
			}
		}
		return false
	})
	return f
}

// ByScope adds a predicate matching users who have the specified scope.
func (f *UserFilter) ByScope(scope domain.Scope) *UserFilter {
	f.predicates = append(f.predicates, func(u *domain.UserModel) bool {
		for _, s := range u.Scopes() {
			if s == scope {
				return true
			}
		}
		return false
	})
	return f
}

// Matches evaluate whether a user satisfies all predicates in this filter.
//
// Returns true if all predicates pass, false if any predicate fails.
// Empty filters (no predicates) match all users.
func (f *UserFilter) Matches(u *domain.UserModel) bool {
	for _, pred := range f.predicates {
		if !pred(u) {
			return false
		}
	}
	return true
}

// IsUserFilter implements the storage.UserFilter marker interface.
func (f *UserFilter) IsUserFilter() {}
