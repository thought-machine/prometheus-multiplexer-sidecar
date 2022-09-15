go_library(
    name = "mutate",
    srcs = [
        "mutate.go",
    ],
    visibility = ["//..."],
    deps = [
        "//third_party/go:client_model",
        "//third_party/go:protobuf",
    ],
)

go_test(
    name = "mutate_test",
    srcs = [
        "mutate_test.go",
    ],
    deps = [
        ":mutate",
        "//internal/pkg/parse",
        "//third_party/go:protobuf",
        "//third_party/go:testify",
    ],
)
