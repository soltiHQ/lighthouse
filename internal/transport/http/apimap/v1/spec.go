package apimapv1

import (
	"time"

	restv1 "github.com/soltiHQ/control-plane/api/rest/v1"
	"github.com/soltiHQ/control-plane/domain/model"
)

// Spec maps a domain Spec to its REST DTO.
func Spec(ts *model.Spec) restv1.Spec {
	if ts == nil {
		return restv1.Spec{}
	}
	return restv1.Spec{
		ID:      ts.ID(),
		Name:    ts.Name(),
		Slot:    ts.Slot(),
		Version: ts.Version(),

		KindType:   string(ts.KindType()),
		KindConfig: ts.KindConfig(),

		TimeoutMs:      ts.TimeoutMs(),
		IntervalMs:     ts.IntervalMs(),
		BackoffFirstMs: ts.Backoff().FirstMs,
		BackoffMaxMs:   ts.Backoff().MaxMs,
		BackoffFactor:  ts.Backoff().Factor,

		Jitter:      string(ts.Backoff().Jitter),
		RestartType: string(ts.RestartType()),
		Admission:   string(ts.Admission()),

		Targets:      ts.Targets(),
		TargetLabels: ts.TargetLabels(),
		RunnerLabels: ts.RunnerLabels(),

		CreateSpec: ts.ToCreateSpec(),

		CreatedAt: ts.CreatedAt().Format(time.RFC3339),
		UpdatedAt: ts.UpdatedAt().Format(time.RFC3339),
	}
}
