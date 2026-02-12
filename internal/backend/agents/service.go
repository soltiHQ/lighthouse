package agents

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/soltiHQ/control-plane/domain/model"
	"github.com/soltiHQ/control-plane/internal/backend"
	"github.com/soltiHQ/control-plane/internal/storage"
)

// Service implements agent-related use-cases on top of storage contracts.
type Service struct {
	logger zerolog.Logger
	store  storage.AgentStore
}

// New creates a new agents service.
func New(store storage.AgentStore, logger zerolog.Logger) *Service {
	if store == nil {
		panic("agents.Service: store is nil")
	}
	return &Service{
		logger: logger.With().Str("service", "agents").Logger(),
		store:  store,
	}
}

// List returns a page of agents matching the query.
func (s *Service) List(ctx context.Context, q ListQuery) (*Page, error) {
	res, err := s.store.ListAgents(ctx, q.Filter, storage.ListOptions{
		Limit:  backend.NormalizeListLimit(q.Limit, defaultListLimit),
		Cursor: q.Cursor,
	})
	if err != nil {
		return nil, err
	}

	out := make([]View, 0, len(res.Items))
	for _, a := range res.Items {
		if a == nil {
			continue
		}
		out = append(out, toView(a))
	}
	return &Page{
		Items:      out,
		NextCursor: res.NextCursor,
	}, nil
}

// Get returns a single agent by ID.
func (s *Service) Get(ctx context.Context, id string) (View, error) {
	if id == "" {
		return View{}, storage.ErrInvalidArgument
	}

	a, err := s.store.GetAgent(ctx, id)
	if err != nil {
		return View{}, err
	}
	return toView(a), nil
}

// Upsert an agent.
func (s *Service) Upsert(ctx context.Context, model *model.Agent) error {
	return s.store.UpsertAgent(ctx, model)
}

// PatchLabels replaces labels for an agent (control-plane owned).
func (s *Service) PatchLabels(ctx context.Context, req PatchLabels) (View, error) {
	if req.ID == "" {
		return View{}, storage.ErrInvalidArgument
	}

	agent, err := s.store.GetAgent(ctx, req.ID)
	if err != nil {
		return View{}, err
	}
	if agent == nil {
		return View{}, storage.ErrInternal
	}

	replaceLabels(agent, req.Labels)

	if err = s.store.UpsertAgent(ctx, agent); err != nil {
		return View{}, err
	}
	return toView(agent), nil
}

// replaceLabels replaces the entire labels set. No merge semantics.
func replaceLabels(a *model.Agent, labels map[string]string) {
	for k := range a.LabelsAll() {
		a.LabelDelete(k)
	}
	for k, v := range labels {
		if k == "" || v == "" {
			continue
		}
		a.LabelAdd(k, v)
	}
}
