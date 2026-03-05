#!/bin/bash
# Test script for SDKs in WSL
# Run this in WSL: bash test_sdks_wsl.sh

set -e

echo "=========================================="
echo "  InsureTech SDK Testing Script (WSL)"
echo "=========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if we're in WSL
if ! grep -qi microsoft /proc/version; then
    print_error "This script should be run in WSL (Windows Subsystem for Linux)"
    exit 1
fi

print_success "Running in WSL"
echo ""

# ==========================================
# Test Go SDK
# ==========================================

print_step "Testing Go SDK..."
echo ""

if [ ! -d "sdks/insuretech-go-sdk" ]; then
    print_error "Go SDK directory not found: sdks/insuretech-go-sdk"
    exit 1
fi

cd sdks/insuretech-go-sdk

# Check Go installation
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Install with: sudo apt install golang-go"
    exit 1
fi

GO_VERSION=$(go version)
print_success "Go installed: $GO_VERSION"

# Check go.mod
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found"
    exit 1
fi

print_success "go.mod found"
echo ""
cat go.mod
echo ""

# Check go.work
if [ -f "go.work" ]; then
    print_warning "go.work file found - this might cause issues"
    echo ""
    cat go.work
    echo ""
fi

# List directory structure
print_step "Go SDK directory structure:"
find . -type f -name "*.go" | head -20
echo ""

# Try to list modules
print_step "Checking Go modules..."
GOWORK=off go list ./... 2>&1 || print_warning "go list failed"
echo ""

# Try to build
print_step "Attempting to build Go SDK..."
if GOWORK=off go build ./... 2>&1; then
    print_success "Go SDK build successful!"
else
    print_error "Go SDK build failed!"
    echo ""
    print_step "Detailed error:"
    GOWORK=off go build -v ./... 2>&1 || true
    echo ""
    print_step "Checking for common issues..."
    
    # Check if pkg directory exists
    if [ ! -d "pkg" ]; then
        print_error "pkg/ directory not found"
    else
        print_success "pkg/ directory exists"
        ls -la pkg/
    fi
    
    # Check module path
    MODULE_PATH=$(grep "^module " go.mod | awk '{print $2}')
    print_step "Module path: $MODULE_PATH"
    
    # Check if files match module path
    print_step "Checking import paths in Go files..."
    grep -r "^package " pkg/ | head -10
fi

cd ../..
echo ""

# ==========================================
# Test TypeScript SDK
# ==========================================

print_step "Testing TypeScript SDK..."
echo ""

if [ ! -d "sdks/insuretech-typescript-sdk" ]; then
    print_error "TypeScript SDK directory not found: sdks/insuretech-typescript-sdk"
    exit 1
fi

cd sdks/insuretech-typescript-sdk

# Check Node.js installation
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Install with: sudo apt install nodejs npm"
    exit 1
fi

NODE_VERSION=$(node --version)
NPM_VERSION=$(npm --version)
print_success "Node.js installed: $NODE_VERSION"
print_success "npm installed: $NPM_VERSION"

# Check package.json
if [ ! -f "package.json" ]; then
    print_error "package.json not found"
    exit 1
fi

print_success "package.json found"
echo ""
cat package.json | head -30
echo ""

# List directory structure
print_step "TypeScript SDK directory structure:"
find src -type f -name "*.ts" | head -20
echo ""

# Check for problematic files
print_step "Checking for files with colons in names..."
find . -name "*:*" 2>/dev/null || print_success "No files with colons found"
echo ""

# Check imports in index files
print_step "Checking service imports..."
if [ -f "src/services/index.ts" ]; then
    echo "Contents of src/services/index.ts:"
    cat src/services/index.ts
    echo ""
    
    # Check if imported files exist
    print_step "Verifying imported service files exist..."
    grep "from '\./.*\.service'" src/services/index.ts | while read -r line; do
        # Extract filename from import
        filename=$(echo "$line" | sed -n "s/.*from '\.\\/\\(.*\\)'.*/\\1/p")
        if [ -f "src/services/${filename}.ts" ]; then
            print_success "Found: src/services/${filename}.ts"
        else
            print_error "Missing: src/services/${filename}.ts"
            # Check if file with colon exists
            if [ -f "src/services/${filename//:/:}.ts" ]; then
                print_warning "File exists but has colon in name: src/services/${filename//:/:}.ts"
            fi
        fi
    done
fi
echo ""

# Install dependencies
print_step "Installing dependencies..."
if npm install --legacy-peer-deps 2>&1; then
    print_success "Dependencies installed"
else
    print_error "Failed to install dependencies"
fi
echo ""

# Try to build
print_step "Attempting to build TypeScript SDK..."
if npm run build 2>&1; then
    print_success "TypeScript SDK build successful!"
else
    print_error "TypeScript SDK build failed!"
    echo ""
    print_step "Detailed error:"
    npm run build 2>&1 || true
    echo ""
    
    # Try typecheck separately
    print_step "Running typecheck..."
    npm run typecheck 2>&1 || true
fi

cd ../..
echo ""

# ==========================================
# Summary
# ==========================================

echo "=========================================="
echo "  Test Summary"
echo "=========================================="
echo ""
echo "Review the output above to identify issues."
echo ""
echo "Common issues:"
echo "1. Go SDK: Module path mismatch or go.work file conflicts"
echo "2. TypeScript SDK: Files with colons in names (e.g., 'products:search.service.ts')"
echo "3. Missing dependencies or incorrect import paths"
echo ""
echo "Next steps:"
echo "1. Fix any identified issues in the SDK generators"
echo "2. Regenerate SDKs"
echo "3. Re-run this test script"
echo "4. Update GitHub workflow once tests pass"
echo ""
