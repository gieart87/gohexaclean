# API-First Development Workflow

## 📋 Overview

This project follows **API-First** (Contract-First) development approach using OpenAPI 3.0 specification. This ensures:

- ✅ **Design before implementation** - API contract is defined first
- ✅ **Documentation as code** - Spec is the single source of truth
- ✅ **Client-Server agreement** - Clear contract between frontend and backend
- ✅ **Auto-generated docs** - Swagger UI from spec
- ✅ **Validation** - Request/response validation against spec

## 🎯 Workflow

### Step 1: Design API in OpenAPI Spec

Create or update OpenAPI YAML file in `api/openapi/`:

```yaml
# api/openapi/user-api.yaml
openapi: 3.0.3

paths:
  /users:
    post:
      summary: Create user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserResponse'

components:
  schemas:
    CreateUserRequest:
      type: object
      required:
        - email
        - name
        - password
      properties:
        email:
          type: string
          format: email
        # ...
```

### Step 2: Validate OpenAPI Spec

Use OpenAPI tools to validate your spec:

```bash
# Using openapi-cli
npm install -g @redocly/cli
redocly lint api/openapi/user-api.yaml

# Or use online validator
# https://editor.swagger.io/
```

### Step 3: View Documentation

Start the HTTP server and access Swagger UI:

```bash
# Run server
make run-http

# Open browser
open http://localhost:8080/api/v1/swagger
```

### Step 4: Generate Code (Optional)

You can generate server stubs or client SDKs from the spec:

```bash
# Generate Go server stubs (oapi-codegen)
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

oapi-codegen -package generated -generate types,server,spec \
  api/openapi/user-api.yaml > internal/adapter/inbound/http/generated/api.go

# Generate TypeScript client
openapi-generator-cli generate \
  -i api/openapi/user-api.yaml \
  -g typescript-axios \
  -o clients/typescript
```

### Step 5: Implement Handler

Implement the handler following the spec:

```go
// internal/adapter/inbound/http/handler/user_handler.go

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // Parse request according to spec
    var req request.CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Invalid request body", err),
        )
    }

    // Call use case
    user, err := h.userService.CreateUser(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Failed to create user", err),
        )
    }

    // Return response according to spec
    return c.Status(fiber.StatusCreated).JSON(
        response.NewSuccessResponse("User created successfully", user),
    )
}
```

### Step 6: Test Against Spec

Ensure your implementation matches the spec:

```bash
# Manual testing with curl
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
  }'

# Or use Postman with OpenAPI import
# Import api/openapi/user-api.yaml to Postman

# Or use automated API testing
npm install -g dredd
dredd api/openapi/user-api.yaml http://localhost:8080
```

## 📁 File Organization

```
api/
├── openapi/              # OpenAPI specifications
│   ├── user-api.yaml    # User API spec
│   ├── product-api.yaml # Product API spec (future)
│   └── common.yaml      # Shared schemas
└── proto/               # gRPC proto files (different protocol)

docs/
└── swagger/             # Generated Swagger docs
```

## 🔧 Tools & Resources

### Recommended Tools

1. **OpenAPI Editor**
   - [Swagger Editor](https://editor.swagger.io/) - Online editor
   - [Stoplight Studio](https://stoplight.io/studio) - Desktop app
   - VSCode Extension: `openapi-lint`

2. **Code Generation**
   - [oapi-codegen](https://github.com/deepmap/oapi-codegen) - Go server/client
   - [OpenAPI Generator](https://openapi-generator.tech/) - Multi-language
   - [swagger-codegen](https://swagger.io/tools/swagger-codegen/)

3. **Validation & Testing**
   - [Redocly CLI](https://redocly.com/docs/cli/) - Lint & bundle
   - [Spectral](https://stoplight.io/open-source/spectral) - Linting
   - [Dredd](https://dredd.org/) - API testing
   - [Prism](https://stoplight.io/open-source/prism) - Mock server

4. **Documentation**
   - [ReDoc](https://redocly.github.io/redoc/) - Alternative UI
   - [Swagger UI](https://swagger.io/tools/swagger-ui/) - Built-in
   - [Postman](https://www.postman.com/) - Import & test

### Integration with CI/CD

Add to `.github/workflows/ci.yml`:

```yaml
- name: Validate OpenAPI Spec
  run: |
    npm install -g @redocly/cli
    redocly lint api/openapi/*.yaml

- name: Test API against spec
  run: |
    # Start server in background
    ./bin/http-server &
    sleep 5

    # Run API tests
    npm install -g dredd
    dredd api/openapi/user-api.yaml http://localhost:8080
```

## 📝 Best Practices

### 1. Versioning

Always version your APIs:

```yaml
servers:
  - url: http://localhost:8080/api/v1
    description: Version 1
  - url: http://localhost:8080/api/v2
    description: Version 2 (future)
```

### 2. Reusable Components

Define common schemas in `components`:

```yaml
components:
  schemas:
    Error:
      type: object
      properties:
        message:
          type: string
        code:
          type: string

    Pagination:
      type: object
      properties:
        page:
          type: integer
        limit:
          type: integer
        total:
          type: integer
```

### 3. Security Schemes

Always define security:

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

security:
  - BearerAuth: []
```

### 4. Examples

Provide examples for all schemas:

```yaml
CreateUserRequest:
  type: object
  properties:
    email:
      type: string
      example: user@example.com  # ← Example
```

### 5. Descriptions

Document everything:

```yaml
paths:
  /users:
    post:
      summary: Create user
      description: |
        Creates a new user account. Email must be unique.
        Password will be hashed before storage.
```

## 🚀 Adding New Endpoints

### Step-by-Step Guide

1. **Add to OpenAPI spec**:
```yaml
# api/openapi/user-api.yaml
paths:
  /users/{id}/activate:
    post:
      summary: Activate user
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: User activated
```

2. **Update DTOs** (if needed):
```go
// internal/dto/request/user_request.go
type ActivateUserRequest struct {
    Reason string `json:"reason"`
}
```

3. **Update use case**:
```go
// internal/port/inbound/user_service_port.go
type UserServicePort interface {
    // ... existing methods
    ActivateUser(ctx context.Context, id uuid.UUID) error
}

// internal/app/user_service.go
func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) error {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return err
    }

    user.Activate()
    return s.userRepo.Update(ctx, user)
}
```

4. **Implement handler**:
```go
// internal/adapter/inbound/http/handler/user_handler.go
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
    id, _ := uuid.Parse(c.Params("id"))

    if err := h.userService.ActivateUser(c.Context(), id); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(
            response.NewErrorResponse("Failed to activate user", err),
        )
    }

    return c.JSON(
        response.NewSuccessResponse("User activated", nil),
    )
}
```

5. **Add route**:
```go
// internal/adapter/inbound/http/router/router.go
protected.Post("/:id/activate", userHandler.ActivateUser)
```

6. **Test in Swagger UI**:
```
Open: http://localhost:8080/api/v1/swagger
Test the new endpoint
```

## 📊 Comparing Approaches

### ❌ Code-First (Traditional)
```
1. Write handler code
2. Implement endpoint
3. Write tests
4. Document API (often outdated)
```

**Problems**:
- Documentation often lags behind
- No clear contract
- Hard to coordinate with frontend

### ✅ API-First (Recommended)
```
1. Design API in OpenAPI spec
2. Review & validate spec
3. Generate client/server code (optional)
4. Implement handlers
5. Documentation auto-generated
```

**Benefits**:
- Spec is single source of truth
- Frontend can start before backend ready
- Always up-to-date documentation
- Better team collaboration

## 🎓 Learning Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [API Design Guide](https://cloud.google.com/apis/design)
- [Best Practices](https://swagger.io/resources/articles/best-practices-in-api-design/)
- [OpenAPI Tutorial](https://oai.github.io/Documentation/)

## 🔗 Quick Links

- **Swagger UI**: http://localhost:8080/api/v1/swagger
- **OpenAPI Spec**: http://localhost:8080/api/v1/swagger/spec
- **Editor**: https://editor.swagger.io/

---

**Remember**: Spec first, code second! 🎯
