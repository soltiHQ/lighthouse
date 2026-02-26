package taskspec

import (
	"context"

	"github.com/soltiHQ/control-plane/domain/model"
	"github.com/soltiHQ/control-plane/internal/service"
	"github.com/soltiHQ/control-plane/internal/storage"
)

// Service provides task spec management operations.
type Service struct {
	store storage.Storage
}

// New creates a new task spec service.
func New(store storage.Storage) *Service {
	if store == nil {
		panic("taskspec.Service: store is nil")
	}
	return &Service{store: store}
}

// List returns a page of task specs matching the query.
func (s *Service) List(ctx context.Context, q ListQuery) (*Page, error) {
	res, err := s.store.ListSpecs(ctx, q.Filter, storage.ListOptions{
		Limit:  service.NormalizeListLimit(q.Limit, defaultListLimit),
		Cursor: q.Cursor,
	})
	if err != nil {
		return nil, err
	}

	out := make([]*model.Spec, 0, len(res.Items))
	for _, ts := range res.Items {
		if ts == nil {
			continue
		}
		out = append(out, ts.Clone())
	}
	return &Page{
		Items:      out,
		NextCursor: res.NextCursor,
	}, nil
}

// Get returns a single task spec by ID.
func (s *Service) Get(ctx context.Context, id string) (*model.Spec, error) {
	if id == "" {
		return nil, storage.ErrInvalidArgument
	}
	ts, err := s.store.GetSpec(ctx, id)
	if err != nil {
		return nil, err
	}
	return ts.Clone(), nil
}

// Create persists a new task spec.
func (s *Service) Create(ctx context.Context, ts *model.Spec) error {
	if ts == nil {
		return storage.ErrInvalidArgument
	}
	return s.store.UpsertSpec(ctx, ts)
}

// Update persists changes to an existing task spec and increments its version.
func (s *Service) Update(ctx context.Context, ts *model.Spec) error {
	if ts == nil {
		return storage.ErrInvalidArgument
	}
	// Verify it exists
	if _, err := s.store.GetSpec(ctx, ts.ID()); err != nil {
		return err
	}
	ts.IncrementVersion()
	return s.store.UpsertSpec(ctx, ts)
}

// Delete removes a task spec and all associated rollouts.
func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return storage.ErrInvalidArgument
	}
	if err := s.store.DeleteRolloutsBySpec(ctx, id); err != nil {
		return err
	}
	return s.store.DeleteSpec(ctx, id)
}

// Deploy creates or updates rollout records for all target agents with status pending.
func (s *Service) Deploy(ctx context.Context, specID string) error {
	ts, err := s.store.GetSpec(ctx, specID)
	if err != nil {
		return err
	}

	for _, agentID := range ts.Targets() {
		ssID := "ss-" + specID + "-" + agentID
		existing, err := s.store.GetRollout(ctx, ssID)
		if err == nil {
			// Update existing rollout
			existing.MarkPending(ts.Version())
			if err := s.store.UpsertRollout(ctx, existing); err != nil {
				return err
			}
			continue
		}
		// Create new rollout
		ss, err := model.NewSyncState(specID, agentID, ts.Version())
		if err != nil {
			return err
		}
		if err := s.store.UpsertRollout(ctx, ss); err != nil {
			return err
		}
	}
	return nil
}

// RolloutsBySpec returns all rollouts for a given spec.
func (s *Service) RolloutsBySpec(ctx context.Context, specID string) ([]*model.SyncState, error) {
	if specID == "" {
		return nil, storage.ErrInvalidArgument
	}

	// Use a filter to match by spec ID
	res, err := s.store.ListRollouts(ctx, nil, storage.ListOptions{Limit: storage.MaxListLimit})
	if err != nil {
		return nil, err
	}

	out := make([]*model.SyncState, 0)
	for _, ss := range res.Items {
		if ss.SpecID() == specID {
			out = append(out, ss.Clone())
		}
	}
	return out, nil
}
