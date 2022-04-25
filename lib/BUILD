load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(name = "gazelle")

go_library(
    name = "lib",
    srcs = ["common.go"],
    importpath = "github.com/hika019/File_transfer/lib",
    visibility = ["//visibility:public"],

)