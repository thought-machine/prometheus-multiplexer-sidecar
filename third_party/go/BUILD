package(default_visibility = ["PUBLIC"] )

go_module(
    name = "go-flags",
    licences = ["BSD-3-Clause"],
    module = "github.com/thought-machine/go-flags",
    version = "v1.6.1",
    deps = [
        ":x_sys",
    ],
)

go_module(
    name = "x_sys",
    install = ["..."],
    module = "golang.org/x/sys",
    version = "v0.0.0-20210510120138-977fb7262007",
)

go_module(
    name = "difflib",
    install = ["..."],
    licences = ["bsd-3-clause"],
    module = "github.com/pmezard/go-difflib",
    version = "v1.0.0",
)

go_module(
    name = "objx",
    licences = ["mit"],
    module = "github.com/stretchr/objx",
    version = "v0.4.0",
)

go_module(
    name = "spew",
    install = ["spew"],
    module = "github.com/davecgh/go-spew",
    version = "v1.1.1",
)

go_module(
    name = "yaml.v3",
    licences = ["apache-2.0"],
    module = "gopkg.in/yaml.v3",
    version = "v3.0.1",
)

go_module(
    name = "testify",
    install = [
        "assert",
        "mock",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "v1.8.0",
    deps = [
        ":difflib",
        ":objx",
        ":spew",
        ":yaml.v3",
    ],
)

go_module(
    name = "protobuf",
    install = [
        "proto",
        "encoding/prototext",
        "encoding/protowire",
        "internal/...",
        "runtime/protoimpl",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "runtime/protoiface",
        "types/known/timestamppb",
        "types/known/anypb",
        "types/known/durationpb",
    ],
    licences = ["bsd-3-clause"],
    module = "google.golang.org/protobuf",
    version = "v1.25.0",
)

go_module(
    name = "protobuf-github",
    install = [
        "proto",
        "ptypes",
        "ptypes/any",
        "ptypes/duration",
        "ptypes/timestamp",
    ],
    licences = ["bsd-3-clause"],
    module = "github.com/golang/protobuf",
    version = "v1.4.3",
    deps = [
        ":protobuf",
    ]
)

go_module(
    name = "golang_protobuf_extensions",
    install = ["..."],
    licences = ["apache-2.0"],
    module = "github.com/matttproud/golang_protobuf_extensions",
    version = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    deps = [
        ":protobuf-github",
    ],
)

go_module(
    name = "template",
    licences = ["bsd-3-clause"],
    module = "github.com/alecthomas/template",
    version = "a0175ee3bccc567396460bf5acd36800cb10c49c",
)

go_module(
    name = "units",
    licences = ["mit"],
    module = "github.com/alecthomas/units",
    version = "2efee857e7cfd4f3d0138cc3cbb1b4966962b93a",
)

go_module(
    name = "prometheus_common",
    install = [
        "model",
        "expfmt",
        "internal/...",
    ],
    licences = ["apache-2.0"],
    module = "github.com/prometheus/common",
    deps = [
        ":client_model",
        ":protobuf-github",
        ":protobuf",
        ":golang_protobuf_extensions",
    ],
    version = "v0.10.0",
)

go_module(
    name = "client_model",
    install = ["..."],
    licences = ["apache-2.0"],
    module = "github.com/prometheus/client_model",
    version = "v0.2.0",
    deps = [
        ":protobuf-github",
    ],
)


go_module(
    name = "logrus",
    licences = ["mit"],
    module = "github.com/sirupsen/logrus",
    version = "v1.6.0",
    deps = [
        ":x_crypto",
        ":x_sys",
    ],
)

go_module(
    name = "x_crypto",
    install = ["..."],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/crypto",
    strip = ["acme/autocert"],
    version = "v0.0.0-20210817164053-32db794688a5",
    deps = [
        ":x_sys",
        ":x_term",
    ],
)

go_module(
    name = "x_term",
    install = ["..."],
    module = "golang.org/x/term",
    version = "v0.0.0-20210615171337-6886f2dfbf5b",
    deps = [":x_sys"],
)

go_mod_download(
    name = "protobuf_download",
    licences = ["bsd-3-clause"],
    module = "github.com/golang/protobuf",
    version = "v1.4.3",
)

go_mod_download(
    name = "mockgen_download",
    licences = ["apache-2.0"],
    module = "github.com/golang/mock",
    version = "v1.6.0",
)

go_module(
    name = "mock",
    download = ":mockgen_download",
    install = ["..."],
    module = "github.com/golang/mock",
    deps = [":x_tools"],
)

go_module(
    name = "mockgen",
    binary = True,
    download = ":mockgen_download",
    install = ["mockgen"],
    module = "github.com/golang/mock",
    deps = [":x_tools"],
)

go_module(
    name = "x_tools",
    install = ["..."],
    module = "golang.org/x/tools",
    deps = [
        ":goldmark",
        ":x_sys",
        ":x_mod",
        ":x_net",
        ":x_sync",
    ],
    version = "v0.1.1",
)

go_module(
    name = "x_net",
    install = ["..."],
    module = "golang.org/x/net",
    deps = [
        ":x_text",
        ":x_term",
    ],
    version = "v0.0.0-20210405180319-a5a99cb37ef4",
)

go_module(
    name = "x_sync",
    install = ["..."],
    module = "golang.org/x/sync",
    version = "v0.0.0-20220907140024-f12130a52804",
)

go_module(
    name = "x_text",
    install = ["..."],
    module = "golang.org/x/text",
    strip = [
        "cmd",
        "message/pipeline",
    ],
    version = "v0.3.7",
)

go_module(
    name = "x_mod",
    install = ["..."],
    module = "golang.org/x/mod",
    deps = [
        ":x_xerrors",
        ":x_crypto",
    ],
    version = "v0.4.2",
)

go_module(
    name = "x_xerrors",
    install = ["..."],
    module = "golang.org/x/xerrors",
    version = "v0.0.0-20191011141410-1b5146add898",
)

go_module(
    name = "goldmark",
    install = ["..."],
    licences = ["MIT"],
    module = "github.com/yuin/goldmark",
    version = "v1.4.1",
)