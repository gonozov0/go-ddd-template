.PHONY: help unit_test integration_test e2e_test test lint coverage_report

help:
	cat Makefile

unit_test:
	go test -v ./internal/...

integration_test:
	go test -v ./integration_test/...

e2e_test:
	go test -v ./e2e_test/...

test: unit_test integration_test e2e_test

lint:
	go fmt ./...
	find . -name '*.go' -exec goimports -local go-echo-ddd-template/ -w {} +
	find . -name '*.go' -exec golines -w {} -m 120 \;
	swag fmt ./...
	golangci-lint run ./...


coverage_report:
	go test -coverpkg=./... -count=1 -coverprofile=.coverage.out ./...
	go tool cover -html .coverage.out -o .coverage.html
	open ./.coverage.html
