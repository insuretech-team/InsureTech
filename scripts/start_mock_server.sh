#!/usr/bin/env bash
# start_mock_server.sh — Start Prism mock server for local frontend development
#
# Usage:
#   ./scripts/start_mock_server.sh [--port=4010] [--dynamic]
#
# Prerequisites: Node.js 18+
# Prism is auto-installed via npx on first run.
# Pin a known-good Prism version because the latest release currently fails
# to start under Node 24 on this machine due to a missing transitive module.
#
# What this does:
#   Starts a Prism HTTP mock server that reads api/openapi.yaml and returns
#   example responses for every endpoint. Frontend teams can point their apps
#   at http://localhost:4010 without needing a running backend.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
OPENAPI_SPEC="$PROJECT_ROOT/api/openapi.yaml"

PORT=4010
DYNAMIC=false

for arg in "$@"; do
    case "$arg" in
        --port=*) PORT="${arg#*=}" ;;
        --dynamic) DYNAMIC=true ;;
        *) echo "Unknown argument: $arg"; exit 1 ;;
    esac
done

GREEN='\033[0;32m'; CYAN='\033[0;36m'; YELLOW='\033[1;33m'; RESET='\033[0m'

echo ""
echo -e "${CYAN}========================================${RESET}"
echo -e "${CYAN}   InsureTech Prism Mock Server${RESET}"
echo -e "${CYAN}========================================${RESET}"

if [ ! -f "$OPENAPI_SPEC" ]; then
    echo -e "  ${YELLOW}⚠ OpenAPI spec not found at: $OPENAPI_SPEC${RESET}"
    echo -e "  Run ./run_api_pipeline.sh first to generate the spec."
    exit 1
fi

echo -e "  ${GREEN}✓${RESET} OpenAPI spec: $OPENAPI_SPEC"
echo -e "  ${GREEN}✓${RESET} Mock server:  http://localhost:${PORT}"
if [ "$DYNAMIC" = true ]; then
    echo -e "  ${GREEN}✓${RESET} Mode: dynamic (randomised responses)"
else
    echo -e "  ${GREEN}✓${RESET} Mode: static (first example from spec)"
fi
echo ""
echo -e "  Frontend usage:"
echo -e "    Set your API base URL to: ${CYAN}http://localhost:${PORT}${RESET}"
echo -e "    All endpoints return example data from the OpenAPI spec."
echo -e "    No authentication required (mock server ignores Bearer tokens)."
echo ""
echo -e "  Press Ctrl+C to stop the mock server."
echo ""

PRISM_ARGS=("mock" "$OPENAPI_SPEC" "--port" "$PORT" "--cors")
if [ "$DYNAMIC" = true ]; then
    PRISM_ARGS+=("--dynamic")
fi

PRISM_VERSION="@stoplight/prism-cli@5.14.2"
npx --yes "$PRISM_VERSION" "${PRISM_ARGS[@]}"
