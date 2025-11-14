#!/bin/bash

# Generate Go code from OpenAPI specification

echo "Generating Go code from OpenAPI spec..."

# Check if oapi-codegen is installed
if ! command -v oapi-codegen &> /dev/null; then
    echo "Installing oapi-codegen..."
    go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

# Output directory
OUTPUT_DIR="internal/adapter/inbound/http/generated"
mkdir -p "$OUTPUT_DIR"

# Generate types, server interface, and spec for Fiber
echo "Generating types, server interface, and spec for Fiber..."
oapi-codegen -package generated -generate types,fiber,spec \
  api/openapi/user-api.yaml > "$OUTPUT_DIR/server.gen.go"

echo "✅ OpenAPI code generation complete!"
echo ""
echo "Generated files:"
echo "  - $OUTPUT_DIR/server.gen.go    (Types, Fiber ServerInterface, and spec)"
echo ""
echo "📝 Note:"
echo "   Generated ServerInterface is designed for Fiber (uses *fiber.Ctx):"
echo ""
echo "   type ServerInterface interface {"
echo "       HealthCheck(c *fiber.Ctx) error"
echo "       Login(c *fiber.Ctx) error"
echo "       ListUsers(c *fiber.Ctx, params ListUsersParams) error"
echo "       CreateUser(c *fiber.Ctx) error"
echo "       GetUserById(c *fiber.Ctx, id string) error"
echo "       UpdateUser(c *fiber.Ctx, id string) error"
echo "       DeleteUser(c *fiber.Ctx, id string) error"
echo "   }"
echo ""
echo "   You can implement this interface in your handlers or use generated"
echo "   types directly. All requests are automatically validated against"
echo "   the OpenAPI spec!"
