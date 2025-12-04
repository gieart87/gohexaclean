# Error Handling Guide

This guide explains the error handling strategy in GoHexaClean, following Clean Architecture principles with proper separation of concerns.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Error Handling Layers                    │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  1. Domain Errors (internal/domain/errors.go)               │
│     └─> Business logic errors (pure domain)                 │
│                                                               │
│  2. Infrastructure Errors (internal/infra/<component>/)     │
│     ├─> db/errors.go: Database-specific errors              │
│     ├─> cache/errors.go: Cache-specific errors              │
│     ├─> broker/errors.go: Message broker errors             │
│     └─> asynq/errors.go: Task queue errors                  │
│                                                               │
│  3. HTTP Error Wrapper (pkg/errors/)                        │
│     ├─> errors.go: AppError struct & constructors           │
│     └─> mapper.go: All errors → HTTP status mapping         │
│                                                               │
│  4. Validation Errors (DTO layer)                           │
│     └─> Request validation using ozzo-validation            │
│                                                               │
│  5. Response Helpers (pkg/response/)                        │
│     └─> Standardized JSON error responses                   │
│                                                               │
└─────────────────────────────────────────────────────────────┘
```

## Best Practice Error Placement

| **Error Type** | **Location** | **Example** |
|----------------|--------------|-------------|
| Domain/Business Logic | `internal/domain/errors.go` | `ErrUserNotFound`, `ErrInvalidCredentials` |
| HTTP Status Mapping | `pkg/errors/errors.go` | `AppError`, `BadRequest()`, `NotFound()` |
| Error Mapper | `pkg/errors/mapper.go` | `MapDomainError()` |
| Validation | DTO layer | `CreateUserRequest.Validate()` |
| Database Infrastructure | `internal/infra/db/errors.go` | `ErrDBConnection`, `ErrDBTimeout` |
| Cache Infrastructure | `internal/infra/cache/errors.go` | `ErrCacheConnection`, `ErrCacheTimeout` |
| Message Broker | `internal/infra/broker/errors.go` | `ErrBrokerConnection`, `ErrBrokerPublish` |
| Task Queue (Asynq) | `internal/infra/asynq/errors.go` | `ErrTaskEnqueue`, `ErrTaskProcess` |
| Response Helpers | `pkg/response/` | `NewErrorResponse()`, `NewValidationErrorResponse()` |

## Error Types

### 1. Domain Errors

**Location:** `internal/domain/errors.go`

Domain errors represent business logic violations and are independent of any delivery mechanism (HTTP, gRPC, etc.).

```go
// User-related errors
ErrUserNotFound       = errors.New("user not found")
ErrUserAlreadyExists  = errors.New("user already exists")
ErrInvalidCredentials = errors.New("invalid credentials")

// Generic errors
ErrInvalidInput   = errors.New("invalid input")
ErrUnauthorized   = errors.New("unauthorized")
ErrForbidden      = errors.New("forbidden")
ErrInternalServer = errors.New("internal server error")
```

**When to use:**
- ✅ Business rule violations
- ✅ Domain entity state errors
- ✅ Use case/service layer errors
- ❌ NOT for validation errors (use DTO validation)
- ❌ NOT for infrastructure errors (use specific infra errors)

### 2. Infrastructure Errors

Infrastructure errors are specific to each infrastructure component and are located in their respective packages.

#### 2.1 Database Errors

**Location:** `internal/infra/db/errors.go`

```go
ErrDBConnection     = errors.New("database connection failed")
ErrDBTimeout        = errors.New("database operation timeout")
ErrDBTransaction    = errors.New("database transaction failed")
ErrDBMigration      = errors.New("database migration failed")
ErrDBRecordNotFound = errors.New("record not found in database")
ErrDBDuplicateKey   = errors.New("duplicate key violation")
ErrDBConstraint     = errors.New("database constraint violation")
```

#### 2.2 Cache Errors

**Location:** `internal/infra/cache/errors.go`

```go
ErrCacheConnection  = errors.New("cache connection failed")
ErrCacheTimeout     = errors.New("cache operation timeout")
ErrCacheKeyNotFound = errors.New("key not found in cache")
ErrCacheMarshal     = errors.New("failed to marshal cache data")
ErrCacheUnmarshal   = errors.New("failed to unmarshal cache data")
ErrCacheExpired     = errors.New("cache entry expired")
```

#### 2.3 Message Broker Errors

**Location:** `internal/infra/broker/errors.go`

```go
ErrBrokerConnection    = errors.New("broker connection failed")
ErrBrokerPublish       = errors.New("failed to publish message")
ErrBrokerSubscribe     = errors.New("failed to subscribe to topic")
ErrBrokerTimeout       = errors.New("broker operation timeout")
ErrBrokerChannelClosed = errors.New("broker channel closed")
ErrBrokerAck           = errors.New("failed to acknowledge message")
ErrBrokerNack          = errors.New("failed to negatively acknowledge message")
```

#### 2.4 Task Queue Errors (Asynq)

**Location:** `internal/infra/asynq/errors.go`

```go
ErrTaskEnqueue   = errors.New("failed to enqueue task")
ErrTaskProcess   = errors.New("failed to process task")
ErrTaskTimeout   = errors.New("task processing timeout")
ErrTaskRetry     = errors.New("task retry limit exceeded")
ErrTaskDuplicate = errors.New("duplicate task detected")
ErrWorkerStart   = errors.New("failed to start worker")
ErrWorkerStop    = errors.New("failed to stop worker")
```

**When to use infrastructure errors:**
- ✅ Connection failures
- ✅ Timeout issues
- ✅ Infrastructure-specific operations
- ✅ Technical errors that don't belong to business logic

### 3. HTTP Error Wrapper

**Location:** `pkg/errors/errors.go`

Wraps errors with HTTP status codes for HTTP/REST API responses.

```go
type AppError struct {
    Code    int    `json:"code"`    // HTTP status code
    Message string `json:"message"` // User-friendly message
    Err     error  `json:"-"`       // Original error (not exposed)
}
```

**Constructors:**
```go
BadRequest(message string, err error) *AppError           // 400
Unauthorized(message string, err error) *AppError         // 401
Forbidden(message string, err error) *AppError            // 403
NotFound(message string, err error) *AppError             // 404
Conflict(message string, err error) *AppError             // 409
InternalServerError(message string, err error) *AppError  // 500
```

### 3. Error Mapper

**Location:** `pkg/errors/mapper.go`

Automatically maps domain errors to HTTP status codes.

```go
// Automatic mapping
appErr := errors.MapDomainError(domainErr)

// With custom message
appErr := errors.MapDomainErrorWithCustomMessage(domainErr, "Custom message")

// Get status code only
statusCode := errors.GetHTTPStatusFromDomainError(domainErr)
```

**Mapping Table:**

| Domain Error | HTTP Status | Status Code |
|--------------|-------------|-------------|
| `ErrUserNotFound` | Not Found | 404 |
| `ErrUserAlreadyExists` | Conflict | 409 |
| `ErrInvalidCredentials` | Unauthorized | 401 |
| `ErrUnauthorized` | Unauthorized | 401 |
| `ErrForbidden` | Forbidden | 403 |
| `ErrInvalidInput` | Bad Request | 400 |
| Other errors | Internal Server Error | 500 |

### 4. Validation Errors

**Location:** `internal/dto/request/*.go`

Request DTOs handle input validation using `ozzo-validation`.

```go
// Example: CreateUserRequest
func (r CreateUserRequest) Validate() error {
    return validation.ValidateStruct(&r,
        validation.Field(&r.Email,
            validation.Required.Error("email is required"),
            is.Email.Error("email must be a valid email address"),
        ),
        validation.Field(&r.Name,
            validation.Required.Error("name is required"),
            validation.Length(3, 100).Error("name must be between 3 and 100 characters"),
        ),
        validation.Field(&r.Password,
            validation.Required.Error("password is required"),
            validation.Length(6, 0).Error("password must be at least 6 characters"),
        ),
    )
}
```

## Error Handling Patterns

### Pattern 1: Service Layer Returns Domain Errors

**Service Layer** (`internal/app/user_service.go`):
```go
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        // Return domain error
        return nil, domain.ErrUserNotFound
    }
    return response.NewUserResponse(user), nil
}
```

### Pattern 2: HTTP Handler Maps Domain Errors

**HTTP Handler** (simple approach):
```go
func (h *Handler) GetUser(c *fiber.Ctx) error {
    id, _ := uuid.Parse(c.Params("id"))

    user, err := h.userService.GetUserByID(c.Context(), id)
    if err != nil {
        // Manual mapping
        if errors.Is(err, domain.ErrUserNotFound) {
            return c.Status(fiber.StatusNotFound).JSON(
                response.NewErrorResponse("User not found", err),
            )
        }
        return c.Status(fiber.StatusInternalServerError).JSON(
            response.NewErrorResponse("Internal server error", err),
        )
    }

    return c.Status(fiber.StatusOK).JSON(
        response.NewSuccessResponse("User retrieved successfully", user),
    )
}
```

**HTTP Handler** (using mapper - recommended):
```go
func (h *Handler) GetUser(c *fiber.Ctx) error {
    id, _ := uuid.Parse(c.Params("id"))

    user, err := h.userService.GetUserByID(c.Context(), id)
    if err != nil {
        // Use error mapper
        appErr := pkgErrors.MapDomainError(err)
        return c.Status(appErr.Code).JSON(
            response.NewErrorResponse(appErr.Message, err),
        )
    }

    return c.Status(fiber.StatusOK).JSON(
        response.NewSuccessResponse("User retrieved successfully", user),
    )
}
```

### Pattern 3: Validation Errors

**HTTP Handler**:
```go
func (h *Handler) Register(c *fiber.Ctx) error {
    var req request.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid request body", err),
        )
    }

    // Validate request
    if err := req.Validate(); err != nil {
        return c.Status(fiber.StatusUnprocessableEntity).JSON(
            response.NewValidationErrorResponse(
                "Validation failed",
                response.ParseValidationErrors(err),
            ),
        )
    }

    user, err := h.userService.CreateUser(c.Context(), &req)
    if err != nil {
        appErr := pkgErrors.MapDomainError(err)
        return c.Status(appErr.Code).JSON(
            response.NewErrorResponse(appErr.Message, err),
        )
    }

    return c.Status(fiber.StatusCreated).JSON(
        response.NewSuccessResponse("User created successfully", user),
    )
}
```

## Response Format

### Success Response

```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

### Error Response

```json
{
  "success": false,
  "message": "User not found",
  "error": "user not found"
}
```

### Validation Error Response

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "email": "email is required",
    "password": "password must be at least 6 characters"
  }
}
```

## Best Practices

### 1. ✅ DO: Separate Concerns

```go
// ✅ GOOD: Domain layer returns domain errors
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, domain.ErrUserNotFound  // Domain error
    }
    return response.NewUserResponse(user), nil
}

// ✅ GOOD: HTTP layer handles HTTP concerns
func (h *Handler) GetUser(c *fiber.Ctx) error {
    user, err := h.userService.GetUserByID(c.Context(), id)
    if err != nil {
        appErr := pkgErrors.MapDomainError(err)  // Convert to HTTP error
        return c.Status(appErr.Code).JSON(...)
    }
    return c.Status(fiber.StatusOK).JSON(...)
}
```

### 2. ❌ DON'T: Mix HTTP Concerns in Domain

```go
// ❌ BAD: Don't return HTTP status codes from service
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (int, *response.UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return 404, nil, err  // ❌ HTTP status in service layer
    }
    return 200, response.NewUserResponse(user), nil
}
```

### 3. ✅ DO: Use Error Wrapping

```go
// ✅ GOOD: Wrap errors with context
func (s *UserService) CreateUser(ctx context.Context, req *request.CreateUserRequest) (*response.LoginResponse, error) {
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)  // ✅ Wrap with context
    }
    return response, nil
}
```

### 4. ✅ DO: Validate at the Boundary

```go
// ✅ GOOD: Validate in handler (boundary)
func (h *Handler) Register(c *fiber.Ctx) error {
    var req request.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(...)
    }

    // Validate at boundary
    if err := req.Validate(); err != nil {
        return c.Status(fiber.StatusUnprocessableEntity).JSON(...)
    }

    user, err := h.userService.CreateUser(c.Context(), &req)
    // ... handle response
}
```

### 5. ✅ DO: Use Sentinel Errors

```go
// ✅ GOOD: Use sentinel errors for comparison
var (
    ErrUserNotFound = errors.New("user not found")
)

// Check using errors.Is
if errors.Is(err, domain.ErrUserNotFound) {
    // Handle specific error
}
```

### 6. ❌ DON'T: Compare Error Strings

```go
// ❌ BAD: String comparison
if err.Error() == "user not found" {
    // Fragile and error-prone
}

// ✅ GOOD: Use errors.Is
if errors.Is(err, domain.ErrUserNotFound) {
    // Type-safe comparison
}
```

## Adding New Error Types

### Step 1: Add Domain Error

**File:** `internal/domain/errors.go`
```go
var (
    // ... existing errors
    ErrProductNotFound = errors.New("product not found")
    ErrInsufficientStock = errors.New("insufficient stock")
)
```

### Step 2: Update Error Mapper

**File:** `pkg/errors/mapper.go`
```go
func MapDomainError(err error) *AppError {
    switch {
    // ... existing mappings
    case errors.Is(err, domain.ErrProductNotFound):
        return NotFound("Product not found", err)
    case errors.Is(err, domain.ErrInsufficientStock):
        return BadRequest("Insufficient stock", err)
    default:
        return InternalServerError("Internal server error", err)
    }
}
```

### Step 3: Use in Service

**File:** `internal/app/product_service.go`
```go
func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*response.ProductResponse, error) {
    product, err := s.productRepo.FindByID(ctx, id)
    if err != nil {
        return nil, domain.ErrProductNotFound
    }
    return response.NewProductResponse(product), nil
}
```

### Step 4: Use in Handler

**File:** `internal/adapter/inbound/http/handler/product/get_product_handler.go`
```go
func (h *Handler) GetProduct(c *fiber.Ctx) error {
    id, _ := uuid.Parse(c.Params("id"))

    product, err := h.productService.GetProductByID(c.Context(), id)
    if err != nil {
        appErr := pkgErrors.MapDomainError(err)
        return c.Status(appErr.Code).JSON(
            response.NewErrorResponse(appErr.Message, err),
        )
    }

    return c.Status(fiber.StatusOK).JSON(
        response.NewSuccessResponse("Product retrieved successfully", product),
    )
}
```

## Error Handling for Different Protocols

### HTTP/REST (Fiber)

```go
func (h *Handler) GetUser(c *fiber.Ctx) error {
    user, err := h.userService.GetUserByID(c.Context(), id)
    if err != nil {
        appErr := pkgErrors.MapDomainError(err)
        return c.Status(appErr.Code).JSON(
            response.NewErrorResponse(appErr.Message, err),
        )
    }
    return c.Status(fiber.StatusOK).JSON(...)
}
```

### gRPC

```go
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, err := s.userService.GetUserByID(ctx, uuid.MustParse(req.Id))
    if err != nil {
        // Convert domain error to gRPC status
        if errors.Is(err, domain.ErrUserNotFound) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        return nil, status.Error(codes.Internal, "internal server error")
    }
    return &pb.GetUserResponse{User: toProtoUser(user)}, nil
}
```

## Testing Error Handling

### Unit Test Example

```go
func TestGetUserByID_NotFound(t *testing.T) {
    // Setup
    mockRepo := new(MockUserRepository)
    service := app.NewUserService(mockRepo, nil, nil, nil, nil)

    // Mock repository to return error
    mockRepo.On("FindByID", mock.Anything, mock.Anything).
        Return(nil, domain.ErrUserNotFound)

    // Execute
    user, err := service.GetUserByID(context.Background(), uuid.New())

    // Assert
    assert.Nil(t, user)
    assert.ErrorIs(t, err, domain.ErrUserNotFound)
}
```

### Integration Test Example

```go
func TestRegisterHandler_ValidationError(t *testing.T) {
    // Setup
    app := fiber.New()
    handler := setupHandler()
    app.Post("/auth/register", handler.Register)

    // Invalid request (missing required fields)
    reqBody := `{"email": "invalid-email"}`
    req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")

    // Execute
    resp, _ := app.Test(req)

    // Assert
    assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}
```

## Troubleshooting

### Problem: Errors not being mapped correctly

**Solution:** Ensure you're using `errors.Is()` for comparison, not direct equality.

```go
// ❌ BAD
if err == domain.ErrUserNotFound {
    // Won't work with wrapped errors
}

// ✅ GOOD
if errors.Is(err, domain.ErrUserNotFound) {
    // Works with wrapped errors
}
```

### Problem: Lost error context

**Solution:** Always wrap errors with context using `fmt.Errorf` with `%w`:

```go
// ✅ GOOD
if err := s.userRepo.Create(ctx, user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

### Problem: HTTP status codes in service layer

**Solution:** Keep service layer protocol-agnostic. Return domain errors only.

```go
// ❌ BAD - Service layer
return nil, fiber.NewError(fiber.StatusNotFound, "user not found")

// ✅ GOOD - Service layer
return nil, domain.ErrUserNotFound

// ✅ GOOD - Handler layer
appErr := pkgErrors.MapDomainError(err)
return c.Status(appErr.Code).JSON(...)
```

## Summary

1. **Domain errors** (`internal/domain/errors.go`) - Business logic errors
2. **HTTP wrapper** (`pkg/errors/errors.go`) - HTTP status code mapping
3. **Error mapper** (`pkg/errors/mapper.go`) - Automatic domain → HTTP conversion
4. **Validation** - DTO layer using ozzo-validation
5. **Response helpers** (`pkg/response/`) - Standardized JSON responses

**Key Principle:** Keep domain layer independent of delivery mechanisms. Convert errors at the adapter/handler layer.

---

For more information, see:
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Error Handling Best Practices](https://go.dev/blog/error-handling-and-go)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
