// Package proxy provides outbound communication with agents.
//
// The control-plane calls INTO agents to query tasks, submit work, etc.
// Each agent exposes either an HTTP or a gRPC endpoint; the proxy package
// selects the transport based on the endpoint type reported by the agent.
package proxy

import (
	"context"
	"strings"

	"github.com/soltiHQ/control-plane/domain/kind"
)

// Task is the proxy-internal representation of an agent task.
type Task struct {
	ID        string
	Slot      string
	Status    string
	Attempt   int
	CreatedAt int64
	UpdatedAt int64
	Error     string
}

// TaskFilter holds optional filters and pagination for listing tasks.
type TaskFilter struct {
	Slot   string
	Status string
	Limit  int
	Offset int
}

// TaskListResult is the result of a ListTasks call.
type TaskListResult struct {
	Tasks []Task
	Total int
}

// AgentProxy is the interface for outbound communication with an agent.
type AgentProxy interface {
	ListTasks(ctx context.Context, filter TaskFilter) (*TaskListResult, error)
}

// New creates an HTTP or gRPC proxy based on the agent's endpoint type.
func New(endpoint string, epType kind.EndpointType) (AgentProxy, error) {
	switch epType {
	case kind.EndpointHTTP:
		return &httpProxy{endpoint: strings.TrimRight(endpoint, "/")}, nil
	case kind.EndpointGRPC:
		return &grpcProxy{endpoint: endpoint}, nil
	default:
		return nil, ErrUnsupportedEndpointType
	}
}
