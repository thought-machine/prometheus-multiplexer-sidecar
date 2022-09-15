go_library(
    name = "cache",
    srcs = [
        "cache.go",
    ],
    visibility = ["//..."],
)

go_test(
    name = "cache_test",
    srcs = [
        "cache_test.go",
    ],
    deps = [
        ":cache",
        "//internal/pkg/cache/mocks",
        "//third_party/go:mock",
        "//third_party/go:testify",
    ],
)
