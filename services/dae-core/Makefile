MAIN_ENTRY=dae-core
BIN_DIR=bin
APP_NAME=dae-core

# Build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(MAIN_ENTRY)/main.go

build-linux:
	@echo "Building $(APP_NAME) for Linux..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux ./cmd/$(MAIN_ENTRY)/main.go

# Development
run:
	@FIRESTORE_EMULATOR_HOST=localhost:8080 GOOGLE_CLOUD_PROJECT=dae-project go run ./cmd/$(MAIN_ENTRY)/main.go

fmt:
	@echo "Formatting code..."
	gofmt -w -s .
	goimports -w .

lint:
	@echo "Running linters..."
	golangci-lint run ./...

# Testing
test:
	@echo "Running unit tests..."
	go test ./... -race -count=1 -short

test-integration:
	@echo "Running integration tests..."
	go test ./... -race -count=1

test-verbose:
	go test ./... -race -count=1 -v

cover:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

cover-html:
	@echo "Generating HTML coverage report..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy && go mod vendor

# Code Generation
gen:
	@echo "Generating protobuf code..."
	@mkdir -p ../../proto/gen
	@protoc -I ../../proto \
		-I ../../proto/third_party \
		--go_out=../../proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=../../proto/gen --go-grpc_opt=paths=source_relative \
		--validate_out=../../proto/gen --validate_opt=lang=go,paths=source_relative\
		../../proto/*.proto
# Docker
docker-up:
	@echo "Starting Docker services..."
	cd ../.. && docker compose up -d --build

docker-down:
	@echo "Stopping Docker services..."
	cd ../.. && docker compose down -v

docker-logs:
	cd ../.. && docker compose logs -f

# Cleanup
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

# Help
help:
	@echo "Available targets:"
	@echo "  build              - Build the application"
	@echo "  build-linux        - Build for Linux"
	@echo "  run                - Run the application"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linters"
	@echo "  test               - Run unit tests"
	@echo "  test-integration   - Run integration tests"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  cover              - Run tests with coverage"
	@echo "  cover-html         - Generate HTML coverage report"
	@echo "  tidy               - Tidy dependencies"
	@echo "  gen                - Generate protobuf code"
	@echo "  docker-up          - Start Docker services"
	@echo "  docker-down        - Stop Docker services"
	@echo "  docker-logs        - View Docker logs"
	@echo "  clean              - Clean build artifacts"

.PHONY: build build-linux run fmt lint test test-integration test-verbose \
        cover cover-html tidy gen docker-up docker-down docker-logs clean help
