package domain

type Role string
type Scope string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

const (
	ScopeAgentsRead  Scope = "agents:read"
	ScopeAgentsWrite Scope = "agents:write"
	ScopeUsersRead   Scope = "users:read"
)

// HasRole checks whether a user has a role.
func (u *UserModel) HasRole(role Role) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasScope checks whether a user has scope.
func (u *UserModel) HasScope(scope Scope) bool {
	for _, s := range u.scopes {
		if s == scope {
			return true
		}
	}
	return false
}
