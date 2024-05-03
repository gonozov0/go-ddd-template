#!/bin/bash

temp_dir=$(mktemp -d)
cp -r ./docs/* "$temp_dir"

go generate ./...

diff -r "$temp_dir" ./docs
if [ $? -ne 0 ]; then
    echo "Swagger documentation files have changed. Please update them in your repository."
    exit 1
fi

rm -rf "$temp_dir"
