package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gcp "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/xds-relay/internal/app/transport"
	"github.com/envoyproxy/xds-relay/internal/pkg/log"
	bootstrapv1 "github.com/envoyproxy/xds-relay/pkg/api/bootstrap/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestAdminServer_DefaultHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := defaultHandler([]Handler{{
		"/foo",
		"does nothing",
		http.HandlerFunc(nil),
		true,
	}})

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "admin commands are:\n  /foo: does nothing\n", rr.Body.String())
}

func TestAdminServer_DefaultHandler_NotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/not-implemented", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := defaultHandler([]Handler{{
		"/foo",
		"does nothing",
		http.HandlerFunc(nil),
		true,
	}})

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Equal(t, "404 page not found\n", rr.Body.String())
}

func TestAdminServer_ConfigDumpHandler(t *testing.T) {
	for _, url := range []string{"/server_info", "/server_info/"} {
		req, err := http.NewRequest("GET", url, nil)
		assert.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := configDumpHandler(&bootstrapv1.Bootstrap{
			Server: &bootstrapv1.Server{Address: &bootstrapv1.SocketAddress{
				Address:   "127.0.0.1",
				PortValue: 9991,
			}},
			OriginServer: nil,
			Logging:      nil,
			Cache:        nil,
		})

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t,
			`{
  "server": {
    "address": {
      "address": "127.0.0.1",
      "port_value": 9991
    }
  }
}
`,
			rr.Body.String())
	}
}

func TestGetParam(t *testing.T) {
	path := "/cache/foo_production_*"
	cacheKey := getParam(path, "/cache")
	assert.Equal(t, "foo_production_*", cacheKey)
}

func TestGetParam_Empty(t *testing.T) {
	prefix := "/cache"
	path := "/cache/"
	cacheKey := getParam(path, prefix)
	assert.Equal(t, "", cacheKey)

	path = "/cache"
	cacheKey = getParam(path, prefix)
	assert.Equal(t, "", cacheKey)
}

func TestGetBoolQuery(t *testing.T) {
	isVerbose, err := getBoolQueryValue(url.Values{
		"verbose": []string{"true"},
	}, "verbose")
	assert.True(t, isVerbose)
	assert.NoError(t, err)
}

func TestGetBoolQuery_Empty(t *testing.T) {
	queryValue, err := getBoolQueryValue(url.Values{}, "abc")
	assert.False(t, queryValue)
	assert.Error(t, err)
}

func TestGetBoolQuery_Malformed(t *testing.T) {
	isVerbose, err := getBoolQueryValue(url.Values{
		"verbose":   []string{"abc"},
		"something": []string{"true"},
	}, "verbose")
	assert.False(t, isVerbose)
	assert.Error(t, err)
}

func TestAdminServer_LogLevelHandler(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	logger := log.NewMock("error", &buf)
	assert.Equal(t, 0, buf.Len())

	logger.Error(ctx, "foo")
	logger.Debug(ctx, "bar")
	output := buf.String()
	assert.Contains(t, output, "foo")
	assert.NotContains(t, output, "bar")

	req, err := http.NewRequest("POST", logURL+"/debug", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := logLevelHandler(logger)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, rr.Body.String(), "Current log level: debug\n")
	logger.Debug(ctx, "bar")
	output = buf.String()
	assert.Contains(t, output, "bar")

	req, err = http.NewRequest("POST", logURL+"/info", nil)
	assert.NoError(t, err)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Current log level: info\n")
	logger.Debug(ctx, "baz")
	logger.Info(ctx, "qux")
	output = buf.String()
	assert.NotContains(t, output, "baz")
	assert.Contains(t, output, "qux")
}

func TestAdminServer_LogLevelHandler_GetLevel(t *testing.T) {
	for _, url := range []string{logURL, logURL} {
		var buf bytes.Buffer
		logger := log.NewMock("error", &buf)
		assert.Equal(t, 0, buf.Len())

		req, err := http.NewRequest("POST", url, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := logLevelHandler(logger)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, rr.Body.String(), "Current log level: error\n")

		req, err = http.NewRequest("POST", logURL+"/info", nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		req, err = http.NewRequest("POST", url, nil)
		assert.NoError(t, err)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Current log level: info\n")
	}
}

func TestMarshalResources(t *testing.T) {
	listener := &listenerv3.Listener{
		Name: "lds resource",
	}
	listenerAny, err := ptypes.MarshalAny(listener)
	assert.NoError(t, err)
	marshalled := marshalResources([]*any.Any{
		listenerAny,
	})
	assert.NotNil(t, marshalled)
	assert.Equal(t, 1, len(marshalled.Listeners))
	assert.Equal(t, "lds resource", marshalled.Listeners[0].(*listenerv3.Listener).Name)
}

func TestMarshalDiscoveryResponse(t *testing.T) {
	listener := &listenerv3.Listener{
		Name: "lds resource",
	}
	listenerAny, err := ptypes.MarshalAny(listener)
	assert.NoError(t, err)
	resp := discoveryv3.DiscoveryResponse{
		VersionInfo: "1",
		TypeUrl:     "type.googleapis.com/envoy.config.listener.v3.Listener",
		Resources: []*any.Any{
			listenerAny,
		},
	}
	marshalled := marshalDiscoveryResponse(transport.NewResponseV3(&gcp.Request{}, &resp))
	assert.NotNil(t, marshalled)
	assert.Equal(t, resp.VersionInfo, marshalled.VersionInfo)
	assert.Equal(t, resp.TypeUrl, marshalled.TypeURL)
	assert.Equal(t, listener.Name, marshalled.Resources.Listeners[0].(*listenerv3.Listener).Name)
}
