# GoHexaClean

> Production-ready Golang boilerplate with Hexagonal + Clean Architecture

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A modern, scalable, and maintainable Go microservice boilerplate combining **Hexagonal Architecture** (Ports & Adapters) with **Clean Architecture** principles. Built for production use with comprehensive support for both HTTP (Fiber/Gin) and gRPC protocols.

## Features

- ✅ **Hybrid Architecture**: Hexagonal + Clean Architecture
- ✅ **Dual Protocol Support**: HTTP (Fiber) & gRPC
- ✅ **API-First Development**: OpenAPI 3.0 specification with Swagger UI
- ✅ **Framework Agnostic**: Easy to switch between Fiber, Gin, Echo, or Chi
- ✅ **PostgreSQL + Redis**: Production-ready database setup
- ✅ **Dependency Injection**: Clean DI container pattern
- ✅ **Structured Logging**: Using Uber's Zap
- ✅ **JWT Authentication**: Built-in auth middleware
- ✅ **Docker Ready**: Multi-stage Dockerfile & docker-compose
- ✅ **Database Migrations**: SQL migration support
- ✅ **SOLID Principles**: Highly testable and maintainable
- ✅ **Configuration**: Environment-based config with YAML support
- ✅ **Observability Ready**: Telemetry, metrics, and health checks
- ✅ **Interactive API Docs**: Auto-generated Swagger documentation

## Architecture

This boilerplate implements a **hybrid architecture** combining the best of both worlds:

```
┌─────────────────────────────────────────────────────────────┐
│                     External World                           │
│         (HTTP, gRPC, CLI, Kafka, REST API, DB)              │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────▼───────────────┐
        │   Inbound Adapters (Drivers)  │
        │   - HTTP Handlers (Fiber)     │
        │   - gRPC Handlers             │
        │   - CLI Commands              │
        └───────────────┬───────────────┘
                        │
        ┌───────────────▼───────────────┐
        │    Inbound Ports (Interface)  │
        │    - Service Interfaces       │
        └───────────────┬───────────────┘
                        │
        ┌───────────────▼───────────────┐
        │   Application Layer (Core)    │
        │   - Use Cases / Services      │
        │   - Business Logic            │
        └───────────────┬───────────────┘
                        │
        ┌───────────────▼───────────────┐
        │      Domain Layer (Core)      │
        │      - Entities               │
        │      - Domain Logic           │
        │      - Domain Errors          │
        └───────────────┬───────────────┘
                        │
        ┌───────────────▼───────────────┐
        │   Outbound Ports (Interface)  │
        │   - Repository Interfaces     │
        │   - Service Interfaces        │
        └───────────────┬───────────────┘
                        │
        ┌───────────────▼───────────────┐
        │  Outbound Adapters (Driven)   │
        │  - PostgreSQL Repository      │
        │  - Redis Cache                │
        │  - External APIs              │
        └───────────────┬───────────────┘
                        │
┌───────────────────────▼─────────────────────────────────────┐
│                  Infrastructure                              │
│              (DB, Cache, Config, Logger)                    │
└─────────────────────────────────────────────────────────────┘
```

### Key Concepts

1. **Domain Layer**: Core business logic, framework-independent
2. **Application Layer**: Use cases implementing business workflows
3. **Ports**: Interfaces defining contracts (Inbound for driving, Outbound for driven)
4. **Adapters**: Concrete implementations (HTTP, gRPC, Database, etc.)
5. **Infrastructure**: Shared concerns (Config, Logging, DB connection)

## Project Structure

```
gohexaclean/
├── api/
│   ├── openapi/                    # OpenAPI 3.0 specifications
│   │   ├── health-api.yaml        # Health check API spec
│   │   └── user-api.yaml          # User management API spec
│   └── proto/                      # Protocol Buffer definitions
│       └── user.proto
├── cmd/
│   ├── http/                       # HTTP server entry point
│   │   └── main.go
│   └── grpc/                       # gRPC server entry point
│       └── main.go
├── config/                         # Configuration files
│   └── app.yaml
├── docs/                           # Documentation
│   ├── API_FIRST_WORKFLOW.md
│   ├── SWAGGER_GUIDE.md
│   └── OPENAPI_FIBER_INTEGRATION.md
├── scripts/                        # Build and utility scripts
│   ├── generate-openapi.sh        # Auto-generate from OpenAPI specs
│   └── generate-proto.sh          # Generate from protobuf
├── internal/
│   ├── domain/                     # Domain entities & business logic
│   │   ├── user.go
│   │   └── errors.go
│   ├── dto/                        # Data Transfer Objects
│   │   ├── request/
│   │   └── response/
│   ├── app/                        # Use cases / Application services
│   │   └── user_service.go
│   ├── port/                       # Interfaces (Ports)
│   │   ├── inbound/               # Input ports (driven by adapters)
│   │   │   └── user_service_port.go
│   │   └── outbound/              # Output ports (drive adapters)
│   │       ├── repository/
│   │       └── service/
│   ├── adapter/                    # Adapters implementation
│   │   ├── inbound/
│   │   │   ├── http/              # HTTP adapter (Fiber)
│   │   │   │   ├── generated/     # Auto-generated from OpenAPI
│   │   │   │   │   ├── healthapi/ # Health API generated code
│   │   │   │   │   └── userapi/   # User API generated code
│   │   │   │   ├── middleware/
│   │   │   │   ├── handler/       # Per-endpoint handler files
│   │   │   │   │   ├── health/
│   │   │   │   │   │   ├── handler.go
│   │   │   │   │   │   └── health_check_handler.go
│   │   │   │   │   ├── user/
│   │   │   │   │   │   ├── handler.go
│   │   │   │   │   │   ├── login_handler.go
│   │   │   │   │   │   ├── register_handler.go
│   │   │   │   │   │   ├── admin_list_users_handler.go
│   │   │   │   │   │   ├── admin_get_user_handler.go
│   │   │   │   │   │   ├── admin_update_user_handler.go
│   │   │   │   │   │   └── admin_delete_user_handler.go
│   │   │   │   │   └── swagger_handler.go
│   │   │   │   └── router/
│   │   │   │       └── router.go
│   │   │   └── grpc/              # gRPC adapter
│   │   │       ├── handler/
│   │   │       └── interceptor/
│   │   └── outbound/
│   │       ├── pgsql/             # PostgreSQL adapter
│   │       ├── redis/             # Redis adapter
│   │       └── telemetry/         # Telemetry services
│   ├── infra/                      # Shared infrastructure
│   │   ├── db/
│   │   │   ├── connection.go
│   │   │   ├── migrations/
│   │   │   └── seeders/
│   │   ├── config/
│   │   ├── logger/
│   │   └── cache/
│   └── bootstrap/                  # Dependency injection
│       └── container.go
├── pkg/                            # Shared packages
│   ├── response/
│   ├── errors/
│   └── utils/
├── test/                           # Tests
│   ├── unit/
│   ├── integration/
│   ├── mocks/
│   └── fixtures/
├── scripts/                        # Utility scripts
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)
- Protocol Buffer Compiler (for gRPC)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/gieart87/gohexaclean.git
cd gohexaclean
```

2. **Install dependencies**

```bash
make deps
```

3. **Generate Protocol Buffers (for gRPC)**

```bash
make proto-install
make proto
```

4. **Copy environment file**

```bash
cp .env.example .env
```

5. **Update configuration**

Edit `.env` or `config/app.yaml` with your database credentials.

### Running with Docker

```bash
# Start all services (PostgreSQL, Redis, HTTP, gRPC)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Running Locally

1. **Start PostgreSQL and Redis**

```bash
docker-compose up -d postgres redis
```

2. **Run migrations**

```bash
make migrate-up
```

3. **Seed database (optional)**

```bash
make seed
```

4. **Run HTTP server**

```bash
make run-http
# Server runs on http://localhost:8080
```

5. **Run gRPC server (in another terminal)**

```bash
make run-grpc
# Server runs on localhost:50051
```

## API Documentation

### 📖 Interactive Swagger UI

Access the interactive API documentation at:

```
http://localhost:8080/api/v1/swagger
```

Features:
- **Try it out**: Test endpoints directly from browser
- **Authentication**: Test with JWT tokens
- **Request/Response examples**: See all schemas
- **OpenAPI 3.0 compliant**: Industry standard

### 📝 OpenAPI Specifications

View the raw OpenAPI specs at:

```
http://localhost:8080/api/v1/swagger/spec
```

Or find them in the repository:
- `api/openapi/health-api.yaml` - Health Check API
- `api/openapi/user-api.yaml` - User Management API

### 🔧 API-First Development Workflow

This boilerplate uses **OpenAPI-first** approach with auto-generated code:

```bash
# 1. Create/Edit OpenAPI spec
vim api/openapi/product-api.yaml

# 2. Generate Fiber handlers and types
make openapi

# 3. Implement handlers
# Generated: internal/adapter/inbound/http/generated/productapi/
# Implement: internal/adapter/inbound/http/handler/product/

# 4. Register routes
# Auto-register: productapi.RegisterHandlers(api, productHandler)
```

**Benefits:**
- ✅ Type-safe request/response handling
- ✅ Auto-validation from OpenAPI schema
- ✅ Auto-generated Swagger documentation
- ✅ Contract-first development
- ✅ One file per endpoint for better organization

**Generated Structure:**
```
handler/
├── health/
│   ├── handler.go              # Implements healthapi.ServerInterface
│   └── health_check_handler.go # GET /health
└── user/
    ├── handler.go                    # Implements userapi.ServerInterface
    ├── login_handler.go              # POST /auth/login (public)
    ├── register_handler.go           # POST /users (public)
    ├── admin_list_users_handler.go   # GET /users (protected)
    ├── admin_get_user_handler.go     # GET /users/{id} (protected)
    ├── admin_update_user_handler.go  # PUT /users/{id} (protected)
    └── admin_delete_user_handler.go  # DELETE /users/{id} (protected)
```

See [docs/API_FIRST_WORKFLOW.md](docs/API_FIRST_WORKFLOW.md) and [docs/OPENAPI_FIBER_INTEGRATION.md](docs/OPENAPI_FIBER_INTEGRATION.md) for detailed guides.

### Adding New API Endpoints

This project follows **API-First** approach. See [docs/API_FIRST_WORKFLOW.md](docs/API_FIRST_WORKFLOW.md) for:
- How to design APIs in OpenAPI
- Code generation from spec
- Testing against spec
- Best practices

### HTTP Endpoints

#### Health Check
```bash
GET /api/v1/health
```

#### Authentication
```bash
# Register
POST /api/v1/users
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "password123"
}

# Login
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### User Management (Protected)
```bash
# List users
GET /api/v1/users?page=1&limit=10
Authorization: Bearer <token>

# Get user
GET /api/v1/users/:id
Authorization: Bearer <token>

# Update user
PUT /api/v1/users/:id
Authorization: Bearer <token>
{
  "name": "Jane Doe"
}

# Delete user
DELETE /api/v1/users/:id
Authorization: Bearer <token>
```

### gRPC

Use [grpcurl](https://github.com/fullstorydev/grpcurl) to test gRPC endpoints:

```bash
# List services
grpcurl -plaintext localhost:50051 list

# Create user
grpcurl -plaintext -d '{"email":"test@example.com","name":"Test","password":"pass123"}' \
  localhost:50051 user.UserService/CreateUser

# Get user
grpcurl -plaintext -d '{"id":"<uuid>"}' \
  localhost:50051 user.UserService/GetUser
```

## Development

### Available Make Commands

```bash
make help                # Show all available commands
make run-http            # Run HTTP server
make run-grpc            # Run gRPC server
make build               # Build both servers
make test                # Run tests
make test-coverage       # Generate coverage report
make lint                # Run linter
make fmt                 # Format code
make proto               # Generate protobuf files
make docker-up           # Start docker containers
make docker-down         # Stop docker containers
make clean               # Clean build artifacts
```

### Switching HTTP Frameworks

The architecture is designed to be framework-agnostic. To switch from Fiber to Gin:

1. Update `internal/adapter/inbound/http/handler/` to use Gin's context
2. Update `internal/adapter/inbound/http/router/` to use Gin's router
3. Update `cmd/http/main.go` to initialize Gin instead of Fiber
4. Update dependencies in `go.mod`

The business logic (Domain, Application, Ports) remains unchanged!

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make test-coverage
```

## Configuration

Configuration can be provided via:

1. **YAML file** (`config/app.yaml`)
2. **Environment variables** (`.env` file)
3. **Environment-specific YAML** (`config/app.dev.yaml`, `config/app.prod.yaml`)

Environment variables take precedence over YAML configuration.

## Design Patterns Used

- **Hexagonal Architecture** (Ports & Adapters)
- **Clean Architecture** (Layer separation)
- **Repository Pattern** (Data access abstraction)
- **Dependency Injection** (Container pattern)
- **Factory Pattern** (Entity creation)
- **Strategy Pattern** (Swappable adapters)

## Best Practices

✅ **SOLID Principles**
✅ **Domain-Driven Design**
✅ **12-Factor App**
✅ **Dependency Injection**
✅ **Interface Segregation**
✅ **Error Handling**
✅ **Structured Logging**
✅ **Graceful Shutdown**

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Hexagonal Architecture (Alistair Cockburn)
- Clean Architecture principles (Robert C. Martin)
- Domain-Driven Design (Eric Evans)

## Support

If you find this project helpful, please give it a ⭐️!

For questions and support, please open an issue.

---

**Built with ❤️ by [gieart87](https://github.com/gieart87)**
