package kind

// SyncStatus describes the synchronization state of a Spec on a specific agent.
type SyncStatus uint8

const (
	SyncStatusPending SyncStatus = iota // created/updated, not yet pushed
	SyncStatusSynced                    // desiredVersion == actualVersion
	SyncStatusDrift                     // export shows mismatch
	SyncStatusFailed                    // push error
	SyncStatusUnknown                   // agent unreachable
)

// String returns the human-readable sync status label.
func (s SyncStatus) String() string {
	switch s {
	case SyncStatusPending:
		return "pending"
	case SyncStatusSynced:
		return "synced"
	case SyncStatusDrift:
		return "drift"
	case SyncStatusFailed:
		return "failed"
	case SyncStatusUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}
