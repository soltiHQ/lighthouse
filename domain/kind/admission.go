package kind

// AdmissionStrategy defines how the controller admits a new task into a slot.
type AdmissionStrategy string

const (
	AdmissionDropIfRunning AdmissionStrategy = "dropIfRunning"
	AdmissionReplace       AdmissionStrategy = "replace"
	AdmissionQueue         AdmissionStrategy = "queue"
)
