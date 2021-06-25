package transport

import (
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/envoyproxy/xds-relay/internal/pkg/stats"
)

const (
	nodeID          = "1"
	resourceName    = "route"
	resourceVersion = "version"
	cluster         = "cluster"
	region          = "region"
	zone            = "zone"
	subzone         = "subzone"
)

var requestV3 = discoveryv3.DiscoveryRequest{
	Node: &corev3.Node{
		Id: nodeID,
		Metadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{"a": nil},
		},
		Cluster: cluster,
		Locality: &corev3.Locality{
			Region:  region,
			Zone:    zone,
			SubZone: subzone,
		},
	},
	ResourceNames: []string{resourceName},
	TypeUrl:       "typeUrl",
	VersionInfo:   resourceVersion,
	ErrorDetail:   &status.Status{Code: 0},
	ResponseNonce: "1",
}

func TestGetRaw(t *testing.T) {
	requestv3 := NewRequestV3(&requestV3)
	assert.Equal(t, requestv3.GetRaw().V3, &requestV3)
}

func TestCreateWatch(t *testing.T) {
	requestv3 := NewRequestV3(&requestV3)
	scope := stats.NewMockScope("mockwatch")
	assert.NotNil(t, requestv3.CreateWatch(scope).GetChannel().V3)
}
