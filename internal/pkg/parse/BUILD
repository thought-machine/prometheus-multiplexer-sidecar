go_library(
    name = "parse",
    srcs = [
        "parse.go",
    ],
    visibility = ["//..."],
    deps = [
        "//third_party/go:client_model",
        "//third_party/go:prometheus_common",
    ],
)

go_test(
    name = "parse_test",
    srcs = [
        "parse_test.go",
    ],
    deps = [
        ":parse",
        "//third_party/go:protobuf",
        "//third_party/go:testify",
    ],
)
