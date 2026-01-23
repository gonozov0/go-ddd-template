# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Service Overview

This is a **Go template project** for Yandex Education services, providing a complete foundation for microservices within the Yandex Arcadia build system. The template implements a DDD (Domain-Driven Design) architecture with users, products, orders, and deliveries as example domains.

**Key Components:**
- **server**: Main API server (HTTP REST + gRPC)
- **consumers**: YDB topic message consumers (async event processing)
- **crons**: Scheduled tasks (e.g., handling created orders)
- **migrator**: Database migration runner (Postgres + YDB)

## Build System

This project can be built and tested using standard Go tooling. While it lives in Yandex Arcadia monorepo, native Go commands work perfectly.

**Quick Start:**
```bash
# If you encounter missing dependencies, generate proto files first:
make protogen

# Then use standard Go commands:
go build ./...
go test ./...
```

**Initial Setup:**
```bash
# For macOS
make install_on_mac

# For Linux
make install_on_linux
```

## Development Commands

**Testing:**
```bash
# Unit tests only
make unit_test
# or: go test -v ./internal/... -count=1

# Integration tests only (expects postgres/ydb/tvmtool running)
make integration_test
# or: go test -v ./integration_tests/... -p=1 -count=1

# Integration tests with dependencies (starts docker-compose)
make integration_test_with_deps

# All tests
make test
```

**Linting:**
```bash
# Run all linters and formatters
make lint

# Individual steps:
# - go fmt ./...           (standard go formatting)
# - wsl                    (whitespace linter)
# - goimports              (import organizer)
# - golines -m 120         (line length formatter)

# Note: Linters run automatically on pre-commit via lint.sh
```

**Database Migrations:**
```bash
# PostgreSQL migrations (requires .env with DB credentials)
make postgres_migrate_up                    # Apply all migrations
make postgres_migrate_down count=1          # Rollback N migrations
make postgres_create_migration name=foo     # Create new migration

# YDB migrations
make ydb_migrate_up
make ydb_migrate_down
make ydb_create_migration name=foo

# Run migrator directly
make run_migrator db=postgres
make run_migrator db=ydb
```

**Local Development:**
```bash
# Start local infrastructure (Postgres + YDB)
docker-compose up -d

# Run server
make run_server
# or: go run ./cmd/server

# Run crons
go run ./cmd/crons

# Run consumers
go run ./cmd/consumers
```

## Architecture

This service follows **Domain-Driven Design (DDD)** principles with clear layer separation:

### Directory Structure

```
cmd/                    # Entry points for each service component
├── server/            # Main API server (HTTP + gRPC)
├── consumers/         # Event consumers (YDB topics)
├── crons/             # Scheduled jobs
└── migrator/          # Database migrations

internal/
├── application/       # Application services (HTTP/gRPC handlers, use case orchestration)
│   ├── users/        # User management endpoints
│   ├── products/     # Product management endpoints
│   ├── orders/       # Order management endpoints
│   │   └── crons/    # Order-related cron jobs
│   ├── deliveries/   # Delivery management endpoints
│   ├── checks/       # Health checks and diagnostics
│   ├── http_server.go   # HTTP server setup (gRPC-gateway + Swagger)
│   ├── grpc_server.go   # gRPC server setup
│   ├── metric_server.go # Metrics endpoint
│   └── interfaces.go    # Repository interfaces
├── service/          # Domain services (business logic orchestration)
│   ├── users/        # User business logic
│   ├── products/     # Product business logic
│   ├── orders/       # Order business logic
│   └── deliveries/   # Delivery business logic
├── domain/           # Domain model (aggregates, entities, value objects, domain rules)
│   ├── users/        # User domain aggregate
│   ├── products/     # Product domain aggregate
│   ├── orders/       # Order domain aggregate
│   ├── deliveries/   # Delivery domain aggregate
│   └── shared/       # Shared kernel (shared value objects)
│       └── valueobjects/  # Common value objects (email, user_id, product_id, etc.)
├── infrastructure/   # Infrastructure (repository implementations, external clients)
│   ├── users/        # User repository implementations
│   │   ├── postgres.go   # PostgreSQL repository
│   │   └── in_memory.go  # In-memory repository for testing
│   ├── products/     # Product repository implementations
│   ├── orders/       # Order repository implementations
│   └── deliveries/   # Delivery repository implementations
├── config.go         # Configuration loading
├── run_servers.go    # Server initialization and lifecycle
├── run_consumers.go  # Consumer initialization
├── run_crons.go      # Cron job initialization
└── topics.go         # YDB topic definitions

pkg/                  # Reusable packages
└── auth/             # Primitive authentication via headers and metadata

proto/                # Protocol buffer definitions
├── server/           # Server API definitions
└── shared/           # Shared proto messages

migrations/           # Database migrations
├── postgres/         # PostgreSQL migrations (timestamped)
└── ydb/             # YDB migrations
```

### Key DDD Patterns

1. **Layered Architecture**: `application → service → domain ← infrastructure`
   - **Domain Model** (`internal/domain/`): Core business entities, aggregates, value objects, and domain rules. Pure business logic with no external dependencies.
   - **Domain Services** (`internal/service/`): Business logic orchestration that spans multiple entities or doesn't naturally belong to one aggregate. Manages transaction boundaries.
   - **Application Services** (`internal/application/`): Use case orchestration, API endpoints (HTTP/gRPC handlers), authentication/authorization. Translates between protocol buffers and domain models.
   - **Infrastructure** (`internal/infrastructure/`): Repository implementations (Postgres, YDB, in-memory), external service integrations. Implements interfaces defined in application layer.
   - **Repository Pattern**: Application layer defines repository interfaces, infrastructure implements them
   - Domain model has no dependencies on outer layers (Dependency Inversion Principle)

2. **Bounded Contexts**: Organized by business domains
   - `users/`: User management bounded context
   - `products/`: Product catalog bounded context
   - `orders/`: Order processing bounded context
   - `deliveries/`: Delivery management bounded context
   - `shared/`: Shared kernel with value objects used across contexts

3. **Server Architecture**: Multiple servers running concurrently using `errgroup`
   - **gRPC Server** (default port 8081): Main service endpoints with TVM authentication
   - **HTTP Server** (default port 8080): gRPC-gateway for REST API with Swagger UI at `/swagger/`
   - **Pprof Server** (default port 6060): Go profiling endpoints for performance debugging
   - All servers support graceful shutdown with configurable timeout (default 2s)

4. **Configuration**: All config loaded via environment variables in `internal/config.go`
   - Supports multiple environments: local, dev, testing, production
   - Uses `.env` file for local development (see `.env.example`)
   - Separate configs for server, crons, consumers, databases, TVM, logging, and tracing

5. **Database Access**:
   - **PostgreSQL** for main data storage (via `pgx`) with automatic master failover support
   - **YDB** for async messaging/topics and distributed storage
   - Migration system uses `goose` library
   - Connection pooling and cluster management handled by library wrappers
   - All repositories in `internal/infrastructure/*`

6. **Observability**:
   - Structured logging via `slog` with ErrorBooster integration
   - Distributed tracing with OpenTelemetry
   - Solomon metrics with automatic collection
   - Request ID propagation through `X-Request-Id` header

7. **Async Processing**:
   - **YDB Topics** for async messaging
   - **Topic Writers** for publishing events (e.g., order events)
   - **Consumers** (`cmd/consumers/`) for processing topic messages
   - Example: Order creation publishes to YDB topic, consumer processes asynchronously

8. **Cron Jobs**:
   - Periodic task execution framework
   - Example: `HandleCreatedOrders` cron in `internal/application/orders/crons/`
   - Configurable intervals via environment variables (e.g., `HANDLE_CREATED_ORDERS_INTERVAL`)

9. **Testing**:
   - **Unit tests**: `internal/application/*_test/` use in-memory repositories
   - **Integration tests**: `integration_tests/` use real Postgres/YDB/TVM
   - Test suites use testify for assertions
   - Helper functions in `*_test/helpers/` for test data generation

## Testing Guidelines

This project uses two types of tests: unit tests and integration tests.

### Unit Tests

Unit tests verify business logic at the application layer for each interaction method (HTTP/gRPC servers, cron, consumers). All external dependencies are replaced with mocks (in-memory implementations), middleware is not tested.

**Helper Functions for Tests** (in `shared/helpers/` directory):

- **generate.go** — creation of domain objects using the functional options pattern. Objects with random data are created without options.
- **create.go** — simplified creation of objects in storage.
- **convert.go** — transformation of domain objects into structures for the application layer.

**Recommendation**: pass the suite to helpers to perform checks inside them. Use an alias when importing: `<aggregate_name>unithelper` (for example, `orderunithelpers`).

### Integration Tests

Integration tests verify the operation of a fully deployed service. **Important**: you cannot directly access external systems (DB, etc.) from these tests, as they run in sandbox without network access to external subsystems.

For each aggregate in the `helpers/` directory, an auxiliary suite is implemented for managing test objects.

In integration tests, you **do not need** to implement structures for HTTP interaction; you should use generated structures from proto (they have json tags). But for correct operation of Marshal/Unmarshal operations, you need to use functions from the `grpcutils` package.

## Critical Development Notes

### SQL Query Style

**ALL SQL queries MUST use table aliases.** This prevents silent bugs where wrong tables are referenced.

**Bad (will delete all modules!):**
```sql
DELETE FROM modules
WHERE "venue_season_id" IN
    (SELECT "venue_season_id" FROM venue_seasons WHERE season_id = $1)
```

**Good (will error if field doesn't exist):**
```sql
DELETE FROM modules m
WHERE m.venue_season_id IN
    (SELECT vs.venue_season_id FROM venue_seasons vs WHERE vs.season_id = $1)
```

## Common Development Workflows

### Adding a New API Endpoint

1. Define proto message/service in `proto/server/protobuf.proto`
2. Run `make protogen` to generate Go code
3. Add domain entities/value objects in `internal/domain/<context>/`
4. Implement business logic in `internal/service/<context>/`
5. Create repository methods in `internal/infrastructure/<context>/`
6. Add HTTP/gRPC handlers in `internal/application/<context>/`
7. Wire up in `internal/run_servers.go`
8. Add tests in appropriate test files
9. Run `make lint` before committing
