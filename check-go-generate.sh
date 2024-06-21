#!/bin/bash

temp_dir=$(mktemp -d)
cp -r ./generated/* "$temp_dir"

go generate ./...

diff -r "$temp_dir" ./generated
if [ $? -ne 0 ]; then
    echo "Generated files have changes. Please update them in your local repository."
    exit 1
fi

rm -rf "$temp_dir"
