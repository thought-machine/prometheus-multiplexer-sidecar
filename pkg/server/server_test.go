package server

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/client"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/utils"
	mock_server "github.com/thought-machine/prometheus-multiplexer-sidecar/pkg/server/mocks"
)

const (
	metricPort     = 13434
	path           = "http://localhost/metrics"
	acceptEncoding = "gzip"
	endpoint       = "/metrics"
)

var containerToPortMap = map[string]int{
	"container1": 1,
	"container2": 2,
	"container3": 3,
}

func TestHandleMetrics(t *testing.T) {

	testCases := []struct {
		name           string
		metricResult   []byte
		ok             bool
		encodingType   string
		expectedResult []byte
	}{
		{
			"test expose with invalid metric cache",
			nil,
			false,
			acceptEncoding,
			[]byte(""),
		},

		{
			"test expose with metric cache of gzip request",
			[]byte("test"),
			true,
			acceptEncoding,
			utils.CompressDataToGzip([]byte("testtesttest")),
		},
		{
			"test expose with metric cache of non-gzip request",
			[]byte("test"),
			true,
			"compress",
			[]byte("testtesttest"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := client.NewClient()
			header := http.Header{}
			req, _ := http.NewRequest("GET", path, nil)
			req.Header.Add("Accept-Encoding", tc.encodingType)
			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mockCache := mock_server.NewMockMetricCache(ctr)
			mockCache.EXPECT().GetAndInvalidate(gomock.Any()).Return(tc.metricResult, tc.ok)
			mockCache.EXPECT().GetAndInvalidate(gomock.Any()).Return(tc.metricResult, tc.ok)
			mockCache.EXPECT().GetAndInvalidate(gomock.Any()).Return(tc.metricResult, tc.ok)
			mw := mock_server.NewMockResponseWriter(ctr)
			mw.EXPECT().Header().Return(header)
			if tc.encodingType == acceptEncoding {
				mw.EXPECT().Header().Return(header)
			}
			mw.EXPECT().WriteHeader(http.StatusOK)
			mw.EXPECT().Write(tc.expectedResult)

			server := NewServer(metricPort, mockCache, client, containerToPortMap, endpoint)
			server.HandleMetrics(mw, req)
		})
	}
}

func TestPopulateCacheForContainer(t *testing.T) {
	testCases := []struct {
		name          string
		labelName     string
		containerName string
		port          int
		metricBuff    *bytes.Buffer
		expectedCache []byte
	}{
		{
			"test correct cache update",
			"multiplexer",
			"container1",
			1,
			bytes.NewBuffer([]byte(`# TYPE new_metric untyped
new_metric 22222
`)),
			[]byte(`# TYPE new_metric untyped
new_metric{multiplexer="container1"} 22222
`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mc := mock_server.NewMockMetricClient(ctr)
			mc.EXPECT().ScrapeRawMetrics(context.Background(), tc.port, endpoint).Return(tc.metricBuff, nil)
			mockCache := mock_server.NewMockMetricCache(ctr)
			mockCache.EXPECT().Set(tc.containerName, tc.expectedCache)

			server := NewServer(metricPort, mockCache, mc, containerToPortMap, endpoint)
			server.PopulateCacheForContainer(tc.labelName, tc.containerName, tc.port)
		})
	}
}
