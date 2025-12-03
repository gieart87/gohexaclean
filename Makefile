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

## test: Run all tests with coverage
test:
	@echo "$(COLOR_GREEN)Running tests...$(COLOR_RESET)"
	go test -v -race -coverprofile=coverage.out ./...
	@echo "$(COLOR_GREEN)Calculating coverage...$(COLOR_RESET)"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $$3}'

## test-coverage: Run tests with coverage report
test-coverage: test
	@echo "$(COLOR_GREEN)Generating coverage report...$(COLOR_RESET)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)Coverage report generated: coverage.html$(COLOR_RESET)"
	@open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "Please open coverage.html manually"

## test-unit: Run unit tests with coverage
test-unit:
	@echo "$(COLOR_GREEN)Running unit tests...$(COLOR_RESET)"
	go test -v -cover ./internal/adapter/outbound/pgsql/... ./internal/app/... ./internal/adapter/inbound/http/handler/...
	@echo "$(COLOR_GREEN)Unit tests complete!$(COLOR_RESET)"

## test-integration: Run integration tests
test-integration:
	@echo "$(COLOR_GREEN)Running integration tests...$(COLOR_RESET)"
	go test -v ./test/integration/...

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "$(COLOR_GREEN)Running tests with verbose output...$(COLOR_RESET)"
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

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

## migrate-up: Run database migrations using goose
migrate-up:
	@echo "$(COLOR_GREEN)Running migrations...$(COLOR_RESET)"
	cd internal/infra/db/migrations && goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=gohexaclean sslmode=disable" up

## migrate-down: Rollback database migrations
migrate-down:
	@echo "$(COLOR_GREEN)Rolling back migrations...$(COLOR_RESET)"
	cd internal/infra/db/migrations && goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=gohexaclean sslmode=disable" down

## migrate-status: Check migration status
migrate-status:
	@echo "$(COLOR_GREEN)Checking migration status...$(COLOR_RESET)"
	cd internal/infra/db/migrations && goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=gohexaclean sslmode=disable" status

## migrate-create: Create new migration file
migrate-create:
	@echo "$(COLOR_GREEN)Creating migration...$(COLOR_RESET)"
	@read -p "Enter migration name: " name; \
	cd internal/infra/db/migrations && goose create $${name} sql; \
	echo "$(COLOR_GREEN)Migration created in internal/infra/db/migrations/$(COLOR_RESET)"

## migrate-reset: Reset database (down all, then up all)
migrate-reset:
	@echo "$(COLOR_GREEN)Resetting database...$(COLOR_RESET)"
	cd internal/infra/db/migrations && goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=gohexaclean sslmode=disable" reset
	cd internal/infra/db/migrations && goose postgres "host=localhost port=5432 user=postgres password=postgres dbname=gohexaclean sslmode=disable" up

## seed: Run database seeders
seed:
	@echo "$(COLOR_GREEN)Running seeders...$(COLOR_RESET)"
	@for file in internal/infra/db/seeders/*.sql; do \
		echo "Running seeder: $$file"; \
		PGPASSWORD=postgres psql -h localhost -p 5432 -U postgres -d gohexaclean -f $$file; \
	done
	@echo "$(COLOR_GREEN)Seeders complete!$(COLOR_RESET)"

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

##@ Mocks

## mock-gen: Generate mocks for testing
mock-gen:
	@echo "$(COLOR_GREEN)Generating mocks...$(COLOR_RESET)"
	mockgen -source=internal/port/outbound/repository/user_repository.go -destination=internal/port/outbound/repository/mock/mock_user_repository.go -package=mock
	mockgen -source=internal/port/outbound/service/cache_service.go -destination=internal/port/outbound/service/mock/mock_cache_service.go -package=mock
	mockgen -source=internal/port/inbound/user_service_port.go -destination=internal/port/inbound/mock/mock_user_service.go -package=mock
	@echo "$(COLOR_GREEN)Mocks generated successfully!$(COLOR_RESET)"

##@ Utilities

## install-tools: Install development tools
install-tools:
	@echo "$(COLOR_GREEN)Installing development tools...$(COLOR_RESET)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "$(COLOR_GREEN)Tools installed!$(COLOR_RESET)"

## version: Show version information
version:
	@echo "Go version: $$(go version)"
	@echo "App: $(APP_NAME)"

.DEFAULT_GOAL := help
