#!/bin/bash

temp_dir=$(mktemp -d)
cp -r ./generated/* "$temp_dir"

protoc \
		-I proto \
		-I $(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.5 \
		--go_out=generated --go_opt=paths=source_relative \
		--go-grpc_out=generated --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=generated --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=generated \
		proto/server/protobuf.proto

diff -r "$temp_dir" ./generated
if [ $? -ne 0 ]; then
    echo "Generated files have changes. Please update them in your local repository."
    exit 1
fi

rm -rf "$temp_dir"
