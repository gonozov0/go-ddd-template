#!/bin/sh
set -e

echo "yo fix..."
ya tool yo fix .

echo "go fmt..."
go fmt ./...

echo "wsl..."
# ignore `proto` dir because it's generated code
directories=$(find . -mindepth 1 -maxdepth 1 -type d -not -name "proto" -exec basename {} \; | awk '{print "./"$0"/..."}')
wsl $directories

echo "yoimports..."
find . -name '*.go' ! -path "./proto/*" -exec ya tool yoimports -w {} +

echo "golines..."
find . -name '*.go' ! -path "./proto/*" -exec golines -w -m 120 {} +

# check if there have been any changes
if [ -n "$(arc diff --name-only)" ]; then
  echo "Changes detected after running linters"
  exit 1
fi
