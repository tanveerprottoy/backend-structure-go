## backend-structure-go
backend-structure-go module demonstrates a go http server with a well defined strtucture, which provides abstraction, testable code and using ideas from established architecture patterns.

## Running the app
```cli
go run ./cmd/api/main.go
```
```makefile
make run
```

## testing
unit test:

package wise tests:
user:
    service: go test -v ./internal/api/user/service
    storage: go test -v ./internal/api/user/postgres
product:
    service: go test -v ./internal/api/product/service
    storage: go test -v ./internal/api/product/postgres

test all:
go test ./...
(will also perform storage, integration, e2e based on env values)

storeage test: go test ./test/storage -v (env: STORAGE_TEST_ENABLED=true)

integration test: go test ./test/integration -v (env: INTEGRATION_TEST_ENABLED=true)

e2e test: go test ./test/e2e -v (env: E2E_TEST_ENABLED=true)

package: go test -v -cover ./<directory>/<package>

hint: go test file_name.go