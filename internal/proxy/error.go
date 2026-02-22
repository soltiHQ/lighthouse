package proxy

import "errors"

var (
	// ErrUnsupportedEndpointType indicates the agent reported an unknown endpoint type.
	//
	// This typically means the control-plane is out of sync with the agent's protocol version.
	// Callers should log this as an error and treat the agent as unreachable.
	ErrUnsupportedEndpointType = errors.New("proxy: unsupported endpoint type")
)
