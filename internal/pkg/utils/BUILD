go_library(
    name = "utils",
    srcs = [
        "util.go",
    ],
    visibility = ["//..."],
    deps = [
        "//third_party/go:prometheus_common",
    ],
)

go_test(
    name = "util_test",
    srcs = [
        "util_test.go",
    ],
    deps = [
        ":utils",
        "//third_party/go:testify",
    ],
)
