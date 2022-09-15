go_library(
    name = "server",
    srcs = [
        "server.go",
    ],
    visibility = ["//..."],
    deps = [
        "//internal/pkg/mutate",
        "//internal/pkg/parse",
        "//internal/pkg/utils",
        "//third_party/go:client_model",
        "//third_party/go:logrus",
    ],
)

go_test(
    name = "server_test",
    srcs = [
        "server_test.go",
    ],
    deps = [
        ":server",
        "//internal/pkg/cache",
        "//internal/pkg/client",
        "//internal/pkg/utils",
        "//pkg/server/mocks",
        "//third_party/go:mock",
        "//third_party/go:testify",
    ],
)
