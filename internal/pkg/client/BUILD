go_library(
    name = "client",
    srcs = [
        "client.go",
    ],
    visibility = ["//..."],
    deps = [
        "//internal/pkg/utils",
    ],
)

go_test(
    name = "client_test",
    srcs = [
        "client_test.go",
    ],
    deps = [
        ":client",
        "//internal/pkg/client/mocks",
        "//internal/pkg/utils",
        "//third_party/go:mock",
        "//third_party/go:testify",
    ],
)
