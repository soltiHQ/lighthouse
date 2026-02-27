package restv1

// RolloutSpec embeds Spec with per-agent delivery state.
type RolloutSpec struct {
	Spec
	Entries []RolloutEntry `json:"rollout,omitempty"`
}

// RolloutEntry tracks the delivery state of a spec on a single agent.
type RolloutEntry struct {
	AgentID      string `json:"agent_id"`
	Status       string `json:"status"`
	LastPushedAt string `json:"last_pushed_at,omitempty"`
	LastSyncedAt string `json:"last_synced_at,omitempty"`
	Error        string `json:"error,omitempty"`

	DesiredVersion int `json:"desired_version"`
	ActualVersion  int `json:"actual_version"`
	Attempts       int `json:"attempts,omitempty"`
}
