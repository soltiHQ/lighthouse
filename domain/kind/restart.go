package kind

// RestartType determines whether a task should be automatically restarted.
type RestartType string

const (
	RestartNever     RestartType = "never"
	RestartOnFailure RestartType = "onFailure"
	RestartAlways    RestartType = "always"
)
