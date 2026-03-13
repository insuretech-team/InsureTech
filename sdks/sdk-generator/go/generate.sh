#!/bin/bash
# Generate Go SDK from OpenAPI spec
# This script generates the InsureTech Go SDK

set -e

echo "========================================"
echo "  InsureTech Go SDK Generator"
echo "========================================"
echo ""

# Get the script directory and project root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/../../.." && pwd )"

# Configuration - using relative paths from project root
PROTO_PATH="$PROJECT_ROOT/proto"
API_SPEC_PATH="$PROJECT_ROOT/api/openapi.yaml"
OUTPUT_PATH="$PROJECT_ROOT/sdks/insuretech-go-sdk"
GENERATOR_PATH="$SCRIPT_DIR"

echo "Project Root: $PROJECT_ROOT"
echo "Script Directory: $SCRIPT_DIR"
echo ""

# Validate paths
echo "Validating paths..."

if [ ! -d "$PROTO_PATH" ]; then
    echo "✗ Proto path not found: $PROTO_PATH"
    exit 1
fi
echo "✓ Proto path exists"

if [ ! -f "$API_SPEC_PATH" ]; then
    echo "✗ OpenAPI spec not found: $API_SPEC_PATH"
    exit 1
fi
echo "✓ OpenAPI spec exists"

echo ""

# Create output directory if it doesn't exist
if [ ! -d "$OUTPUT_PATH" ]; then
    echo "Creating output directory: $OUTPUT_PATH"
    mkdir -p "$OUTPUT_PATH"
    echo "✓ Output directory created"
else
    echo "✓ Output directory exists"
fi

echo ""

# Run the generator
echo "Running Go SDK generator..."
echo ""

cd "$GENERATOR_PATH"

# Check if generator is built
if [ ! -f "generator" ]; then
    echo "Building generator..."
    go build -o generator generator.go
    echo "✓ Generator built successfully"
fi

# Run the generator
./generator

echo ""
echo "========================================"
echo "✓ SDK Generation Complete!"
echo "========================================"
echo ""
echo "Output location: $OUTPUT_PATH"
echo ""
echo "Next steps:"
echo "  1. cd $OUTPUT_PATH"
echo "  2. go mod tidy"
echo "  3. go test ./..."
echo ""
