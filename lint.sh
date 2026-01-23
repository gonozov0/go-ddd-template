#!/bin/sh
set -e

echo "go fmt..."
go fmt ./...

echo "wsl..."
# ignore `generated` dir
directories=$(find . -mindepth 1 -maxdepth 1 -type d -not -name "generated" -exec basename {} \; | awk '{print "./"$0"/..."}')
wsl $directories

echo "goimports..."
find . -name "*.go" ! -path "./generated/*" -exec goimports -local go-ddd-template/ -w {} +

echo "golines..."
find . -name '*.go' ! -path "./generated/*" -exec golines -w -m 120 {} +

# check if there have been any changes
if [ -n "$(git diff --name-only)" ]; then
  echo "Changes detected after running linters"
  exit 1
fi
