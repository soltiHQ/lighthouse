package model

import (
	"time"

	"github.com/soltiHQ/control-plane/domain"
	"github.com/soltiHQ/control-plane/domain/kind"
)

var _ domain.Entity[*SyncState] = (*SyncState)(nil)

// SyncState tracks the synchronization state of a Spec on a specific agent.
//
// It records the desired version (what CP wants) vs actual version (what the agent has),
// enabling the sync runner to detect drift and reconcile.
type SyncState struct {
	createdAt time.Time
	updatedAt time.Time

	lastPushedAt time.Time
	lastSyncedAt time.Time

	id         string
	specID string
	agentID    string
	errMsg     string

	desiredVersion int
	actualVersion  int
	attempts       int

	status kind.SyncStatus
}

// NewSyncState creates a new SyncState for a Spec-Agent pair.
func NewSyncState(specID, agentID string, desiredVersion int) (*SyncState, error) {
	if specID == "" || agentID == "" {
		return nil, domain.ErrEmptyID
	}
	now := time.Now()
	return &SyncState{
		createdAt: now,
		updatedAt: now,

		id:         "ss-" + specID + "-" + agentID,
		specID: specID,
		agentID:    agentID,

		desiredVersion: desiredVersion,
		status:         kind.SyncStatusPending,
	}, nil
}

// ID returns the sync state's unique identifier.
func (ss *SyncState) ID() string { return ss.id }

// SpecID returns the associated Spec ID.
func (ss *SyncState) SpecID() string { return ss.specID }

// AgentID returns the target agent ID.
func (ss *SyncState) AgentID() string { return ss.agentID }

// DesiredVersion returns the Spec version that should be on the agent.
func (ss *SyncState) DesiredVersion() int { return ss.desiredVersion }

// ActualVersion returns the Spec version confirmed on the agent.
func (ss *SyncState) ActualVersion() int { return ss.actualVersion }

// Status returns the current sync status.
func (ss *SyncState) Status() kind.SyncStatus { return ss.status }

// LastPushedAt returns when the spec was last pushed to the agent.
func (ss *SyncState) LastPushedAt() time.Time { return ss.lastPushedAt }

// LastSyncedAt returns when the agent last confirmed sync.
func (ss *SyncState) LastSyncedAt() time.Time { return ss.lastSyncedAt }

// Error returns the last error message (if any).
func (ss *SyncState) Error() string { return ss.errMsg }

// Attempts returns the retry counter.
func (ss *SyncState) Attempts() int { return ss.attempts }

// CreatedAt returns the creation timestamp.
func (ss *SyncState) CreatedAt() time.Time { return ss.createdAt }

// UpdatedAt returns the last modification timestamp.
func (ss *SyncState) UpdatedAt() time.Time { return ss.updatedAt }

// MarkPending sets the state to pending with a new desired version.
func (ss *SyncState) MarkPending(desiredVersion int) {
	ss.desiredVersion = desiredVersion
	ss.status = kind.SyncStatusPending
	ss.attempts = 0
	ss.errMsg = ""
	ss.updatedAt = time.Now()
}

// MarkSynced marks the agent as having the correct version.
func (ss *SyncState) MarkSynced(actualVersion int) {
	ss.actualVersion = actualVersion
	ss.status = kind.SyncStatusSynced
	ss.lastSyncedAt = time.Now()
	ss.errMsg = ""
	ss.updatedAt = time.Now()
}

// MarkDrift marks a version mismatch detected via export.
func (ss *SyncState) MarkDrift() {
	ss.status = kind.SyncStatusDrift
	ss.updatedAt = time.Now()
}

// MarkFailed records a push failure.
func (ss *SyncState) MarkFailed(errMsg string) {
	ss.status = kind.SyncStatusFailed
	ss.errMsg = errMsg
	ss.attempts++
	ss.lastPushedAt = time.Now()
	ss.updatedAt = time.Now()
}

// MarkUnknown sets the state when the agent is unreachable.
func (ss *SyncState) MarkUnknown() {
	ss.status = kind.SyncStatusUnknown
	ss.updatedAt = time.Now()
}

// SetLastPushedAt records a push attempt timestamp.
func (ss *SyncState) SetLastPushedAt(t time.Time) {
	ss.lastPushedAt = t
	ss.updatedAt = time.Now()
}

// Clone creates a deep copy of the SyncState.
func (ss *SyncState) Clone() *SyncState {
	return &SyncState{
		createdAt:    ss.createdAt,
		updatedAt:    ss.updatedAt,
		lastPushedAt: ss.lastPushedAt,
		lastSyncedAt: ss.lastSyncedAt,

		id:         ss.id,
		specID: ss.specID,
		agentID:    ss.agentID,
		errMsg:     ss.errMsg,

		desiredVersion: ss.desiredVersion,
		actualVersion:  ss.actualVersion,
		attempts:       ss.attempts,

		status: ss.status,
	}
}
