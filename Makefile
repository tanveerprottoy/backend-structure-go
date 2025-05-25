BIN_DIR = ./bin
APP_PATH = ./cmd/api/main.go

build:
	go build -o .$(BIN_DIR)/app $(APP_PATH)

run:
	go run $(APP_PATH)

test-all:
	go test -v ./...
