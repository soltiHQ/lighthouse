package backend

import (
	"context"

	"github.com/soltiHQ/control-plane/domain"
	"github.com/soltiHQ/control-plane/internal/storage"

	"github.com/rs/zerolog"
)

// AgentList returns a list of all agents.
func AgentList(ctx context.Context, logger zerolog.Logger, store storage.Storage) ([]domain.AgentModel, error) {
	return nil, nil
}
