package kind

// EndpointType describes the transport protocol an agent exposes.
type EndpointType string

const (
	EndpointGRPC EndpointType = "grpc"
	EndpointHTTP EndpointType = "http"
)

// EndpointTypeFromInt maps the proto/JSON integer enum to EndpointType.
func EndpointTypeFromInt(v int) EndpointType {
	if v == 1 {
		return EndpointHTTP
	}
	return EndpointGRPC
}
