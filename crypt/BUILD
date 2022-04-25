load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")

go_library(
    name = "crypt",
    srcs = ["rsa.go", "aes.go"],
    deps = ["//lib"],
    importpath = "github.com/hika019/File_transfer/crypt",
    visibility = ["//visibility:public"],

)


go_binary(
    name = "test",
    embed = [":crypt"],
    srcs = ["test.go"],
    visibility = ["//visibility:public"],
    pure = "on",
)
