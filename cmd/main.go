package main

import (
	"fmt"
	flags "github.com/thought-machine/go-flags"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/cache"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/client"
	util "github.com/thought-machine/prometheus-multiplexer-sidecar/internal/pkg/utils"
	"github.com/thought-machine/prometheus-multiplexer-sidecar/pkg/server"
	"log"
)

var opts struct {
	MetricsEndpoint    string   `short:"p" long:"endpoint" description:"The endpoint the metrics are exposing to." default:"/metrics"`
	ExportMetricsPort  int      `short:"e" long:"export_to" description:"The port the metrics are exposing to." default:"13434"`
	ContainerLabelName string   `short:"n" long:"container_label" description:"The name of the container label which will be appended to multiplexed metrics." default:"container"`
	ScrapeInterval     int      `short:"i" long:"scrape_interval" description:"The time interval for the scraping process in milliseconds." default:"200"`
	ExcludedContainers []string `short:"x" long:"exclude_containers" description:"Containers that can be excluded from the scraping process." default:""`
	ContainerToPortMap []string `short:"m" long:"container_to_port_map" description:"The mapping between container and ports, formatted as <container>:<port>." required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("Could not parse flags: %v", err)
	}

	containerToPortMap, err := util.GenerateContainerToPortMap(opts.ContainerToPortMap)
	if err != nil {
		log.Fatalf("Failed to generate container:port map: %s", err)
	}

	if err := util.ValidateLabelName(opts.ContainerLabelName); err != nil {
		log.Fatalf("Invalid container label name %s : %w", opts.ContainerLabelName, err)
	}

	svr := server.NewServer(opts.ExportMetricsPort, cache.NewMetricCache(), client.NewClient(), containerToPortMap, opts.MetricsEndpoint)
	defer svr.Close()
	svr.Start(opts.ScrapeInterval, opts.ContainerLabelName)
	fmt.Printf("start the server on port: %d", opts.ExportMetricsPort)
	if err := svr.ServeOnPort(); err != nil {
		log.Panicf("Unable to start the server: %v", err)
	}
}
