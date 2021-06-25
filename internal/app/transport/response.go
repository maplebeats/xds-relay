package transport

import (
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes/any"
)

// ResponseVersion holds either one of the v2/v3 DiscoveryRequests
type ResponseVersion struct {
	V3 *discoveryv3.DiscoveryResponse
}

// Response is the generic response interface
type Response interface {
	GetPayloadVersion() string
	GetNonce() string
	GetTypeURL() string
	GetRequest() *RequestVersion
	Get() *ResponseVersion
	GetResources() []*any.Any
}
