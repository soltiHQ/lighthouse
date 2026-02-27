package kind

// TaskKindType represents the execution backend for a task.
type TaskKindType string

const (
	TaskKindSubprocess TaskKindType = "subprocess"
	TaskKindWasm       TaskKindType = "wasm"
	TaskKindContainer  TaskKindType = "container"
)
