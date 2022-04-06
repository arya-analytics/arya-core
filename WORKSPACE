workspace(name = "arya_core")

# |||| BASE ||||

# Load all of it. All of it. Now.
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# |||| HTTP ARCHIVES ||||

# || GO ||

# Load go bazel rules, which gives us access to go bazel SDK.
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "d6b2513456fe2229811da7eb67a444be7785f5323c6708b38d851d2b51e54d83",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.30.0/rules_go-v0.30.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.30.0/rules_go-v0.30.0.zip",
    ],
)

# || GAZELLE ||

# Gazelle allows us to auto-gen BUILD.bazel files across directories.
http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

# |||| GO BUILD RULES ||||

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18")

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

gazelle_dependencies()

# |||| NODEJS BUILD RULES ||||

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "523da2d6b50bc00eaf14b00ed28b1a366b3ab456e14131e9812558b26599125c",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/5.3.1/rules_nodejs-5.3.1.tar.gz"],
)

load("@build_bazel_rules_nodejs//:repositories.bzl", "build_bazel_rules_nodejs_dependencies")

build_bazel_rules_nodejs_dependencies()

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")

node_repositories(
    node_version = "17.8.0",
    yarn_version = "1.22.18",
)

yarn_install(
    name = "npm",
    package_json = "//pkg/ui:package.json",
    yarn_lock = "//pkg/ui:yarn.lock",
)
