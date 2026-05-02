# AI Agents Configuration

This document defines specialized agents and their configurations for the backend-structure-go project.

## Build & Test Agents

### Agent: Builder
**Purpose**: Compile and build the application  
**Commands**:
- Standard build: `make build`
- Development run: `go run ./cmd/api/main.go`
- Binary build: `go build -o ./bin/app ./cmd/api/main.go`

**Environment**: Requires Go 1.x, Make installed

---

### Agent: Test Runner
**Purpose**: Execute unit, integration, and e2e tests  
**Commands**:
- Run all tests: `go test ./...` or `make test-all`
- Package tests: `go test -v ./internal/api/user/service`
- Storage tests: `STORAGE_TEST_ENABLED=true go test ./test/storage -v`
- Integration tests: `INTEGRATION_TEST_ENABLED=true go test ./test/integration -v`
- E2E tests: `E2E_TEST_ENABLED=true go test ./test/e2e -v`
- Single test: `go test -run '^TestName$' ./path/to/package -v`

**Environment**: Requires Docker for storage/integration/e2e suites; testcontainers via go.mod

---

## Code Analysis Agents

### Agent: Architecture Inspector
**Purpose**: Analyze and verify project structure  
**Key Focus Areas**:
- cmd/api: Application entrypoint
- internal/api: Core wiring, DB client, router, validator
- internal/api/<domain> (user, product): Domain services, repositories, DTOs
- internal/api/delivery/http: HTTP handlers and routes
- pkg/*: Utilities (router, server, sqlext, httpext, validation, constants)
- test/: Storage, integration, e2e suites

**Validation Rules**:
- Route ordering in initRoutes: 0=product, 1=user (fixed index order)
- Router: chi v5; patterns in pkg/constant
- DB: sqlext singleton via GetInstance()
- Validation: centrallyinitialized validatorext
- Server: pkg/server.Server with functional options
- Mocks: in-memory mocks under internal/api/<domain>/mock

---

### Agent: Configuration Verifier
**Purpose**: Validate environment and configuration  
**Environment Files**: example.env, deploy.env  
**Required Keys**: DB_HOST, DB_PORT, DB_USERNAME, DB_PASS, DB_NAME, DB_SSL_MODE, PORT

---

## Database Agents

### Agent: Database Manager
**Purpose**: Manage database setup, migrations, and cleanup  
**SQL Resources**: Located in scripts/db/  
**Test Setup**: Database initialization scripts for storage/integration tests  
**Tools**: Testcontainers integration for isolated test databases

---

## Git & Hooks Agents

### Agent: Git Hooks Manager
**Purpose**: Install and manage pre-commit hooks  
**Installation**: `scripts/git-hooks/install-hooks.sh`  
**Hook Location**: scripts/git-hooks/*.sh  
**Purpose**: Automated linting, testing, and code quality checks before commits

---

## Domain-Specific Agents

### Agent: User Domain Specialist
**Purpose**: Work with user domain components  
**Components**:
- Service: internal/api/user/service
- Repository: internal/api/user/repository
- Storage: internal/api/user/postgres
- DTOs: internal/api/user/dto
- Mocks: internal/api/user/mock
- Tests: internal/api/user/*_test.go

---

### Agent: Product Domain Specialist
**Purpose**: Work with product domain components  
**Components**:
- Service: internal/api/product/service
- Repository: internal/api/product/repository
- Storage: internal/api/product/postgres
- DTOs: internal/api/product/dto
- Mocks: internal/api/product/mock
- Tests: internal/api/product/*_test.go

---

## Utility Agents

### Agent: Router & Routing
**Purpose**: Manage HTTP routing and API patterns  
**Tools**:
- Router framework: chi v5
- API patterns: pkg/constant (ApiPattern, V1, ProductsPattern, UsersPattern)
- Route assembly: internal/api/delivery/http

---

### Agent: Validation Manager
**Purpose**: Manage request validation across the application  
**Tool**: validatorext wrapper around go-playground/validator  
**Initialization**: Centralized in config, passed to all components

---

### Agent: Server Lifecycle Manager
**Purpose**: Manage server startup, shutdown, and graceful handling  
**Tool**: pkg/server.Server with functional options  
**Options**: WithReadTimeout, WithWriteTimeout, ConfigureGracefulShutdown

---

## Helper Resources

**Scripts**: scripts/commands.txt lists helper commands  
**Makefile**: Available targets for common tasks  
**Documentation**: .github/copilot-instructions.md for detailed conventions
