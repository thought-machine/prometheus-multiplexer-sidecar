package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/utils"
)

const (
	localPath      = "http://localhost"
	defaultTimeout = 20 * time.Second
	acceptEncoding = "gzip"
	accept         = "application/openmetrics-text; version=0.0.1,text/plain;version=0.0.4;q=0.5,*/*;q=0.1`"
)

var errStatusNotOK = errors.New("received a non-OK status")

// HTTPClient is a client interface that implements functionality for doing HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a thin wrapper around an HTTP client used to scrape the raw metrics from different containers in a pod.
type Client struct {
	httpClient HTTPClient
}

// NewClient instantiates a new client.
func NewClient() *Client {
	return &Client{&http.Client{Timeout: defaultTimeout}}
}

// ScrapeRawMetrics scrapes the metrics on the given port and returns raw metrics.
func (client *Client) ScrapeRawMetrics(ctx context.Context, port int, endpoint string) (*bytes.Buffer, error) {

	var rawMetrics bytes.Buffer
	url := fmt.Sprintf("%s:%d%s", localPath, port, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}
	req.Header.Add("Accept-Encoding", acceptEncoding)
	req.Header.Add("Accept", accept)
	resp, err := client.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to do GET request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("server returned HTTP status %s: %w", resp.Status, errStatusNotOK)
	}

	if resp.Header.Get("Content-Encoding") == "gzip" {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read response body: %w", err)
		}

		data, err = utils.GzipToCompressData(data)
		if err != nil {
			return nil, fmt.Errorf("failed to unzip gzip data: %w", err)
		}
		rawMetrics = *bytes.NewBuffer(data)
		return &rawMetrics, nil
	}

	if _, err := io.Copy(&rawMetrics, resp.Body); err != nil {
		return nil, fmt.Errorf("unable to copy raw metrics from response body: %w", err)
	}

	return &rawMetrics, nil
}
