#!/bin/bash

# Generate Go code from OpenAPI specifications
# This script automatically detects all OpenAPI YAML files in api/openapi/
# and generates corresponding Go code for each

echo "Generating Go code from OpenAPI specs..."

# Check if oapi-codegen is installed
if ! command -v oapi-codegen &> /dev/null; then
    echo "Installing oapi-codegen..."
    go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
fi

# Base output directory
BASE_OUTPUT_DIR="internal/adapter/inbound/http/generated"

# Find all OpenAPI spec files
SPEC_DIR="api/openapi"
if [ ! -d "$SPEC_DIR" ]; then
    echo "Error: Directory $SPEC_DIR not found"
    exit 1
fi

# Counter for generated files
count=0

# Process each YAML file in the openapi directory
for spec_file in "$SPEC_DIR"/*.yaml; do
    if [ ! -f "$spec_file" ]; then
        echo "No OpenAPI spec files found in $SPEC_DIR"
        exit 1
    fi

    # Extract filename without extension (e.g., "user-api.yaml" -> "user-api")
    filename=$(basename "$spec_file" .yaml)

    # Convert to package name (e.g., "user-api" -> "userapi", "health-api" -> "healthapi")
    package_name=$(echo "$filename" | sed 's/-//g')

    # Create output directory
    output_dir="$BASE_OUTPUT_DIR/$package_name"
    mkdir -p "$output_dir"

    # Generate code
    echo "Generating $package_name from $filename..."
    oapi-codegen -package "$package_name" -generate types,fiber,spec \
        "$spec_file" > "$output_dir/server.gen.go"

    if [ $? -eq 0 ]; then
        echo "  ✓ Generated: $output_dir/server.gen.go"
        ((count++))
    else
        echo "  ✗ Failed to generate from $spec_file"
        exit 1
    fi
done

echo ""
echo "✅ OpenAPI code generation complete!"
echo "   Generated $count package(s) from $SPEC_DIR/*.yaml"
echo ""
echo "📝 Note:"
echo "   - Each .yaml file in $SPEC_DIR automatically generates a package"
echo "   - Package name derived from filename (e.g., user-api.yaml → userapi)"
echo "   - Generated code includes types, Fiber ServerInterface, and OpenAPI spec"
echo ""
echo "To add a new API:"
echo "   1. Create new YAML file in $SPEC_DIR (e.g., product-api.yaml)"
echo "   2. Run: make openapi"
echo "   3. Implement handler/product/handler.go with productapi.ServerInterface"
echo "   4. Register in router: productapi.RegisterHandlers(api, productHandler)"
