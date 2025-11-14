# Architecture Documentation

## Overview

GoHexaClean implements a **hybrid architecture** combining **Hexagonal Architecture** (Ports & Adapters) with **Clean Architecture** principles. This approach provides:

- **Framework Independence**: Business logic doesn't depend on frameworks
- **Testability**: Easy to test at all levels
- **Flexibility**: Easy to swap implementations
- **Maintainability**: Clear separation of concerns

## Architecture Layers

### 1. Domain Layer (Core)

**Location**: `internal/domain/`

The innermost layer containing pure business logic and entities.

**Characteristics**:
- No dependencies on external libraries
- Contains domain entities and value objects
- Implements domain-specific business rules
- Defines domain errors

**Example**:
```go
// Domain entity with business logic
type User struct {
    ID        uuid.UUID
    Email     string
    Name      string
    Password  string
    IsActive  bool
}

func (u *User) Deactivate() {
    u.IsActive = false
    u.UpdatedAt = time.Now()
}
```

### 2. Application Layer (Use Cases)

**Location**: `internal/app/`

Orchestrates the flow of data and coordinates domain entities.

**Characteristics**:
- Implements application-specific business rules
- Coordinates between domain and infrastructure
- Depends on ports (interfaces), not implementations
- Transaction boundaries

**Example**:
```go
type UserService struct {
    userRepo     repository.UserRepository
    hashService  service.HashService
    tokenService service.TokenService
}

func (s *UserService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error) {
    // Application logic orchestrating domain and infrastructure
}
```

### 3. Ports (Interfaces)

**Location**: `internal/port/`

Define contracts between layers.

#### Inbound Ports (Driving)
**Location**: `internal/port/inbound/`

Interfaces that define what the application can do (driven by adapters).

```go
type UserServicePort interface {
    CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.UserResponse, error)
    GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error)
}
```

#### Outbound Ports (Driven)
**Location**: `internal/port/outbound/`

Interfaces that define external dependencies (drive adapters).

```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}
```

### 4. Adapters

**Location**: `internal/adapter/`

Concrete implementations of ports.

#### Inbound Adapters (Drivers)
**Location**: `internal/adapter/inbound/`

Handle external requests and drive the application.

- **HTTP Adapter** (`http/`): REST API using Fiber
- **gRPC Adapter** (`grpc/`): gRPC service
- Could add: CLI, GraphQL, Message Queue consumers

#### Outbound Adapters (Driven)
**Location**: `internal/adapter/outbound/`

Implement external dependencies.

- **PostgreSQL** (`pgsql/`): Database repository
- **Redis** (`redis/`): Cache service
- **Telemetry** (`telemetry/`): Token and hash services

### 5. Infrastructure

**Location**: `internal/infra/`

Shared infrastructure concerns.

- **Database**: Connection management, migrations
- **Config**: Configuration loading
- **Logger**: Structured logging
- **Cache**: Redis client setup

### 6. DTOs (Data Transfer Objects)

**Location**: `internal/dto/`

Objects for transferring data across boundaries.

- **Request**: API request models
- **Response**: API response models

Decouples domain entities from API contracts.

## Dependency Flow

```
Inbound Adapter → Inbound Port → Application → Outbound Port → Outbound Adapter
     (HTTP)          (Interface)    (UseCase)     (Interface)      (Database)
```

**Key Rule**: Dependencies point INWARD
- Outer layers depend on inner layers
- Inner layers are independent of outer layers
- Domain layer has ZERO external dependencies

## Communication Flow

### HTTP Request Example

```
1. HTTP Request
   ↓
2. HTTP Handler (Adapter)
   ↓
3. UserServicePort (Inbound Port - Interface)
   ↓
4. UserService (Application - Use Case)
   ↓
5. Domain Entities (Business Logic)
   ↓
6. UserRepository (Outbound Port - Interface)
   ↓
7. PostgreSQL Adapter (Outbound Adapter)
   ↓
8. Database
```

### gRPC Request Example

```
1. gRPC Request
   ↓
2. gRPC Handler (Adapter)
   ↓
3. UserServicePort (Same as HTTP!)
   ↓
4-8. Same flow as HTTP
```

Notice: HTTP and gRPC share the same business logic!

## Benefits of This Architecture

### 1. Framework Independence
You can switch from Fiber to Gin without touching business logic:
```go
// Only change the adapter
internal/adapter/inbound/http/  // Swap Fiber for Gin
internal/app/                   // NO CHANGES NEEDED
internal/domain/                // NO CHANGES NEEDED
```

### 2. Testability
Each layer can be tested independently:
```go
// Unit test - Domain
func TestUser_Deactivate(t *testing.T) {
    user := &User{IsActive: true}
    user.Deactivate()
    assert.False(t, user.IsActive)
}

// Integration test - Application (with mocks)
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo, ...)
    // Test use case
}

// E2E test - Adapter
func TestHTTPHandler_CreateUser(t *testing.T) {
    // Test full HTTP flow
}
```

### 3. Flexibility
Easy to add new adapters:
- Want to add GraphQL? Add a new inbound adapter
- Need Kafka consumer? Add a new inbound adapter
- Using MongoDB instead? Add a new outbound adapter

### 4. Maintainability
Clear boundaries make changes easier:
- Business logic changes: Modify domain/application layers
- API changes: Modify adapters/DTOs only
- Database changes: Modify outbound adapters only

## Design Patterns Used

### 1. Repository Pattern
Abstracts data access logic.

```go
// Port (interface)
type UserRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// Adapter (implementation)
type UserRepositoryPG struct {
    db *sql.DB
}
```

### 2. Dependency Injection
Inject dependencies via constructor.

```go
func NewUserService(
    userRepo repository.UserRepository,
    hashService service.HashService,
) UserServicePort {
    return &UserService{
        userRepo:    userRepo,
        hashService: hashService,
    }
}
```

### 3. Factory Pattern
Create entities with proper initialization.

```go
func NewUser(email, name, password string) *User {
    return &User{
        ID:        uuid.New(),
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
    }
}
```

### 4. Strategy Pattern
Swappable implementations via interfaces.

```go
// Different cache strategies
type CacheService interface {
    Get(key string) (string, error)
}

// Redis implementation
type RedisCacheService struct {}

// In-memory implementation
type MemoryCacheService struct {}
```

## Best Practices

### 1. Keep Domain Pure
```go
// ❌ BAD: Domain depends on framework
type User struct {
    gorm.Model  // Framework dependency!
}

// ✅ GOOD: Pure domain entity
type User struct {
    ID        uuid.UUID
    CreatedAt time.Time
}
```

### 2. Use Interfaces for Ports
```go
// ✅ Always define interfaces in ports
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
}
```

### 3. DTOs for Boundary Crossing
```go
// ❌ BAD: Exposing domain entity directly
func (h *Handler) CreateUser(c *fiber.Ctx) (*domain.User, error)

// ✅ GOOD: Using DTOs
func (h *Handler) CreateUser(c *fiber.Ctx) (*response.UserResponse, error)
```

### 4. One-Way Dependencies
```go
// ✅ Dependencies point inward
Application → Domain (OK)
Adapter → Port (OK)
Domain → Adapter (NOT ALLOWED!)
```

## Adding New Features

### Example: Adding Product Management

1. **Create Domain Entity**
   ```go
   // internal/domain/product.go
   type Product struct {
       ID    uuid.UUID
       Name  string
       Price float64
   }
   ```

2. **Define Outbound Port**
   ```go
   // internal/port/outbound/repository/product_repository.go
   type ProductRepository interface {
       Create(ctx context.Context, product *domain.Product) error
   }
   ```

3. **Define Inbound Port**
   ```go
   // internal/port/inbound/product_service_port.go
   type ProductServicePort interface {
       CreateProduct(ctx context.Context, req *request.CreateProductRequest) error
   }
   ```

4. **Implement Application Service**
   ```go
   // internal/app/product_service.go
   type ProductService struct {
       productRepo repository.ProductRepository
   }
   ```

5. **Implement Adapters**
   ```go
   // internal/adapter/outbound/pgsql/product_repository_pg.go
   type ProductRepositoryPG struct {}

   // internal/adapter/inbound/http/handler/product_handler.go
   type ProductHandler struct {}
   ```

6. **Wire in Container**
   ```go
   // internal/bootstrap/container.go
   container.ProductService = app.NewProductService(
       container.ProductRepository,
   )
   ```

## Common Questions

### Q: Why not just use MVC?
**A**: MVC couples business logic to frameworks. Our architecture keeps business logic independent, making it more testable and maintainable.

### Q: Isn't this over-engineering?
**A**: For simple CRUD apps, maybe. For production microservices that will evolve, this provides necessary structure and flexibility.

### Q: How do I handle transactions?
**A**: Implement a Unit of Work pattern in the application layer, or use database transactions in repositories.

### Q: Where do validations go?
**A**:
- **Input validation**: DTOs/Adapters
- **Business rules**: Domain entities
- **Complex validation**: Application service

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
