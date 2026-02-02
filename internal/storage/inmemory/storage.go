package inmemory

import (
	"context"

	"github.com/soltiHQ/control-plane/domain"
	"github.com/soltiHQ/control-plane/internal/storage"
)

// Compile-time checks that Store implements the required interfaces.
var (
	_ storage.Storage    = (*Store)(nil)
	_ storage.AgentStore = (*Store)(nil)
	_ storage.UserStore  = (*Store)(nil)
)

// Store provides an in-memory implementation of storage.Storage using GenericStore.
type Store struct {
	agents *GenericStore[*domain.AgentModel]
	users  *GenericStore[*domain.UserModel]
}

// New creates a new in-memory store with an empty state.
func New() *Store {
	return &Store{
		agents: NewGenericStore[*domain.AgentModel](),
		users:  NewGenericStore[*domain.UserModel](),
	}
}

// UpsertAgent inserts or fully replaces an agent.
//
// Delegates to GenericStore, which handles cloning and validation.
// Returns storage.ErrInvalidArgument if the agent is nil or has an empty ID.
func (s *Store) UpsertAgent(ctx context.Context, a *domain.AgentModel) error {
	if a == nil {
		return storage.ErrInvalidArgument
	}
	return s.agents.Upsert(ctx, a)
}

// GetAgent retrieves an agent by ID.
//
// Returns a deep clone to prevent external mutations affecting the stored state.
// Returns storage.ErrNotFound if no agent exists, storage.ErrInvalidArgument for empty IDs.
func (s *Store) GetAgent(ctx context.Context, id string) (*domain.AgentModel, error) {
	return s.agents.Get(ctx, id)
}

// ListAgents retrieves agents with filtering and cursor-based pagination.
//
// Filtering:
//   - Pass nil filter to retrieve all agents.
//   - Pass *inmemory.Filter created via NewFilter() for predicate-based filtering.
//   - Passing filters from other storage implementations returns storage.ErrInvalidArgument.
//
// Pagination:
//   - Results are ordered by (UpdatedAt DESC, ID ASC) for stable cursor navigation.
//   - Cursor is an opaque base64-encoded token containing position information.
//   - Invalid or corrupted cursors return storage.ErrInvalidArgument.
//
// All returned agents are deep clones isolated from the internal state.
func (s *Store) ListAgents(ctx context.Context, filter storage.AgentFilter, opts storage.ListOptions) (*storage.AgentListResult, error) {
	var predicate func(*domain.AgentModel) bool

	if filter != nil {
		f, ok := filter.(*Filter)
		if !ok {
			return nil, storage.ErrInvalidArgument
		}
		predicate = f.Matches
	}

	result, err := s.agents.List(ctx, predicate, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteAgent removes an agent by ID.
//
// Returns storage.ErrNotFound if the agent doesn't exist, storage.ErrInvalidArgument for empty IDs.
func (s *Store) DeleteAgent(ctx context.Context, id string) error {
	return s.agents.Delete(ctx, id)
}

// UpsertUser inserts or fully replaces a user.
//
// Delegates to GenericStore, which handles cloning and validation.
// Returns storage.ErrInvalidArgument if the user is nil or has an empty ID.
func (s *Store) UpsertUser(ctx context.Context, u *domain.UserModel) error {
	if u == nil {
		return storage.ErrInvalidArgument
	}
	return s.users.Upsert(ctx, u)
}

// GetUser retrieves a user by ID.
//
// Returns a deep clone to prevent external mutations affecting the stored state.
// Returns storage.ErrNotFound if no user exists, storage.ErrInvalidArgument for empty IDs.
func (s *Store) GetUser(ctx context.Context, id string) (*domain.UserModel, error) {
	return s.users.Get(ctx, id)
}

// GetUserBySubject retrieves a user by their subject identifier.
//
// This method performs a linear scan and is O(n) - acceptable for in-memory implementation
// with small datasets. Production implementations should use indexed lookups.
//
// Returns storage.ErrNotFound if no user with the subject exists, storage.ErrInvalidArgument for empty subject.
func (s *Store) GetUserBySubject(ctx context.Context, subject string) (*domain.UserModel, error) {
	if subject == "" {
		return nil, storage.ErrInvalidArgument
	}

	result, err := s.users.List(ctx, func(u *domain.UserModel) bool {
		return u.Subject() == subject
	}, storage.ListOptions{Limit: 1})
	if err != nil {
		return nil, err
	}
	if len(result.Items) == 0 {
		return nil, storage.ErrNotFound
	}
	return result.Items[0], nil
}

// ListUsers retrieves users with filtering and cursor-based pagination.
//
// Filtering:
//   - Pass nil filter to retrieve all users.
//   - Pass *inmemory.UserFilter created via NewUserFilter() for predicate-based filtering.
//   - Passing filters from other storage implementations returns storage.ErrInvalidArgument.
//
// Pagination:
//   - Results are ordered by (UpdatedAt DESC, ID ASC) for stable cursor navigation.
//   - Cursor is an opaque base64-encoded token containing position information.
//   - Invalid or corrupted cursors return storage.ErrInvalidArgument.
//
// All returned users are deep clones isolated from the internal state.
func (s *Store) ListUsers(ctx context.Context, filter storage.UserFilter, opts storage.ListOptions) (*storage.UserListResult, error) {
	var predicate func(*domain.UserModel) bool

	if filter != nil {
		f, ok := filter.(*UserFilter)
		if !ok {
			return nil, storage.ErrInvalidArgument
		}
		predicate = f.Matches
	}
	result, err := s.users.List(ctx, predicate, opts)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteUser removes a user by ID.
//
// Returns storage.ErrNotFound if the user doesn't exist, storage.ErrInvalidArgument for empty IDs.
func (s *Store) DeleteUser(ctx context.Context, id string) error {
	return s.users.Delete(ctx, id)
}
