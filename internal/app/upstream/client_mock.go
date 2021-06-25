package upstream

import (
	"context"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/xds-relay/internal/pkg/log"
	"github.com/uber-go/tally"
)

// NewMockClientV3 creates a mock implementation for testing
func NewMockClientV3(
	ctx context.Context,
	ldsClient listenerservice.ListenerDiscoveryServiceClient,
	rdsClient routeservice.RouteDiscoveryServiceClient,
	edsClient endpointservice.EndpointDiscoveryServiceClient,
	cdsClient clusterservice.ClusterDiscoveryServiceClient,
	callOptions CallOptions,
	scope tally.Scope) Client {
	return &client{
		ldsClientV3: ldsClient,
		rdsClientV3: rdsClient,
		edsClientV3: edsClient,
		cdsClientV3: cdsClient,
		callOptions: callOptions,
		logger:      log.MockLogger,
		scope:       scope,
		shutdown:    make(<-chan struct{}),
	}
}

// NewMockClientEDS creates a mock implementation for testing for v3 eds
func NewMockClientEDS(
	ctx context.Context,
	edsClientV3 endpointservice.EndpointDiscoveryServiceClient,
	callOptions CallOptions,
	scope tally.Scope) Client {
	return &client{
		edsClientV3: edsClientV3,
		callOptions: callOptions,
		logger:      log.MockLogger,
		scope:       scope,
		shutdown:    make(<-chan struct{}),
	}
}

// NewMockV3 creates a mock client implementation for testing
func NewMockV3(
	ctx context.Context,
	callOptions CallOptions,
	errorOnCreate []error,
	ldsReceiveChan chan *discoveryv3.DiscoveryResponse,
	rdsReceiveChan chan *discoveryv3.DiscoveryResponse,
	edsReceiveChan chan *discoveryv3.DiscoveryResponse,
	cdsReceiveChan chan *discoveryv3.DiscoveryResponse,
	sendCb func(m interface{}) error,
	scope tally.Scope) Client {
	return NewMockClientV3(
		ctx,
		createMockLdsClientV3(errorOnCreate, ldsReceiveChan, sendCb),
		createMockRdsClientV3(errorOnCreate, rdsReceiveChan, sendCb),
		createMockEdsClientV3(errorOnCreate, edsReceiveChan, sendCb),
		createMockCdsClientV3(errorOnCreate, cdsReceiveChan, sendCb),
		callOptions,
		scope,
	)
}

// NewMockEDS creates a mock client implementation for testing v3 eds
func NewMockEDS(
	ctx context.Context,
	callOptions CallOptions,
	errorOnCreate []error,
	edsReceiveChanV3 chan *discoveryv3.DiscoveryResponse,
	sendCb func(m interface{}) error,
	scope tally.Scope) Client {
	return NewMockClientEDS(
		ctx,
		createMockEdsClientV3(errorOnCreate, edsReceiveChanV3, sendCb),
		callOptions,
		scope,
	)
}

func createMockLdsClientV3(
	errorOnCreate []error,
	receiveChan chan *discoveryv3.DiscoveryResponse,
	sendCb func(m interface{}) error) listenerservice.ListenerDiscoveryServiceClient {
	return &mockClientV3{errorOnStreamCreate: errorOnCreate, receiveChan: receiveChan, sendCb: sendCb}
}

func createMockCdsClientV3(
	errorOnCreate []error,
	receiveChan chan *discoveryv3.DiscoveryResponse,
	sendCb func(m interface{}) error) clusterservice.ClusterDiscoveryServiceClient {
	return &mockClientV3{errorOnStreamCreate: errorOnCreate, receiveChan: receiveChan, sendCb: sendCb}
}

func createMockRdsClientV3(
	errorOnCreate []error,
	receiveChan chan *discoveryv3.DiscoveryResponse,
	sendCb func(m interface{}) error) routeservice.RouteDiscoveryServiceClient {
	return &mockClientV3{errorOnStreamCreate: errorOnCreate, receiveChan: receiveChan, sendCb: sendCb}
}

func createMockEdsClientV3(
	errorOnCreate []error,
	receiveChan chan *discoveryv3.DiscoveryResponse,
	sendCb func(m interface{}) error) endpointservice.EndpointDiscoveryServiceClient {
	return &mockClientV3{errorOnStreamCreate: errorOnCreate, receiveChan: receiveChan, sendCb: sendCb}
}
