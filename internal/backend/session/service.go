package sessions

import (
	"context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/soltiHQ/control-plane/internal/backend"
	"github.com/soltiHQ/control-plane/internal/storage"
)

// Service implements session-related use-cases on top of storage contracts.
type Service struct {
	logger zerolog.Logger
	store  storage.SessionStore
}

// New creates a new sessions service.
func New(store storage.SessionStore, logger zerolog.Logger) *Service {
	if store == nil {
		panic("sessions.Service: store is nil")
	}
	return &Service{
		logger: logger.With().Str("service", "sessions").Logger(),
		store:  store,
	}
}

// Get returns a single session by ID.
func (s *Service) Get(ctx context.Context, id string) (View, error) {
	if id == "" {
		return View{}, storage.ErrInvalidArgument
	}
	sess, err := s.store.GetSession(ctx, id)
	if err != nil {
		return View{}, err
	}
	return toView(sess), nil
}

// ListByUser returns all sessions for a user.
func (s *Service) ListByUser(ctx context.Context, q ListByUserQuery) (*Page, error) {
	if q.UserID == "" {
		return nil, storage.ErrInvalidArgument
	}

	items, err := s.store.ListSessionsByUser(ctx, q.UserID)
	if err != nil {
		return nil, err
	}

	limit := backend.NormalizeListLimit(q.Limit, defaultListLimit)
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}

	out := make([]View, 0, len(items))
	for _, sess := range items {
		if sess == nil {
			continue
		}
		out = append(out, toView(sess))
	}
	return &Page{Items: out}, nil
}

// Delete deletes a single session by ID.
func (s *Service) Delete(ctx context.Context, req Delete) error {
	if req.ID == "" {
		return storage.ErrInvalidArgument
	}
	err := s.store.DeleteSession(ctx, req.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return err
	}
	return nil
}

// DeleteByUser deletes all sessions for a user.
func (s *Service) DeleteByUser(ctx context.Context, req DeleteByUser) error {
	if req.UserID == "" {
		return storage.ErrInvalidArgument
	}
	return s.store.DeleteSessionsByUser(ctx, req.UserID)
}

// Revoke marks a session as revoked (idempotent).
func (s *Service) Revoke(ctx context.Context, req Revoke) error {
	if req.ID == "" || req.At.IsZero() {
		return storage.ErrInvalidArgument
	}

	err := s.store.RevokeSession(ctx, req.ID, req.At)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return err
	}
	return nil
}
