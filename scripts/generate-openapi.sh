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

# Generate types (models) - Framework agnostic
echo "Generating types..."
oapi-codegen -package generated -generate types \
  api/openapi/user-api.yaml > "$OUTPUT_DIR/types.gen.go"

# Generate spec (OpenAPI spec as Go code)
echo "Generating spec..."
oapi-codegen -package generated -generate spec \
  api/openapi/user-api.yaml > "$OUTPUT_DIR/spec.gen.go"

echo "✅ OpenAPI code generation complete!"
echo ""
echo "Generated files:"
echo "  - $OUTPUT_DIR/types.gen.go     (Request/Response types - Framework agnostic)"
echo "  - $OUTPUT_DIR/spec.gen.go      (OpenAPI spec embedded in Go)"
echo ""
echo "📝 Note:"
echo "   We only generate types and spec (framework-agnostic)."
echo "   Your Fiber handlers in internal/adapter/inbound/http/handler/"
echo "   are the actual implementation."
echo ""
echo "   Generated types can be used in your Fiber handlers for validation."
