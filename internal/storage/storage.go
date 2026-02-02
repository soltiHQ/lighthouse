// Package storage defines persistence interfaces for control-plane domain objects.
package storage

import (
	"context"

	"github.com/soltiHQ/control-plane/domain"
)

// AgentListResult contains a page of agent results with pagination support.
type AgentListResult = ListResult[*domain.AgentModel]

// UserListResult contains a page of user results with pagination support.
type UserListResult = ListResult[*domain.UserModel]

// AgentStore defines persistence operations for agent domain objects.
type AgentStore interface {
	// UpsertAgent creates a new agent or replaces an existing one.
	//
	// If an agent with the same ID exists, it is fully replaced.
	// Otherwise, a new agent record is created.
	//
	// Returns:
	//   - ErrInvalidArgument if the agent violates storage-level invariants.
	//   - ErrInternal for unexpected storage failures.
	UpsertAgent(ctx context.Context, a *domain.AgentModel) error

	// GetAgent retrieves an agent by its unique identifier.
	//
	// Returns:
	//   - ErrNotFound if no agent with the given ID exists.
	//   - ErrInvalidArgument if the ID format is invalid.
	//   - ErrInternal for unexpected storage failures.
	GetAgent(ctx context.Context, id string) (*domain.AgentModel, error)

	// ListAgents retrieves agents matching the provided filter with pagination support.
	//
	// Results are ordered by (UpdatedAt DESC, ID ASC) to ensure:
	//   - Recently updated agents appear first.
	//   - Stable ordering when UpdatedAt values are identical.
	//   - Cursor-based pagination works correctly across request.
	//
	// The filter parameter is implementation-specific. Pass nil to retrieve all agents.
	// Use filter constructors from the concrete storage package (e.g., inmemory.NewFilter()).
	//
	// Pagination is cursor-based to handle large result sets safely.
	// Clients should:
	//   1. Make an initial request with an empty Cursor.
	//   2. Check AgentListResult.NextCursor.
	//   3. If non-empty, pass it as Cursor in the next request.
	//   4. Repeat until the NextCursor is empty.
	//
	// Returns:
	//   - ErrInvalidArgument if a filter type is incompatible or the cursor is malformed.
	//   - ErrInternal for unexpected storage failures.
	ListAgents(ctx context.Context, filter AgentFilter, opts ListOptions) (*AgentListResult, error)

	// DeleteAgent removes an agent by its unique identifier.
	//
	// Returns:
	//   - ErrNotFound if no agent with the given ID exists.
	//   - ErrInvalidArgument if the ID format is invalid.
	//   - ErrInternal for unexpected storage failures.
	DeleteAgent(ctx context.Context, id string) error
}

// UserStore defines persistence operations for user domain objects.
type UserStore interface {
	// UpsertUser creates a new user or replaces an existing one.
	//
	// If a user with the same ID exists, it is fully replaced.
	// Otherwise, a new user record is created.
	//
	// Returns:
	//   - ErrInvalidArgument if the user is nil or violates storage-level invariants.
	//   - ErrInternal for unexpected storage failures.
	UpsertUser(ctx context.Context, u *domain.UserModel) error

	// GetUser retrieves a user by its unique identifier.
	//
	// Returns:
	//   - ErrNotFound if no user with the given ID exists.
	//   - ErrInvalidArgument if the ID format is invalid.
	//   - ErrInternal for unexpected storage failures.
	GetUser(ctx context.Context, id string) (*domain.UserModel, error)

	// GetUserBySubject retrieves a user by their subject identifier (e.g., OIDC sub claim).
	//
	// Subject is typically populated from the identity provider's subject claim
	// and serves as a stable, unique identifier for authentication purposes.
	//
	// Returns:
	//   - ErrNotFound if no user with the given subject exists.
	//   - ErrInvalidArgument if the subject is empty.
	//   - ErrInternal for unexpected storage failures.
	GetUserBySubject(ctx context.Context, subject string) (*domain.UserModel, error)

	// ListUsers retrieves users matching the provided filter with pagination support.
	//
	// Results are ordered by (UpdatedAt DESC, ID ASC) to ensure:
	//   - Recently updated users appear first.
	//   - Stable ordering when UpdatedAt values are identical.
	//   - Cursor-based pagination works correctly across requests.
	//
	// The filter parameter is implementation-specific. Pass nil to retrieve all users.
	// Use filter constructors from the concrete storage package (e.g., inmemory.NewUserFilter()).
	//
	// Pagination is cursor-based to handle large result sets safely.
	// Clients should:
	//   1. Make an initial request with an empty Cursor.
	//   2. Check UserListResult.NextCursor.
	//   3. If non-empty, pass it as Cursor in the next request.
	//   4. Repeat until the NextCursor is empty.
	//
	// Returns:
	//   - ErrInvalidArgument if a filter type is incompatible or the cursor is malformed.
	//   - ErrInternal for unexpected storage failures.
	ListUsers(ctx context.Context, filter UserFilter, opts ListOptions) (*UserListResult, error)

	// DeleteUser removes a user by its unique identifier.
	//
	// Returns:
	//   - ErrNotFound if no user with the given ID exists.
	//   - ErrInvalidArgument if the ID format is invalid.
	//   - ErrInternal for unexpected storage failures.
	DeleteUser(ctx context.Context, id string) error
}

// Storage aggregates all domain-specific storage capabilities.
type Storage interface {
	AgentStore
	UserStore
}
