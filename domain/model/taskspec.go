package model

import (
	"time"

	"github.com/soltiHQ/control-plane/domain"
	"github.com/soltiHQ/control-plane/domain/kind"
)

var _ domain.Entity[*Spec] = (*Spec)(nil)

// BackoffConfig holds backoff parameters for task restart delays.
type BackoffConfig struct {
	Jitter  kind.JitterStrategy
	FirstMs int64
	MaxMs   int64
	Factor  float64
}

// Spec represents a desired task specification managed by the control-plane.
//
// A Spec defines what task should run on which agents. It is the "desired state"
// in the reconciliation model — the sync runner compares it against what agents actually have.
//
// The spec fields mirror the agent's CreateSpec format: slot, kind, timeout, restart,
// backoff, admission, and runner labels.
type Spec struct {
	// CP metadata
	id           string
	name         string
	version      int
	targets      []string          // concrete agent IDs
	targetLabels map[string]string // label selector for dynamic targeting
	createdAt    time.Time
	updatedAt    time.Time

	// Spec (mirrors agent CreateSpec)
	slot         string
	kindType     kind.TaskKindType
	kindConfig   map[string]any // e.g. {command, args, env, cwd, failOnNonZero} for subprocess
	timeoutMs    int64
	restartType  kind.RestartType
	intervalMs   int64 // only for RestartAlways
	backoff      BackoffConfig
	admission    kind.AdmissionStrategy
	runnerLabels map[string]string
}

// NewSpec creates a new Spec domain entity with sensible defaults.
func NewSpec(id, name, slot string) (*Spec, error) {
	if id == "" {
		return nil, domain.ErrEmptyID
	}
	if slot == "" {
		return nil, domain.ErrFieldEmpty
	}
	now := time.Now()
	return &Spec{
		createdAt: now,
		updatedAt: now,

		id:      id,
		name:    name,
		slot:    slot,
		version: 1,

		targets:      nil,
		targetLabels: make(map[string]string),
		kindType:     kind.TaskKindSubprocess,
		kindConfig:   make(map[string]any),
		timeoutMs:    30000,
		restartType:  kind.RestartNever,
		intervalMs:   0,
		backoff: BackoffConfig{
			Jitter:  kind.JitterNone,
			FirstMs: 1000,
			MaxMs:   5000,
			Factor:  2.0,
		},
		admission:    kind.AdmissionDropIfRunning,
		runnerLabels: make(map[string]string),
	}, nil
}

// --- Getters ---

func (ts *Spec) ID() string            { return ts.id }
func (ts *Spec) Name() string           { return ts.name }
func (ts *Spec) Slot() string           { return ts.slot }
func (ts *Spec) Version() int           { return ts.version }
func (ts *Spec) CreatedAt() time.Time   { return ts.createdAt }
func (ts *Spec) UpdatedAt() time.Time   { return ts.updatedAt }
func (ts *Spec) KindType() kind.TaskKindType       { return ts.kindType }
func (ts *Spec) TimeoutMs() int64                   { return ts.timeoutMs }
func (ts *Spec) RestartType() kind.RestartType      { return ts.restartType }
func (ts *Spec) IntervalMs() int64                  { return ts.intervalMs }
func (ts *Spec) Backoff() BackoffConfig             { return ts.backoff }
func (ts *Spec) Admission() kind.AdmissionStrategy  { return ts.admission }

// KindConfig returns a defensive copy of the kind configuration.
func (ts *Spec) KindConfig() map[string]any {
	out := make(map[string]any, len(ts.kindConfig))
	for k, v := range ts.kindConfig {
		out[k] = v
	}
	return out
}

// Targets returns a copy of the target agent IDs.
func (ts *Spec) Targets() []string {
	out := make([]string, len(ts.targets))
	copy(out, ts.targets)
	return out
}

// TargetLabels returns a defensive copy of the target label selector.
func (ts *Spec) TargetLabels() map[string]string {
	out := make(map[string]string, len(ts.targetLabels))
	for k, v := range ts.targetLabels {
		out[k] = v
	}
	return out
}

// RunnerLabels returns a defensive copy of the runner labels.
func (ts *Spec) RunnerLabels() map[string]string {
	out := make(map[string]string, len(ts.runnerLabels))
	for k, v := range ts.runnerLabels {
		out[k] = v
	}
	return out
}

// --- Setters ---

func (ts *Spec) SetName(name string) {
	ts.name = name
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetSlot(slot string) {
	ts.slot = slot
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetKindType(kt kind.TaskKindType) {
	ts.kindType = kt
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetKindConfig(cfg map[string]any) {
	cp := make(map[string]any, len(cfg))
	for k, v := range cfg {
		cp[k] = v
	}
	ts.kindConfig = cp
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetTimeoutMs(ms int64) {
	ts.timeoutMs = ms
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetRestartType(rt kind.RestartType) {
	ts.restartType = rt
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetIntervalMs(ms int64) {
	ts.intervalMs = ms
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetBackoff(b BackoffConfig) {
	ts.backoff = b
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetAdmission(a kind.AdmissionStrategy) {
	ts.admission = a
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetTargets(targets []string) {
	cp := make([]string, len(targets))
	copy(cp, targets)
	ts.targets = cp
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetTargetLabels(labels map[string]string) {
	cp := make(map[string]string, len(labels))
	for k, v := range labels {
		cp[k] = v
	}
	ts.targetLabels = cp
	ts.updatedAt = time.Now()
}

func (ts *Spec) SetRunnerLabels(labels map[string]string) {
	cp := make(map[string]string, len(labels))
	for k, v := range labels {
		cp[k] = v
	}
	ts.runnerLabels = cp
	ts.updatedAt = time.Now()
}

// IncrementVersion bumps the version number and updates the timestamp.
func (ts *Spec) IncrementVersion() {
	ts.version++
	ts.updatedAt = time.Now()
}

// ToCreateSpec builds a map[string]any in the agent's CreateSpec JSON format.
//
// Example output:
//
//	{"slot":"worker","kind":{"subprocess":{"command":"sleep","args":["30"]}},"timeoutMs":60000,
//	 "restart":{"type":"never"},"backoff":{"jitter":"none","firstMs":1000,"maxMs":5000,"factor":2.0},
//	 "admission":"dropIfRunning"}
func (ts *Spec) ToCreateSpec() map[string]any {
	// kind
	kindCfg := make(map[string]any, len(ts.kindConfig))
	for k, v := range ts.kindConfig {
		kindCfg[k] = v
	}

	// restart
	restart := map[string]any{"type": string(ts.restartType)}
	if ts.restartType == kind.RestartAlways && ts.intervalMs > 0 {
		restart["intervalMs"] = ts.intervalMs
	}

	spec := map[string]any{
		"slot":      ts.slot,
		"kind":      map[string]any{string(ts.kindType): kindCfg},
		"timeoutMs": ts.timeoutMs,
		"restart":   restart,
		"backoff": map[string]any{
			"jitter":  string(ts.backoff.Jitter),
			"firstMs": ts.backoff.FirstMs,
			"maxMs":   ts.backoff.MaxMs,
			"factor":  ts.backoff.Factor,
		},
		"admission": string(ts.admission),
	}
	if len(ts.runnerLabels) > 0 {
		labels := make(map[string]string, len(ts.runnerLabels))
		for k, v := range ts.runnerLabels {
			labels[k] = v
		}
		spec["labels"] = labels
	}
	return spec
}

// Clone creates a deep copy of the Spec.
func (ts *Spec) Clone() *Spec {
	kindConfig := make(map[string]any, len(ts.kindConfig))
	for k, v := range ts.kindConfig {
		kindConfig[k] = v
	}
	targets := make([]string, len(ts.targets))
	copy(targets, ts.targets)
	targetLabels := make(map[string]string, len(ts.targetLabels))
	for k, v := range ts.targetLabels {
		targetLabels[k] = v
	}
	runnerLabels := make(map[string]string, len(ts.runnerLabels))
	for k, v := range ts.runnerLabels {
		runnerLabels[k] = v
	}

	return &Spec{
		id:           ts.id,
		name:         ts.name,
		version:      ts.version,
		targets:      targets,
		targetLabels: targetLabels,
		createdAt:    ts.createdAt,
		updatedAt:    ts.updatedAt,

		slot:         ts.slot,
		kindType:     ts.kindType,
		kindConfig:   kindConfig,
		timeoutMs:    ts.timeoutMs,
		restartType:  ts.restartType,
		intervalMs:   ts.intervalMs,
		backoff:      ts.backoff,
		admission:    ts.admission,
		runnerLabels: runnerLabels,
	}
}
