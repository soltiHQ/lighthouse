package storage

// AgentFilterSeal is a marker type used to tag valid AgentFilter implementations.
type AgentFilterSeal struct{}

// UserFilterSeal is a marker type used to tag valid UserFilter implementations.
type UserFilterSeal struct{}

// RoleFilterSeal is a marker type used to tag valid RoleFilter implementations.
type RoleFilterSeal struct{}

// AgentFilter defines a backend-specific query object for agents.
//
// A filter must be constructed by the same storage backend that consumes it.
// Passing a filter created for a different backend must return ErrInvalidArgument.
type AgentFilter interface {
	IsAgentFilter(AgentFilterSeal)
}

// UserFilter defines a backend-specific query object for users.
//
// A filter must be constructed by the same storage backend that consumes it.
// Passing a filter created for a different backend must return ErrInvalidArgument.
type UserFilter interface {
	IsUserFilter(UserFilterSeal)
}

// RoleFilter defines a backend-specific query object for roles.
//
// A filter must be constructed by the same storage backend that consumes it.
// Passing a filter created for a different backend must return ErrInvalidArgument.
type RoleFilter interface {
	IsRoleFilter(RoleFilterSeal)
}
