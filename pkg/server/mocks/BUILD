subinclude("///third_party/github/pleasings//go:mock")

go_mock(
    name = "mocks",
    interfaces = [
        "ResponseWriter",
        "MetricClient",
        "MetricCache",
    ],
    package = "github.com/thought-machine/prometheus-multiplexer-sidecar/pkg/server",
    src_lib = "//pkg/server",
    visibility = ["PUBLIC"],
)
