include .env
export

help:
	cat Makefile

protogen:
	protoc \
		-I proto \
		-I $(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.27.5 \
		--go_out=generated --go_opt=paths=source_relative \
		--go-grpc_out=generated --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=generated --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=generated \
		proto/server/protobuf.proto

install:
	bash install.sh

install_on_linux:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	make install

install_on_mac:
	brew install goose
	brew install protobuf
	make install

unit_test:
	go test -v ./internal/... -count=1

integration_test:
	# only 1 process because of postgres running
	go test -v ./integration_tests/... -p=1 -count=1

integration_test_with_deps:
	@bash -c ' \
		set -e; \
		cleanup() { \
			echo "Cleaning up..."; \
			docker-compose down; \
		}; \
		trap cleanup EXIT INT TERM HUP; \
		docker-compose up -d; \
		./wait_for_ydb.sh; \
		make integration_test \
	'

test: unit_test integration_test_with_deps

lint:
	go fmt ./...
	wsl --fix $(shell find . -mindepth 1 -maxdepth 1 -type d -not -name "generated" -exec basename {} \; | awk '{print "./"$$0"/..."}') || true
	find . -name "*.go" ! -path "./generated/*" -exec goimports -local go-ddd-template -w {} +
	find . -name '*.go' ! -path "./generated/*" -exec golines -w -m 120 {} +

coverage_report:
	go test -p=1 -coverpkg=./... -count=1 -coverprofile=.coverage.out ./...
	go tool cover -html .coverage.out -o .coverage.html
	open ./.coverage.html

cpu_profile:
	go test -cpuprofile=profiles/cpu.prof  ./... -p=1 -count=1
	go tool pprof -http=:6061 profiles/cpu.prof

mem_profile:
	go test -memprofile=profiles/mem.prof ./... -p=1 -count=1
	go tool pprof -http=:6061 profiles/mem.prof

POSTGRES_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOSTS):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=$(if $(filter $(POSTGRES_SSL),true),require,disable)

postgres_migrate_up:
	goose postgres -dir migrations/postgres "$(POSTGRES_URL)" up

postgres_migrate_down:
	goose postgres -dir migrations/postgres "$(POSTGRES_URL)" down $(count)

postgres_create_migration:
ifndef name
	$(error name parameter is required. Usage: make postgres_create_migration name=init)
endif
	goose postgres -dir migrations/postgres "$(POSTGRES_URL)" create $(name) sql

YDB_URL=grpc://$(YDB_ENDPOINT)/$(YDB_DATABASE)?go_query_mode=scripting&go_fake_tx=scripting&go_query_bind=declare,numeric&token=$(YDB_TOKEN)

ydb_migrate_up:
	goose ydb -dir migrations/ydb "$(YDB_URL)" up

ydb_migrate_down:
	goose ydb -dir migrations/ydb "$(YDB_URL)" reset

ydb_create_migration:
ifndef name
	$(error name parameter is required. Usage: make ydb_create_migration name=init)
endif
	goose ydb -dir migrations/ydb "$(YDB_URL)" create $(name) sql

run_migrator:
ifndef db
	$(error name parameter is required. Usage: make ydb_create_migration db=postgres)
endif
	go run ./cmd/migrator/ -db=$(db)

run_server:
	go run ./cmd/server/
