# Copilot instructions for backend-structure-go

Build & run
- Local run: go run ./cmd/api/main.go
- Make: make run (build: make build)
- Binary build: go build -o ./bin/app ./cmd/api/main.go

Tests
- Run all: go test ./... (Makefile target: test-all)
- Package tests (examples):
  - go test -v ./internal/api/user/service
  - go test -v ./internal/api/user/postgres
- Storage/integration/e2e suites (env toggles):
  - STORAGE_TEST_ENABLED=true go test ./test/storage -v
  - INTEGRATION_TEST_ENABLED=true go test ./test/integration -v
  - E2E_TEST_ENABLED=true go test ./test/e2e -v
- Run a single test by name: go test -run '^TestName$' ./path/to/package -v

Environment
- Files: example.env, deploy.env
- Env keys used by config: DB_HOST, DB_PORT, DB_USERNAME, DB_PASS, DB_NAME, DB_SSL_MODE, PORT
- Integration/storage tests use testcontainers (see go.mod); Docker required when enabled

High-level architecture
- cmd/api: application entrypoint that starts internal/api.App
- internal/api: wiring (DB client, router, validator), component initialization, and initRoutes
- internal/api/<domain> (user, product): domain-level use cases, services, repository, postgres storage, DTOs, mocks
- internal/api/delivery/http: HTTP handlers, DTOs and route assembly
- pkg/*: utilities (router wrapper, server with graceful shutdown, sqlext DB client, httpext, validation extensions, constants)
- test/: higher-level test suites (storage, integration, e2e)

Key conventions
- Route ordering: initRoutes and route.MountAll expect handlers in fixed index order (0: product, 1: user). Preserve this when adding handlers.
- Router: chi v5; API patterns in pkg/constant (ApiPattern, V1, ProductsPattern, UsersPattern).
- DB client: sqlext provides a singleton via GetInstance(opts) — do not instantiate multiple clients.
- Validation: validatorext wraps go-playground/validator and is initialized centrally in config and passed to components.
- Server: pkg/server.Server uses functional options (WithReadTimeout, WithWriteTimeout) and ConfigureGracefulShutdown.
- Tests: some packages have package-scoped tests; use env toggles to enable storage/integration/e2e suites.
- Mocks: in-memory mocks located under internal/api/<domain>/mock for fast unit tests.

Scripts & hooks
- Git hooks: scripts/git-hooks/*.sh — run scripts/git-hooks/install-hooks.sh to install
- SQL helpers: scripts/db contains sample SQL used by tests/setup
- Repo helper commands: scripts/commands.txt and Makefile targets

AI assistant files
- No AI assistant config files (CLAUDE.md, .cursorrules, AGENTS.md, etc.) detected. Add repository-specific notes here if desired.

Summary
- Added .github/copilot-instructions.md with build/test commands, high-level architecture, and repo-specific conventions. If you want adjustments or more coverage (naming rules, additional domains), say which areas to expand.
