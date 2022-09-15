package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock_client "github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/client/mocks"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/utils"
)

const (
	metricPort = 13434
	endpoint   = "/metrics"
)

func TestScrapeRawMetrics(t *testing.T) {

	doerr := errors.New("doErrResp")
	testCases := []struct {
		name           string
		response       *http.Response
		err            error
		expectedResult *bytes.Buffer
		expectedErr    error
	}{
		{
			"reads a gzip response",
			&http.Response{
				Header:     map[string][]string{"Content-Encoding": {"gzip"}},
				Body:       ioutil.NopCloser(bytes.NewReader(utils.CompressDataToGzip([]byte("This is test 1234.")))),
				Status:     "200 OK",
				StatusCode: 200,
			},
			nil,
			bytes.NewBuffer([]byte("This is test 1234.")),
			nil,
		},
		{
			"reads a compress response",
			&http.Response{
				Header:     map[string][]string{"Content-Encoding": {"compress"}},
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("This is test 1234."))),
				Status:     "200 OK",
				StatusCode: 200,
			},
			nil,
			bytes.NewBuffer([]byte("This is test 1234.")),
			nil,
		},

		// 'Do' error test.
		{
			"generate Do error",
			nil,
			doerr,
			nil,
			doerr,
		},

		{
			"non-successful response",
			&http.Response{
				StatusCode: 404,
				Status:     "404 Not Found",
				Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
			},
			nil,
			nil,
			errStatusNotOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ctr := gomock.NewController(t)
			defer ctr.Finish()
			mc := mock_client.NewMockHTTPClient(ctr)
			req, err := http.NewRequest("GET", fmt.Sprintf("%s:%d%s", localPath, metricPort, endpoint), nil)
			require.NoError(t, err)
			req.Header.Add("Accept-Encoding", acceptEncoding)
			req.Header.Add("Accept", accept)

			mc.EXPECT().Do(req).Return(tc.response, tc.err)
			client := Client{httpClient: mc}
			metric, err := client.ScrapeRawMetrics(context.Background(), metricPort, endpoint)
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedResult, metric)
		})
	}
}
