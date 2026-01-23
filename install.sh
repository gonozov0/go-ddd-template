#!/bin/sh
set -e

go install golang.org/x/tools/cmd/goimports@latest
go install github.com/segmentio/golines@latest
go install github.com/bombsimon/wsl/v5/cmd/wsl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
