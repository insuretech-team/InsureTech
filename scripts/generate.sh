#!/bin/bash
# Unified Proto Generation Script for InsureTech (Linux/Mac)
# Verifies tools, loads BUF_TOKEN, then delegates to inject_gorm_tags.go

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
VERBOSE=false
SKIP_REGISTRY=false

while [[ "$#" -gt 0 ]]; do
    case $1 in
        --skip-registry) SKIP_REGISTRY=true ;;
        --verbose) VERBOSE=true ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

echo "==========================================="
echo "InsureTech Proto Generation"
echo "==========================================="
echo "OS: $(uname -s), Arch: $(uname -m)"
echo ""

# Step 1: Verify tools
echo "[1/2] Verifying tools..."
check_tool() {
    if ! command -v "$1" &> /dev/null; then
        echo "  Installing $1..."
        eval "$2"
    else
        echo "  OK: $1"
    fi
}

check_tool "buf" "go install github.com/bufbuild/buf/cmd/buf@latest"
check_tool "protoc-gen-go" "go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"

# Load BUF_TOKEN from .env if it exists
if [ -f "$REPO_ROOT/.env" ]; then
    BUF_TOKEN=$(grep '^BUF_TOKEN=' "$REPO_ROOT/.env" 2>/dev/null | sed 's/^BUF_TOKEN=//')
    if [ -n "$BUF_TOKEN" ]; then
        export BUF_TOKEN
        echo "  Loaded BUF_TOKEN from .env"
    fi
fi

# Step 2: Run unified Go script (buf generate + GORM injection + registry)
echo ""
echo "[2/2] Running generation pipeline..."

cd "$REPO_ROOT"
GO_ARGS="--generate"
if [ "$SKIP_REGISTRY" = false ]; then
    GO_ARGS="$GO_ARGS --registry"
fi
if [ "$VERBOSE" = true ]; then
    GO_ARGS="$GO_ARGS --verbose"
fi
go run ./scripts/inject_gorm_tags.go $GO_ARGS

echo ""
echo "==========================================="
echo "Proto generation complete!"
echo "==========================================="
