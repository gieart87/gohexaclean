# Quick Start Guide

## 🚀 Get Started in 5 Minutes

This guide will help you get GoHexaClean up and running quickly.

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional, but recommended)

## Option 1: Docker (Recommended for Quick Start)

### Step 1: Clone and Setup

```bash
# Clone repository
git clone https://github.com/gieart87/gohexaclean.git
cd gohexaclean

# Copy environment file
cp .env.example .env
```

### Step 2: Generate Protocol Buffers

```bash
# Install protobuf tools
make proto-install

# Generate proto files
./scripts/generate-proto.sh
```

### Step 3: Start Everything with Docker

```bash
# Start all services (PostgreSQL, Redis, HTTP, gRPC)
docker-compose up -d

# Check logs
docker-compose logs -f
```

That's it! Your services are now running:
- **HTTP Server**: http://localhost:8080
- **gRPC Server**: localhost:50051
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### Step 4: Test the API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Option 2: Local Development

### Step 1: Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install development tools
make install-tools

# Generate proto files
make proto
```

### Step 2: Start Infrastructure

```bash
# Start only PostgreSQL and Redis
docker-compose up -d postgres redis
```

### Step 3: Run Database Migrations

```bash
# Run migrations (manual for now)
psql -h localhost -U postgres -d gohexaclean -f internal/infra/db/migrations/001_create_users_table.sql

# Optional: Seed database
make seed
```

### Step 4: Run the Application

```bash
# Terminal 1: Run HTTP server
make run-http

# Terminal 2: Run gRPC server (optional)
make run-grpc
```

## Project Structure Overview

```
gohexaclean/
├── cmd/                    # Entry points (HTTP & gRPC)
├── internal/
│   ├── domain/            # Business entities
│   ├── app/               # Use cases
│   ├── port/              # Interfaces
│   ├── adapter/           # Implementations
│   ├── infra/             # Infrastructure
│   └── bootstrap/         # DI container
├── api/proto/             # gRPC definitions
├── config/                # Configuration
└── pkg/                   # Shared utilities
```

## Common Tasks

### Run Tests

```bash
make test
```

### Format Code

```bash
make fmt
```

### Build Binaries

```bash
make build
```

### View Logs

```bash
# Docker logs
docker-compose logs -f

# Specific service
docker-compose logs -f http-server
```

### Stop Services

```bash
docker-compose down
```

### Clean Everything

```bash
docker-compose down -v
make clean
```

## Testing gRPC

Install grpcurl:
```bash
# macOS
brew install grpcurl

# Linux
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

Test gRPC endpoints:
```bash
# List services
grpcurl -plaintext localhost:50051 list

# Create user via gRPC
grpcurl -plaintext -d '{
  "email":"grpc@example.com",
  "name":"gRPC User",
  "password":"password123"
}' localhost:50051 user.UserService/CreateUser
```

## Next Steps

1. **Read Architecture Docs**: [ARCHITECTURE.md](ARCHITECTURE.md)
2. **Explore API**: Try all endpoints in README
3. **Add Features**: Follow the architecture guide
4. **Write Tests**: See test examples in `test/`

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker-compose ps

# Restart database
docker-compose restart postgres
```

### Proto Generation Failed

```bash
# Install protoc first
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler

# Then install Go plugins
make proto-install
```

### Cannot Connect to Redis

Redis is optional. If you don't need caching, the app will work without it with a no-op cache service.

## Need Help?

- Check [README.md](../README.md) for detailed documentation
- See [ARCHITECTURE.md](ARCHITECTURE.md) for design decisions
- Open an issue on GitHub

Happy coding! 🎉
