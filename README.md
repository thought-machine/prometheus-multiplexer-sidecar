# Prometheus Multiplexer Sidecar

A simple binary that can be deployed as a sidecar in Kubernetes for collecting metrics from
different containers within a Pod and serving these on a single endpoint by mutating the metrics
with the container name as a label.

This was originally started as an intern project by [Chengcheng Xing](https://github.com/Xichc1127) during her internship at [Thought Machine](https://thoughtmachine.net/).

## Introduction

This binary is a new tool for scraping metrics with Prometheus. It works with multi-container pods,
collects metrics of the containers and serve them to the Prometheus. The open-source Prometheus
library is used during the deployment [[link]](https://github.com/prometheus/prometheus).

Specific libraries include:

- Encoding & decoding [[expfmt]](https://pkg.go.dev/github.com/prometheus/common/expfmt)
- Metric parsing and representation
  [[client_model]](https://pkg.go.dev/github.com/prometheus/client_model@v0.2.0/go)

## Options

| **Short Flag** |      **Long Flag**      |                                **Description**                                 | **Default** |
| :------------: | :---------------------: | :----------------------------------------------------------------------------: | :---------: |
|       -p       |       --endpoint        |                   The endpoint the metrics are exposing to.                    |  /metrics   |
|       -e       |       --export_to       |                     The port the metrics are exposing to.                      |    13434    |
|       -n       |    --container_label    | The name of the container label which will be appended to multiplexed metrics. |  container  |
|       -i       |    --scrape_interval    |          The time interval for the scraping process in milliseconds.           |     200     |
|       -x       |  --exclude_containers   |           Containers that can be excluded from the scraping process.           |     ""      |
|       -m       | --container_to_port_map | The mapping between container and ports, formatted as \<container\>:\<port\>.  |     N/A     |

## Set Up Your Prometheus Multiplexed Sidecar

### Adding It As A Container In Your Server

In your server's yaml file, deploy this multiplexer as one of its container:

```aidl
-containers:
    ...
    - name: prom-multiplexer-sidecar
          image: <IMAGE_PATH>
          command:
            - "/metrics-multiplexer-sidecar"
            - "--endpoint=<ENDPOINT>"
            - "--export_to=<PORT>"
            - "--container_label=<LABLE_NAME>"
            - "--scrape_interval=<TIME_INTERVAL>"
            - "--container_to_port_map=<CONTAINER_NAME1>:<PORT_VALUE1>"
            - "--container_to_port_map=<CONTAINER_NAME2>:<PORT_VALUE2>"
            ...
          resources:
            requests:
              memory: 32Mi
              cpu: 2m
            limits:
              memory: 32Mi
              cpu: 50m
          ports:
            - containerPort: <PORT>
              name: <PORT_NAME>
```

The sidecar should be the only one who owns the metric port. Edit the ports of your other containers
within this service to make sure no containers have the same port value.

## Test Your Setup

After deploying the sidecar, you can hit the metrics endpoint and check if the metrics from
different containers are collected correctly.

## Diagram

The sequence diagram which explains how the multiplexer sidecar works can be viewed here:
![sequence diagram](./doc/SequenceDiagram.png?raw=true "sequence diagram for prometheus multiplexer sidecar")
