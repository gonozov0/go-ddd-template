install:
	go install github.com/segmentio/golines@latest
	go install ../../../vendor/github.com/bombsimon/wsl/v4/cmd/wsl # from arcadia root
	arc config pre-commit-hook.education/lib/gopkg.lint education/lib/gopkg/lint.sh

lint:
	ya tool yo fix .
	go fmt ./...
	@# ignore `proto` dir because it's generated code
	wsl --fix $(shell find . -mindepth 1 -maxdepth 1 -type d -not -name "proto" -exec basename {} \; | awk '{print "./"$$0"/..."}') || true
	find . -name '*.go' ! -path "./proto/*" -exec ya tool yoimports -w {} +
	find . -name '*.go' ! -path "./proto/*" -exec golines -w -m 120 {} +
