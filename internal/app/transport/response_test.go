package transport

import (
	"testing"

	discoveryv2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	nonce   = "nonce"
	typeURL = "typeUrl"
)

var responseV2 = discoveryv2.DiscoveryResponse{
	VersionInfo: resourceVersion,
	TypeUrl:     typeURL,
	Resources:   []*anypb.Any{},
	Nonce:       nonce,
}
var responseV3 = discoveryv3.DiscoveryResponse{
	VersionInfo: resourceVersion,
	TypeUrl:     typeURL,
	Resources:   []*anypb.Any{},
	Nonce:       nonce,
}

func TestGetPayloadVersion(t *testing.T) {
	responsev3 := NewResponseV3(&requestV3, &responseV3)
	assert.Equal(t, responsev3.GetPayloadVersion(), resourceVersion)
}

func TestGetNonce(t *testing.T) {
	responsev3 := NewResponseV3(&requestV3, &responseV3)
	assert.Equal(t, responsev3.GetNonce(), nonce)
}

func TestResponseGetTypeURL(t *testing.T) {
	responsev3 := NewResponseV3(&requestV3, &responseV3)
	assert.Equal(t, responsev3.GetTypeURL(), typeURL)
}

func TestGetRequest(t *testing.T) {
	responsev3 := NewResponseV3(&requestV3, &responseV3)
	assert.Equal(t, responsev3.GetRequest().V3, &requestV3)
}

func TestGetResponse(t *testing.T) {
	responsev3 := NewResponseV3(&requestV3, &responseV3)
	assert.Equal(t, responsev3.Get().V3, &responseV3)
}

func TestGetResources(t *testing.T) {
	responsev3 := NewResponseV3(&requestV3, &responseV3)
	assert.Equal(t, responsev3.GetResources(), responseV3.Resources)
}
