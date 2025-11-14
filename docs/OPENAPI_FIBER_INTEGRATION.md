# OpenAPI + Fiber Integration Guide

## 📦 Generated Code

Setelah menjalankan `make openapi`, file berikut akan di-generate:

```
internal/adapter/inbound/http/generated/
└── server.gen.go    # All-in-one: Types, ServerInterface (Fiber), Spec, RegisterHandlers
```

**Satu file berisi semuanya:**
1. **Request/Response Types** - Models dari OpenAPI spec
2. **ServerInterface** - Interface dengan method signature menggunakan `*fiber.Ctx`
3. **RegisterHandlers** - Fungsi untuk auto-register routes ke Fiber
4. **GetSwagger()** - OpenAPI spec embedded in Go code

> **Note:** Jika sebelumnya ada `types.gen.go` atau `spec.gen.go` terpisah, hapus file-file tersebut. Sekarang semuanya ada di satu file `server.gen.go`.

## 🎯 Two Integration Options

### Option 1: Manual Routing (Current - Recommended)

**Keuntungan:**
- ✅ Full control atas routing dan middleware
- ✅ Mudah customize per-endpoint
- ✅ Tidak perlu refactor existing handlers
- ✅ Tetap mendapat type safety dari generated types

**Cara pakai:**

```go
// internal/adapter/inbound/http/handler/user_handler.go

import (
    "github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated"
    "github.com/gofiber/fiber/v2"
)

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // 1. Parse request menggunakan generated type
    var req generated.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid request body", err),
        )
    }

    // 2. Validation otomatis dari OpenAPI:
    //    - req.Email: valid email format
    //    - req.Name: min 3 chars, max 100 chars
    //    - req.Password: min 6 chars

    // 3. Convert ke domain DTO
    user, err := h.userService.CreateUser(c.Context(), &dto.CreateUserRequest{
        Email:    string(req.Email),
        Name:     req.Name,
        Password: req.Password,
    })

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Failed to create user", err),
        )
    }

    // 4. Return response
    return c.Status(fiber.StatusCreated).JSON(
        response.NewSuccessResponse("User created successfully", user),
    )
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
    // Parse query params menggunakan generated type
    var params generated.ListUsersParams
    if err := c.QueryParser(&params); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid query parameters", err),
        )
    }

    page := 1
    if params.Page != nil {
        page = *params.Page
    }

    limit := 10
    if params.Limit != nil {
        limit = *params.Limit
    }

    users, total, err := h.userService.ListUsers(c.Context(), page, limit)
    // ... implementation
}
```

**Router (manual):**

```go
// internal/adapter/inbound/http/router/router.go

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, ...) {
    api := app.Group("/api/v1")

    // Auth routes
    auth := api.Group("/auth")
    auth.Post("/login", userHandler.Login)

    // User routes
    users := api.Group("/users")
    users.Post("/", userHandler.CreateUser)

    // Protected routes
    protected := users.Group("")
    protected.Use(middleware.AuthMiddleware(tokenService))
    protected.Get("/", userHandler.ListUsers)
    protected.Get("/:id", userHandler.GetUser)
    protected.Put("/:id", userHandler.UpdateUser)
    protected.Delete("/:id", userHandler.DeleteUser)
}
```

---

### Option 2: Auto-Generated Routing (Alternative)

**Keuntungan:**
- ✅ Routes auto-generated dari OpenAPI spec
- ✅ Guarantee bahwa semua endpoints dari spec terimplementasi
- ✅ Consistent routing pattern

**Kelemahan:**
- ⚠️ Less flexibility untuk custom middleware per-endpoint
- ⚠️ Requires handler refactoring

**Cara pakai:**

#### Step 1: Implement ServerInterface

```go
// internal/adapter/inbound/http/handler/openapi_handler.go

package handler

import (
    "github.com/gieart87/gohexaclean/internal/adapter/inbound/http/generated"
    "github.com/gieart87/gohexaclean/internal/adapter/inbound/http/response"
    "github.com/gieart87/gohexaclean/internal/dto/request"
    "github.com/gieart87/gohexaclean/internal/port/inbound"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    openapi_types "github.com/oapi-codegen/runtime/types"
)

// OpenAPIHandler implements generated.ServerInterface
type OpenAPIHandler struct {
    userService inbound.UserServicePort
}

func NewOpenAPIHandler(userService inbound.UserServicePort) *OpenAPIHandler {
    return &OpenAPIHandler{
        userService: userService,
    }
}

// Implement ServerInterface methods

func (h *OpenAPIHandler) HealthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status":  "ok",
        "message": "Service is running",
    })
}

func (h *OpenAPIHandler) Login(c *fiber.Ctx) error {
    var req generated.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid request body", err),
        )
    }

    token, user, err := h.userService.Login(c.Context(), &request.LoginRequest{
        Email:    string(req.Email),
        Password: req.Password,
    })

    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(
            response.NewErrorResponse("Invalid credentials", err),
        )
    }

    return c.JSON(response.NewSuccessResponse("Login successful", fiber.Map{
        "token": token,
        "user":  user,
    }))
}

func (h *OpenAPIHandler) CreateUser(c *fiber.Ctx) error {
    var req generated.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid request body", err),
        )
    }

    user, err := h.userService.CreateUser(c.Context(), &request.CreateUserRequest{
        Email:    string(req.Email),
        Name:     req.Name,
        Password: req.Password,
    })

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Failed to create user", err),
        )
    }

    return c.Status(fiber.StatusCreated).JSON(
        response.NewSuccessResponse("User created successfully", user),
    )
}

func (h *OpenAPIHandler) ListUsers(c *fiber.Ctx, params generated.ListUsersParams) error {
    page := 1
    if params.Page != nil {
        page = *params.Page
    }

    limit := 10
    if params.Limit != nil {
        limit = *params.Limit
    }

    users, total, err := h.userService.ListUsers(c.Context(), page, limit)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(
            response.NewErrorResponse("Failed to list users", err),
        )
    }

    totalPages := (int(total) + limit - 1) / limit

    return c.JSON(response.NewPaginatedResponse(
        "Users retrieved successfully",
        users,
        page,
        limit,
        total,
        totalPages,
    ))
}

func (h *OpenAPIHandler) GetUserById(c *fiber.Ctx, id openapi_types.UUID) error {
    user, err := h.userService.GetUserByID(c.Context(), uuid.UUID(id))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(
            response.NewErrorResponse("User not found", err),
        )
    }

    return c.JSON(response.NewSuccessResponse("User retrieved successfully", user))
}

func (h *OpenAPIHandler) UpdateUser(c *fiber.Ctx, id openapi_types.UUID) error {
    var req generated.UpdateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid request body", err),
        )
    }

    user, err := h.userService.UpdateUser(c.Context(), uuid.UUID(id), &request.UpdateUserRequest{
        Name: req.Name,
    })

    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Failed to update user", err),
        )
    }

    return c.JSON(response.NewSuccessResponse("User updated successfully", user))
}

func (h *OpenAPIHandler) DeleteUser(c *fiber.Ctx, id openapi_types.UUID) error {
    if err := h.userService.DeleteUser(c.Context(), uuid.UUID(id)); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Failed to delete user", err),
        )
    }

    return c.JSON(response.NewSuccessResponse("User deleted successfully", nil))
}
```

#### Step 2: Register Routes using RegisterHandlers

```go
// internal/adapter/inbound/http/router/router.go

func SetupRoutes(
    app *fiber.App,
    userService inbound.UserServicePort,
    tokenService service.TokenService,
    log *logger.Logger,
) {
    api := app.Group("/api/v1")

    // Swagger documentation
    swaggerHandler := handler.NewSwaggerHandler()
    api.Get("/swagger", swaggerHandler.ServeSwaggerUI)
    api.Get("/swagger/spec", func(c *fiber.Ctx) error {
        spec, err := os.ReadFile("api/openapi/user-api.yaml")
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": "Failed to load API specification",
            })
        }
        c.Set("Content-Type", "application/x-yaml")
        return c.Send(spec)
    })

    // Auto-register all routes from OpenAPI spec
    openAPIHandler := handler.NewOpenAPIHandler(userService)

    // Option A: Basic registration
    generated.RegisterHandlers(api, openAPIHandler)

    // Option B: With middleware
    generated.RegisterHandlersWithOptions(api, openAPIHandler, generated.FiberServerOptions{
        BaseURL: "",
        Middlewares: []generated.MiddlewareFunc{
            // Global middlewares for all OpenAPI routes
            middleware.AuthMiddleware(tokenService),
        },
    })
}
```

## 🔄 Migration Path

Jika ingin migrate dari **Option 1** ke **Option 2**:

1. Create `openapi_handler.go` yang implements `generated.ServerInterface`
2. Copy logic dari existing handlers ke methods di `OpenAPIHandler`
3. Update `router.go` untuk gunakan `RegisterHandlers`
4. Remove manual route definitions
5. Test semua endpoints masih berfungsi

Namun **tidak perlu migrate** jika sudah nyaman dengan Option 1!

## 📝 Generated Types Reference

### Request Types

```go
// Login request
type LoginRequest struct {
    Email    openapi_types.Email `json:"email"`
    Password string              `json:"password"`
}

// Create user request
type CreateUserRequest struct {
    Email    openapi_types.Email `json:"email"`
    Name     string              `json:"name"`
    Password string              `json:"password"`
}

// Update user request
type UpdateUserRequest struct {
    Name string `json:"name"`
}

// Query params for list users
type ListUsersParams struct {
    Page  *int `form:"page,omitempty" json:"page,omitempty"`
    Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}
```

### Response Types

```go
// Single user response
type UserResponse struct {
    Data      *User      `json:"data,omitempty"`
    Message   *string    `json:"message,omitempty"`
    Success   *bool      `json:"success,omitempty"`
    Timestamp *time.Time `json:"timestamp,omitempty"`
}

// Paginated users response
type PaginatedUserResponse struct {
    Data    *[]User `json:"data,omitempty"`
    Message *string `json:"message,omitempty"`
    Meta    *struct {
        Limit      *int   `json:"limit,omitempty"`
        Page       *int   `json:"page,omitempty"`
        Total      *int64 `json:"total,omitempty"`
        TotalPages *int   `json:"total_pages,omitempty"`
    } `json:"meta,omitempty"`
    Success   *bool      `json:"success,omitempty"`
    Timestamp *time.Time `json:"timestamp,omitempty"`
}

// Error response
type ErrorResponse struct {
    Error     *string                 `json:"error,omitempty"`
    Errors    *map[string]interface{} `json:"errors"`
    Message   *string                 `json:"message,omitempty"`
    Success   *bool                   `json:"success,omitempty"`
    Timestamp *time.Time              `json:"timestamp,omitempty"`
}

// User model
type User struct {
    Id        *openapi_types.UUID  `json:"id,omitempty"`
    Email     *openapi_types.Email `json:"email,omitempty"`
    Name      *string              `json:"name,omitempty"`
    IsActive  *bool                `json:"is_active,omitempty"`
    CreatedAt *time.Time           `json:"created_at,omitempty"`
    UpdatedAt *time.Time           `json:"updated_at,omitempty"`
}
```

## ✅ Best Practices

1. **Use Generated Types** - Selalu gunakan generated types untuk request/response
2. **Validate in OpenAPI** - Definisikan validation rules di OpenAPI spec, bukan di code
3. **Single Source of Truth** - OpenAPI spec adalah contract, jangan modify generated code
4. **Regenerate After Spec Changes** - Selalu run `make openapi` setelah update spec
5. **Check Interface** - ServerInterface berguna untuk ensure semua endpoints terimplementasi

## 🎓 When to Use Each Option

### Use Option 1 (Manual Routing) when:
- Existing handlers sudah banyak dan kompleks
- Butuh fine-grained control per endpoint
- Custom middleware berbeda per route
- Team sudah familiar dengan manual routing

### Use Option 2 (Auto-Generated) when:
- Starting new project from scratch
- Want strict adherence to OpenAPI contract
- Prefer convention over configuration
- Want guaranteed coverage of all spec endpoints

## 🔗 Related Docs

- [API_FIRST_WORKFLOW.md](API_FIRST_WORKFLOW.md) - API-First development workflow
- [SWAGGER_GUIDE.md](SWAGGER_GUIDE.md) - Swagger UI usage guide
- [oapi-codegen Fiber docs](https://github.com/oapi-codegen/oapi-codegen)

---

**Happy Coding with OpenAPI + Fiber!** 🚀
