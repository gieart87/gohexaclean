.PHONY: help build run test clean docker-up docker-down proto migrate

# Variables
APP_NAME=gohexaclean
HTTP_SERVER=cmd/http/main.go
GRPC_SERVER=cmd/grpc/main.go
PROTO_DIR=api/proto
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Colors for output
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m

## help: Display this help message
help:
	@echo "$(COLOR_BOLD)Available commands:$(COLOR_RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(COLOR_GREEN)%-15s$(COLOR_RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(COLOR_BOLD)%s$(COLOR_RESET)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

## run-http: Run HTTP server
run-http:
	@echo "$(COLOR_GREEN)Starting HTTP server...$(COLOR_RESET)"
	go run $(HTTP_SERVER)

## run-grpc: Run gRPC server
run-grpc:
	@echo "$(COLOR_GREEN)Starting gRPC server...$(COLOR_RESET)"
	go run $(GRPC_SERVER)

## build: Build both HTTP and gRPC servers
build:
	@echo "$(COLOR_GREEN)Building HTTP server...$(COLOR_RESET)"
	go build -o bin/http-server $(HTTP_SERVER)
	@echo "$(COLOR_GREEN)Building gRPC server...$(COLOR_RESET)"
	go build -o bin/grpc-server $(GRPC_SERVER)
	@echo "$(COLOR_GREEN)Build complete!$(COLOR_RESET)"

## build-http: Build HTTP server only
build-http:
	@echo "$(COLOR_GREEN)Building HTTP server...$(COLOR_RESET)"
	go build -o bin/http-server $(HTTP_SERVER)

## build-grpc: Build gRPC server only
build-grpc:
	@echo "$(COLOR_GREEN)Building gRPC server...$(COLOR_RESET)"
	go build -o bin/grpc-server $(GRPC_SERVER)

##@ Testing

## test: Run tests
test:
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	go test -v -race -coverprofile=coverage.out ./...

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "$(COLOR_GREEN)Generating coverage report...$(COLOR_RESET)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)"

## test-unit: Run unit tests
test-unit:
	@echo "$(COLOR_GREEN)Running unit tests...$(COLOR_RESET)"
	go test -v -short ./test/unit/...

## test-integration: Run integration tests
test-integration:
	@echo "$(COLOR_GREEN)Running integration tests...$(COLOR_RESET)"
	go test -v ./test/integration/...

##@ Code Quality

## lint: Run linter
lint:
	@echo "$(COLOR_GREEN)Running linter...$(COLOR_RESET)"
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "$(COLOR_GREEN)Formatting code...$(COLOR_RESET)"
	go fmt ./...
	gofmt -s -w $(GO_FILES)

## vet: Run go vet
vet:
	@echo "$(COLOR_GREEN)Running go vet...$(COLOR_RESET)"
	go vet ./...

##@ Dependencies

## deps: Download dependencies
deps:
	@echo "$(COLOR_GREEN)Downloading dependencies...$(COLOR_RESET)"
	go mod download

## tidy: Tidy dependencies
tidy:
	@echo "$(COLOR_GREEN)Tidying dependencies...$(COLOR_RESET)"
	go mod tidy

## vendor: Vendor dependencies
vendor:
	@echo "$(COLOR_GREEN)Vendoring dependencies...$(COLOR_RESET)"
	go mod vendor

##@ Protocol Buffers

## proto: Generate protobuf files
proto:
	@echo "$(COLOR_GREEN)Generating protobuf files...$(COLOR_RESET)"
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "$(COLOR_GREEN)Protobuf generation complete!$(COLOR_RESET)"

## proto-install: Install protobuf tools
proto-install:
	@echo "$(COLOR_GREEN)Installing protobuf tools...$(COLOR_RESET)"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

##@ OpenAPI

## openapi: Generate code from OpenAPI spec
openapi:
	@echo "$(COLOR_GREEN)Generating code from OpenAPI spec...$(COLOR_RESET)"
	./scripts/generate-openapi.sh

## openapi-install: Install OpenAPI code generator
openapi-install:
	@echo "$(COLOR_GREEN)Installing oapi-codegen...$(COLOR_RESET)"
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest

## openapi-validate: Validate OpenAPI spec
openapi-validate:
	@echo "$(COLOR_GREEN)Validating OpenAPI spec...$(COLOR_RESET)"
	@if command -v redocly >/dev/null 2>&1; then \
		redocly lint api/openapi/*.yaml; \
	else \
		echo "$(COLOR_YELLOW)redocly not installed. Run: npm install -g @redocly/cli$(COLOR_RESET)"; \
	fi

##@ Database

## migrate-up: Run database migrations
migrate-up:
	@echo "$(COLOR_GREEN)Running migrations...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Please implement migration tool (e.g., golang-migrate)$(COLOR_RESET)"

## migrate-down: Rollback database migrations
migrate-down:
	@echo "$(COLOR_GREEN)Rolling back migrations...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)Please implement migration tool (e.g., golang-migrate)$(COLOR_RESET)"

## migrate-create: Create new migration file
migrate-create:
	@echo "$(COLOR_GREEN)Creating migration...$(COLOR_RESET)"
	@read -p "Enter migration name: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	touch internal/infra/db/migrations/$${timestamp}_$${name}.sql; \
	echo "$(COLOR_GREEN)Migration created: internal/infra/db/migrations/$${timestamp}_$${name}.sql$(COLOR_RESET)"

## seed: Run database seeders
seed:
	@echo "$(COLOR_GREEN)Running seeders...$(COLOR_RESET)"
	psql -h localhost -U postgres -d gohexaclean -f internal/infra/db/seeders/001_seed_users.sql

##@ Docker

## docker-up: Start docker containers
docker-up:
	@echo "$(COLOR_GREEN)Starting docker containers...$(COLOR_RESET)"
	docker-compose up -d

## docker-down: Stop docker containers
docker-down:
	@echo "$(COLOR_GREEN)Stopping docker containers...$(COLOR_RESET)"
	docker-compose down

## docker-build: Build docker images
docker-build:
	@echo "$(COLOR_GREEN)Building docker images...$(COLOR_RESET)"
	docker-compose build

## docker-logs: Show docker logs
docker-logs:
	docker-compose logs -f

## docker-clean: Clean docker resources
docker-clean:
	@echo "$(COLOR_GREEN)Cleaning docker resources...$(COLOR_RESET)"
	docker-compose down -v
	docker system prune -f

##@ Cleanup

## clean: Clean build artifacts
clean:
	@echo "$(COLOR_GREEN)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)Clean complete!$(COLOR_RESET)"

## clean-all: Clean everything including dependencies
clean-all: clean
	@echo "$(COLOR_GREEN)Cleaning dependencies...$(COLOR_RESET)"
	rm -rf vendor/
	go clean -modcache

##@ Utilities

## install-tools: Install development tools
install-tools:
	@echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## version: Show version information
version:
	@echo "Go version: $$(go version)"
	@echo "App: $(APP_NAME)"

.DEFAULT_GOAL := help
