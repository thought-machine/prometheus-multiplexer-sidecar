package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/mutate"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/parse"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/utils"
)

const (
	contentEncoding = "gzip"
	contentType     = "text/plain; charset=utf-8"
)

// MetricClient is the metricClient interface for scraping metrics.
type MetricClient interface {
	ScrapeRawMetrics(ctx context.Context, port int, endpoint string) (*bytes.Buffer, error)
}

// MetricCache is the cache interface for metric storage.
type MetricCache interface {
	GetAndInvalidate(containerName string) ([]byte, bool)
	Set(containerName string, metrics []byte)
}

// HTTPServer is a server interface that implements functionality for handling HTTP requests.
type HTTPServer interface {
	ListenAndServe() error
	Close() error
}

// ResponseWriter is the interface for a http response writer.
type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

// Server is a wrapper around an HTTP server and have the functionality to scrape all containers within a pod and return the contents of the cache.
type Server struct {
	httpServer         HTTPServer
	cache              MetricCache
	metricClient       MetricClient
	containerToPortMap map[string]int
	path               string
}

// NewServer instantiates a new server.
func NewServer(metricPort int, cache MetricCache, client MetricClient, containerToPortMap map[string]int, endpoint string) *Server {
	return &Server{
		&http.Server{
			Addr: fmt.Sprintf(":%d", metricPort),
		},
		cache,
		client,
		containerToPortMap,
		endpoint,
	}
}

// HandleMetrics is the handler for exposing metrics. It will fetch all the available metrics from cache,
//invalidate all their entries on cache, and finally serve them to the metric path.
func (server *Server) HandleMetrics(writer http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Warningf("Invalid http %s method for getting metrics from server.", r.Method)
		return
	}

	metrics := make([]byte, 0)
	for containerName := range server.containerToPortMap {
		metric, ok := server.cache.GetAndInvalidate(containerName)
		if ok {
			metrics = append(metrics, metric...)
		} else {
			log.Errorf("Missing metric in container : %s", containerName)
		}
	}

	writer.Header().Set("Content-Type", contentType)

	encodingHeaders := r.Header.Get("Accept-Encoding")
	parts := strings.Split(encodingHeaders, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "gzip" || strings.HasPrefix(part, "gzip;") {
			writer.Header().Set("Content-Encoding", contentEncoding)
			metrics = utils.CompressDataToGzip(metrics)
			break
		}
	}
	writer.WriteHeader(http.StatusOK)

	if len(metrics) == 0 {
		if _, err := writer.Write([]byte("")); err != nil {
			log.Errorf("Failed to write empty metric on path %s: %v", server.path, err)
		}
	} else {
		if _, err := writer.Write(metrics); err != nil {
			log.Errorf("Failed to write metrics data on path %s: %v", server.path, err)
		}
	}
}

// ServeOnPort starts the server on the given port.
func (server *Server) ServeOnPort() error {
	http.HandleFunc(server.path, server.HandleMetrics)
	if err := server.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start the server on the path %s: %v", server.path, err)
	}
	return nil
}

// PopulateCacheForContainer updates the specified metrics on the metric cache.
func (server *Server) PopulateCacheForContainer(labelName string, containerName string, port int) {
	rawMetrics, err := server.metricClient.ScrapeRawMetrics(context.Background(), port, server.path)
	if err != nil {
		log.Errorf("Failed to scrape metrics on path %s for container %s on port %d: %v", server.path, containerName, port, err)
		return
	}
	metricFamilyMap, err := parse.Unmarshal(rawMetrics)
	if err != nil {
		log.Errorf("Failed to unmarshal the metrics on path %s: %v", server.path, err)
		return
	}
	if err = mutate.AppendLabelToMetrics(labelName, containerName, metricFamilyMap); err != nil {
		log.Errorf("Failed to append label %s to metrics on path %s: %v", labelName, server.path, err)
		return
	}
	rawMetricsBuff, err := parse.Marshal(metricFamilyMap)
	if err != nil {
		log.Errorf("Failed to marshal the metrics on path %s: %v", server.path, err)
		return
	}
	server.cache.Set(containerName, rawMetricsBuff.Bytes())
}

// Start starts the server for exposing metrics and listen on each port to scrape the container.
func (server *Server) Start(internalMs int, containerLabelName string) {
	ticker := time.NewTicker(time.Duration(internalMs) * time.Millisecond)
	for container, port := range server.containerToPortMap {
		c := container
		p := port
		go func() {
			for {
				select {
				case <-ticker.C:
					server.PopulateCacheForContainer(containerLabelName, c, p)
				}
			}
		}()
	}
}

// Close is a wrapper function for the Close() function of http.Client.
func (server *Server) Close() {
	server.httpServer.Close()
}
