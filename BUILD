load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")


go_binary(
    name = "client",
    deps = ["//lib"],
    srcs = ["client.go"],
    visibility = ["//visibility:public"],
    pure = "on",
)

go_binary(
    name = "server",
    deps = ["//lib"],
    srcs = ["server.go"],
    visibility = ["//visibility:public"],
    pure = "on",
)

go_binary(
    name = "test",
    deps = ["//lib"],
    srcs = ["test.go"],
)