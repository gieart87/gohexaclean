#!/bin/bash

# Generate Go code from OpenAPI specifications

echo "Generating Go code from OpenAPI specs..."

# Check if oapi-codegen is installed
if ! command -v oapi-codegen &> /dev/null; then
    echo "Installing oapi-codegen..."
    go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

# Output directories
HEALTH_OUTPUT_DIR="internal/adapter/inbound/http/generated/healthapi"
USER_OUTPUT_DIR="internal/adapter/inbound/http/generated/userapi"
mkdir -p "$HEALTH_OUTPUT_DIR"
mkdir -p "$USER_OUTPUT_DIR"

# Generate Health API (health handler)
echo "Generating Health API..."
oapi-codegen -package healthapi -generate types,fiber,spec \
  api/openapi/health-api.yaml > "$HEALTH_OUTPUT_DIR/server.gen.go"

# Generate User API (user handlers)
echo "Generating User API..."
oapi-codegen -package userapi -generate types,fiber,spec \
  api/openapi/user-api.yaml > "$USER_OUTPUT_DIR/server.gen.go"

echo ""
echo "✅ OpenAPI code generation complete!"
echo ""
echo "Generated files:"
echo "  - $HEALTH_OUTPUT_DIR/server.gen.go    (Health ServerInterface)"
echo "  - $USER_OUTPUT_DIR/server.gen.go      (User ServerInterface)"
echo ""
echo "📝 Note:"
echo "   Each package contains its own ServerInterface:"
echo ""
echo "   package healthapi.ServerInterface:"
echo "     - HealthCheck(c *fiber.Ctx) error"
echo ""
echo "   package userapi.ServerInterface:"
echo "     - Login(c *fiber.Ctx) error"
echo "     - CreateUser(c *fiber.Ctx) error"
echo "     - ListUsers(c *fiber.Ctx, params ListUsersParams) error"
echo "     - GetUserById(c *fiber.Ctx, id UUID) error"
echo "     - UpdateUser(c *fiber.Ctx, id UUID) error"
echo "     - DeleteUser(c *fiber.Ctx, id UUID) error"
echo ""
echo "   Implement each interface in separate handler packages!"
