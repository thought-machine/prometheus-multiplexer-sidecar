go_binary(
    name = "metrics-multiplexer-sidecar",
    srcs = [
        "main.go",
    ],
    visibility = ["PUBLIC"],
    deps = [
        "//internal/pkg/cache",
        "//internal/pkg/client",
        "//internal/pkg/utils",
        "//pkg/server",
        "//third_party/go:go-flags",
    ],
)
