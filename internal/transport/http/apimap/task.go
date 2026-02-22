package apimap

import (
	proxyv1 "github.com/soltiHQ/control-plane/api/proxy/v1"
	"github.com/soltiHQ/control-plane/internal/proxy"
)

// TaskFromProxy converts a proxy.Task to a proxyv1.Task.
func TaskFromProxy(t proxy.Task) proxyv1.Task {
	return proxyv1.Task{
		ID:        t.ID,
		Slot:      t.Slot,
		Status:    t.Status,
		Attempt:   t.Attempt,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Error:     t.Error,
	}
}

// TasksFromProxy converts a slice of proxy.Task to a slice of proxyv1.Task.
func TasksFromProxy(tasks []proxy.Task) []proxyv1.Task {
	out := make([]proxyv1.Task, len(tasks))
	for i, t := range tasks {
		out[i] = TaskFromProxy(t)
	}
	return out
}
