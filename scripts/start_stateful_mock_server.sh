#!/usr/bin/env bash
# start_stateful_mock_server.sh — Start the stateful Python mock server
#
# Usage:
#   ./scripts/start_stateful_mock_server.sh [--port=4010]

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SERVER_SCRIPT="$PROJECT_ROOT/scripts/stateful_mock_server.py"
PORT=4010

for arg in "$@"; do
    case "$arg" in
        --port=*) PORT="${arg#*=}" ;;
        *) echo "Unknown argument: $arg"; exit 1 ;;
    esac
done

echo ""
echo "============================================================"
echo "  InsureTech Stateful Mock Server"
echo "============================================================"
echo "  Base URL:       http://localhost:${PORT}"
echo "  Reset state:    POST http://localhost:${PORT}/_mock/reset"
echo "  Inspect state:  GET  http://localhost:${PORT}/_mock/state"
echo "============================================================"
echo ""

python "$SERVER_SCRIPT" --port "$PORT"
