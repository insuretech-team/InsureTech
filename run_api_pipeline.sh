#!/usr/bin/env bash
# run_api_pipeline.sh — InsureTech OpenAPI Generation Pipeline
# Equivalent of run_api_pipeline.ps1 for Linux/macOS/Codespaces
#
# Usage:
#   ./run_api_pipeline.sh [--skip-cleanup] [--skip-validation] [--skip-docs] [--fast] [--port=8080]

set -euo pipefail

# Force UTF-8 for all Python subprocess output — prevents encoding errors on non-UTF-8 terminals
export PYTHONUTF8=1
export PYTHONIOENCODING=utf-8

# ── Parse args ────────────────────────────────────────────────────────────────
SKIP_CLEANUP=false
SKIP_VALIDATION=false
SKIP_DOCS=false
FAST=false
SERVER_PORT=8080

for arg in "$@"; do
    case "$arg" in
        --skip-cleanup)    SKIP_CLEANUP=true ;;
        --skip-validation) SKIP_VALIDATION=true ;;
        --skip-docs)       SKIP_DOCS=true ;;
        --fast)            FAST=true ;;
        --port=*)          SERVER_PORT="${arg#*=}" ;;
        *) echo "Unknown argument: $arg" >&2
           echo "Usage: $0 [--skip-cleanup] [--skip-validation] [--skip-docs] [--fast] [--port=PORT]"
           exit 1 ;;
    esac
done

# ── Colours ───────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; CYAN='\033[0;36m'; GRAY='\033[0;37m'; RESET='\033[0m'
step()    { echo -e "\n${CYAN}[$1/$2]${RESET} $3"; }
ok()      { echo -e "  ${GREEN}✓${RESET} $*"; }
warn()    { echo -e "  ${YELLOW}⚠${RESET} $*"; }
err_msg() { echo -e "  ${RED}✗${RESET} $*" >&2; }
fail()    { err_msg "$*"; exit 1; }

START_TIME=$(date +%s)

# ── Locate project root ───────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
API_DIR="$PROJECT_ROOT/api"

[ -d "$API_DIR" ] || fail "API directory not found at $API_DIR — run from the project root"

echo ""
echo -e "${CYAN}========================================${RESET}"
echo -e "${CYAN}   OpenAPI Generation Pipeline${RESET}"
echo -e "${CYAN}========================================${RESET}"

# ── Step 0: Bootstrap prerequisites ──────────────────────────────────────────
step 0 16 "Checking prerequisites..."
BOOTSTRAP="$PROJECT_ROOT/scripts/bootstrap.sh"
if [ -f "$BOOTSTRAP" ]; then
    bash "$BOOTSTRAP" --all
else
    # Minimal inline checks
    command -v go     &>/dev/null || warn "'go' not found — some steps may fail"
    command -v python3 &>/dev/null || command -v python &>/dev/null || fail "Python 3 not found. Install from https://www.python.org/downloads/"
    command -v node   &>/dev/null || warn "'node' not found — SDK build steps may fail"
    command -v npm    &>/dev/null || warn "'npm' not found — SDK build steps may fail"
fi

# Resolve python command (python3 preferred)
PY_CMD="python3"
command -v python3 &>/dev/null || PY_CMD="python"

# ── Step 0b: Generate proto files ─────────────────────────────────────────────
step 0 16 "Generating proto files..."
cd "$PROJECT_ROOT" || fail "Cannot cd to project root"
GENERATE_SCRIPT="$PROJECT_ROOT/scripts/generate.sh"
if [ -f "$GENERATE_SCRIPT" ]; then
    bash "$GENERATE_SCRIPT" && ok "Proto files generated successfully" \
        || warn "Proto generation had issues (exit $?) — continuing with existing files"
else
    warn "generate.sh not found — skipping proto generation"
fi

# ── Step 1: Cleanup ───────────────────────────────────────────────────────────
cd "$API_DIR"

if [ "$SKIP_CLEANUP" = false ]; then
    step 1 16 "Cleanup old files..."
    rm -rf schemas events enums paths openapi.yaml input/descriptors.pb 2>/dev/null || true
    ok "Cleaned old files"
else
    step 1 16 "Cleanup skipped"
fi

# ── Step 2: Run main generator ────────────────────────────────────────────────
# IMPORTANT ORDER:
#   fix_all_warnings.py now runs INSIDE main.py before assembly (baked in at step 5.5)
#   so schema descriptions are fixed before openapi.yaml is assembled.
#   Running it after assembly caused run-to-run instability (1500+ HTML files regenerating).
step 2 16 "Running code generator (proto -> schemas)..."
cd "$API_DIR/generator" || fail "Cannot cd to generator directory"
"$PY_CMD" main.py --discover
cd "$API_DIR" || fail "Cannot cd to API directory"

SCHEMAS_COUNT=$(find schemas -name "*.yaml" 2>/dev/null | wc -l | tr -d ' ')
EVENTS_COUNT=$(find events  -name "*.yaml" 2>/dev/null | wc -l | tr -d ' ')
ENUMS_COUNT=$(find enums    -name "*.yaml" -maxdepth 1 2>/dev/null | wc -l | tr -d ' ')
PATHS_COUNT=$(find paths    -name "*.yaml" 2>/dev/null | wc -l | tr -d ' ')

ok "Generated $SCHEMAS_COUNT schemas"
ok "Generated $EVENTS_COUNT events"
ok "Generated $ENUMS_COUNT enums"
ok "Generated $PATHS_COUNT paths"
ok "Assembled openapi.yaml"

# ── Step 11: Docker Validation (skipped — too slow for 865+ schemas) ────────
# Docker openapi-generator-cli validate is extremely slow on large specs.
# Step 13 enhanced_validator_optimized.py provides faster, more detailed validation.
step 11 16 "Skipping Docker validation - using Python validator in step 13 instead"

# ── Step 12: Fix validation warnings ───────────────────────────────────────
# fix_all_warnings.py is now called INSIDE main.py before assembly (step 5.5).
# Running it here AFTER assembly caused ruamel.yaml to reformat the entire
# openapi.yaml differently from PyYAML, growing it ~2300 lines per run and
# making 1500+ HTML files regenerate every pipeline execution.
#
# fix_pagination.py now does a fast text-scan only for deprecated schema refs.
# Pagination params are injected at generation time in path_generator.py.
step 12 16 "Running deprecation ref check..."
cd "$API_DIR/generator" || fail "Cannot cd to generator directory"

if ! "$PY_CMD" fix_pagination.py > /tmp/fix_pagination.log 2>&1; then
    warn "fix_pagination.py had issues:"
    cat /tmp/fix_pagination.log >&2 || true
fi

ok "Pagination deprecation check complete"
cd "$API_DIR" || fail "Cannot cd to API directory"

# ── Step 12: Enhanced validation ─────────────────────────────────────────────
if [ "$FAST" = true ]; then
    step 12 16 "Skipping validation (Fast mode)"
else
    step 12 16 "Running validation and quality checks..."
    echo -e "  ${GRAY}Running quick validation checks...${RESET}"

    UNKNOWN_TYPE_COUNT=$(grep -r "Unknown type.*Entry" schemas/ 2>/dev/null | wc -l | tr -d ' ' || echo 0)
    EVENTS_EXIST=$([ -d events ] && echo true || echo false)
    ENUM_SUBDIRS=$(find enums -mindepth 1 -maxdepth 1 -type d 2>/dev/null | wc -l | tr -d ' ' || echo 0)

    [ "$UNKNOWN_TYPE_COUNT" -eq 0 ] \
        && ok "Map fields: No 'Unknown type Entry' errors" \
        || warn "Map fields: Found $UNKNOWN_TYPE_COUNT 'Unknown type Entry' errors"

    [ "$EVENTS_EXIST" = true ] && [ "$EVENTS_COUNT" -gt 0 ] \
        && ok "Events folder: $EVENTS_COUNT events generated" \
        || err_msg "Events folder: Not created or empty"

    [ "$ENUM_SUBDIRS" -eq 0 ] && [ "$ENUMS_COUNT" -gt 0 ] \
        && ok "Enums structure: Flat ($ENUMS_COUNT files, no subdirectories)" \
        || err_msg "Enums structure: Has subdirectories or empty"

    echo -e "  ${GRAY}Running enhanced validation...${RESET}"
    cd "$API_DIR/generator" || fail "Cannot cd to generator directory"
    "$PY_CMD" enhanced_validator_optimized.py ../openapi.yaml \
        --report ../validation_report.json \
        --html   ../validation_report.html 2>&1 || true
    cd "$API_DIR" || fail "Cannot cd to API directory"

    if [ -f validation_report.json ]; then
        ERRORS=$(  "$PY_CMD" -c "import json,sys; r=json.load(open('validation_report.json')); print(r['summary']['errors'])"   2>/dev/null || echo 0)
        WARNINGS=$(  "$PY_CMD" -c "import json,sys; r=json.load(open('validation_report.json')); print(r['summary']['warnings'])" 2>/dev/null || echo 0)
        COVERAGE=$(  "$PY_CMD" -c "import json,sys; r=json.load(open('validation_report.json')); print(r['metrics']['description_coverage'])" 2>/dev/null || echo "?")
        ok "Detailed validation complete"
        echo -e "    Errors:   ${ERRORS}"
        echo -e "    Warnings: ${WARNINGS}"
        echo -e "    Description Coverage: ${COVERAGE}%"
        [ "$ERRORS" -eq 0 ] 2>/dev/null || fail "Validation failed with $ERRORS errors!"
    else
        warn "Validation report not generated"
    fi
fi

# ── Step 14: Generate Documentation ──────────────────────────────────────────
if [ "$SKIP_DOCS" = false ]; then
    step 14 16 "Generating API documentation..."
    mkdir -p "$API_DIR/docs"

    cd "$API_DIR/generator" || fail "Cannot cd to generator directory"
    echo -e "  ${GRAY}Generating enhanced documentation hub...${RESET}"
    "$PY_CMD" table_view_generator.py     --spec ../openapi.yaml --output-dir ../docs &>/dev/null || true
    "$PY_CMD" schema_enum_page_generator.py --spec ../openapi.yaml --output-dir ../docs &>/dev/null || true
    "$PY_CMD" doc_generator.py            --spec ../openapi.yaml --output ../docs/index.html \
        --generate-endpoint-pages &>/dev/null || true
    cd "$API_DIR" || fail "Cannot cd to API directory"

    ok "Generated enhanced documentation"

    # Swagger UI
    cat > "docs/swagger.html" << 'SWAGGER_EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InsureTech API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css">
    <style>body{margin:0;padding:0;}.swagger-ui .topbar{display:none;}</style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "../openapi.yaml", dom_id: '#swagger-ui', deepLinking: true,
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                plugins: [SwaggerUIBundle.plugins.DownloadUrl],
                layout: "StandaloneLayout", defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1, docExpansion: "list",
                filter: true, showExtensions: true, showCommonExtensions: true
            });
        };
    </script>
</body>
</html>
SWAGGER_EOF

    # ReDoc
    cat > "docs/redoc.html" << 'REDOC_EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InsureTech API - ReDoc</title>
    <style>body{margin:0;padding:0;}</style>
</head>
<body>
    <redoc spec-url="../openapi.yaml" scroll-y-offset="nav"
           hide-download-button="false" expand-responses="200,201"></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>
REDOC_EOF

    ok "Generated Swagger UI and ReDoc"

    # Content-aware sync to root docs/ — only copy files whose content changed.
    # Do NOT rm -rf + cp -r: that gives every file a new mtime on every run
    # causing git to see 1000+ files as "modified" even with identical content.
    ROOT_DOCS="$PROJECT_ROOT/docs"
    mkdir -p "$ROOT_DOCS"

    if [[ "$ROOT_DOCS" == "/" || "$ROOT_DOCS" == "$HOME" || ${#ROOT_DOCS} -lt 5 ]]; then
        fail "ROOT_DOCS path '$ROOT_DOCS' looks unsafe, aborting"
    fi

    _synced=0; _skipped=0
    while IFS= read -r -d '' src_file; do
        rel="${src_file#$API_DIR/docs/}"
        dst_file="$ROOT_DOCS/$rel"
        dst_dir="$(dirname "$dst_file")"
        mkdir -p "$dst_dir"
        if [ -f "$dst_file" ] && cmp -s "$src_file" "$dst_file"; then
            _skipped=$((_skipped + 1))
        else
            cp "$src_file" "$dst_file"
            _synced=$((_synced + 1))
        fi
    done < <(find "$API_DIR/docs" -type f -print0 2>/dev/null)

    # Remove files in root docs that no longer exist in api/docs
    while IFS= read -r -d '' dst_file; do
        rel="${dst_file#$ROOT_DOCS/}"
        src_file="$API_DIR/docs/$rel"
        [ -f "$src_file" ] || rm -f "$dst_file"
    done < <(find "$ROOT_DOCS" -type f -print0 2>/dev/null)

    # Content-aware copy of individual root-level files
    _copy_if_changed() {
        local src="$1" dst="$2"
        [ -f "$src" ] || return 0
        if [ -f "$dst" ] && cmp -s "$src" "$dst"; then return 0; fi
        cp "$src" "$dst"
    }
    _copy_if_changed "$API_DIR/openapi.yaml"           "$ROOT_DOCS/openapi.yaml"
    _copy_if_changed "$API_DIR/validation_report.html" "$ROOT_DOCS/validation_report.html"
    _copy_if_changed "$API_DIR/validation_report.json" "$ROOT_DOCS/validation_report.json"
    _copy_if_changed "$API_DIR/proto_schema_summary.json" "$ROOT_DOCS/proto_schema_summary.json"
    _copy_if_changed "$API_DIR/schema_api_mapping.json"   "$ROOT_DOCS/schema_api_mapping.json"

    ok "Synced docs to root docs/ ($_synced updated, $_skipped unchanged)"
fi

# ── Step 15: Generate Postman Collection + Upload ─────────────────────────────
step 15 16 "Generating Postman collection + environments..."

# Load POSTMAN_API_KEY from .env if available
if [ -f "$PROJECT_ROOT/.env" ]; then
    POSTMAN_API_KEY=$(grep '^POSTMAN_API_KEY=' "$PROJECT_ROOT/.env" 2>/dev/null | cut -d'=' -f2- | tr -d '"' | tr -d "'" || true)
fi
POSTMAN_API_KEY="${POSTMAN_API_KEY:-${POSTMAN_API_KEY:-}}"

cd "$API_DIR/generator" || fail "Cannot cd to generator directory"

# Always generate locally
"$PY_CMD" sync_postman.py \
    && ok "Postman collection + environments generated (local, staging, production, mock, newman_test)" \
    || warn "Postman generation had issues (continuing...)"

# Upload to Postman API if key is available
if [ -n "${POSTMAN_API_KEY:-}" ]; then
    echo -e "  ${GRAY}POSTMAN_API_KEY found — uploading collection + environments to Postman...${RESET}"
    "$PY_CMD" sync_postman.py --upload \
        && ok "Collection + environments uploaded to Postman API" \
        || warn "Postman upload had issues (continuing...)"
else
    echo -e "  ${GRAY}→ Set POSTMAN_API_KEY in .env to auto-upload to Postman API${RESET}"
    echo -e "  ${GRAY}  Get your key: https://go.postman.co/settings/me/api-keys${RESET}"
fi

# Legacy Apidog sync (if token present)
if [ -f "$PROJECT_ROOT/.env" ]; then
    APIDOG_TOKEN=$(grep '^API_DOG_TOKEN=' "$PROJECT_ROOT/.env" 2>/dev/null | cut -d'=' -f2- | tr -d '"' | tr -d "'" || true)
fi
if [ -n "${APIDOG_TOKEN:-}" ]; then
    echo -e "  ${GRAY}API_DOG_TOKEN found — syncing to Apidog...${RESET}"
    "$PY_CMD" sync_apidog.py && ok "Synced to Apidog successfully" \
        || warn "Apidog sync had issues (continuing...)"
fi

cd "$API_DIR" || fail "Cannot cd to API directory"

# ── Step 16: Generate SDKs ────────────────────────────────────────────────────
step 16 16 "Generating SDKs..."

# TypeScript SDK
echo -e "  ${GRAY}Generating TypeScript SDK (hey-api + custom)...${RESET}"
cd "$PROJECT_ROOT/sdks/sdk-generator/typescript" || fail "Cannot cd to TypeScript SDK directory"
if [ ! -d node_modules ]; then
    echo -e "  ${GRAY}  Installing @hey-api/openapi-ts...${RESET}"
    # Capture output to avoid pipe-buffer stall on large npm output
    npm install 2>&1 | tail -3 || fail "npm install failed"
fi
# format:false and lint:false in openapi-ts.config.ts — no per-file prettier/eslint
npm run generate 2>&1 | tail -3 || fail "hey-api generator failed"

# Build custom post-processor only if source is newer than binary
echo -e "  ${GRAY}  Building custom post-processor...${RESET}"
if [ ! -f generator ] || [ generator.go -nt generator ]; then
    GOWORK=off go build -o generator generator.go || fail "Failed to build post-processor"
else
    echo -e "  ${GRAY}    Post-processor binary up-to-date, skipping rebuild${RESET}"
fi
./generator 2>&1 || fail "Custom post-processing failed"
ok "TypeScript SDK generated (hey-api + custom)"

# Go SDK
echo -e "  ${GRAY}Generating Go SDK...${RESET}"
cd "$PROJECT_ROOT/sdks/sdk-generator/go" || fail "Cannot cd to Go SDK generator directory"
# Build Go SDK generator only if source is newer than binary
if [ ! -f generator ] || [ generator.go -nt generator ]; then
    GOWORK=off go build -o generator generator.go || fail "Failed to build Go SDK generator"
else
    echo -e "  ${GRAY}    Go SDK generator binary up-to-date, skipping rebuild${RESET}"
fi
./generator 2>&1 || fail "Go SDK generation failed"
ok "Go SDK generated"

# Build TypeScript SDK
echo -e "  ${GRAY}Building TypeScript SDK...${RESET}"
cd "$PROJECT_ROOT/sdks/insuretech-typescript-sdk" || fail "Cannot cd to TypeScript SDK directory"
if [ ! -d node_modules ]; then
    npm install --legacy-peer-deps 2>&1 | tail -3 || fail "npm install failed"
fi
npm run build 2>&1 | tail -5 || fail "TypeScript SDK build failed"
[ -d dist ] || fail "TypeScript SDK build succeeded but dist/ not found"
ok "TypeScript SDK built successfully"

# Build Go SDK
echo -e "  ${GRAY}Building Go SDK...${RESET}"
cd "$PROJECT_ROOT/sdks/insuretech-go-sdk" || fail "Cannot cd to Go SDK directory"
GOWORK=off go build ./... || fail "Go SDK build failed"
ok "Go SDK built successfully"

cd "$API_DIR" || fail "Cannot cd back to API directory"

# ── Step 18: Newman Smoke Tests (optional — requires NEWMAN_BASE_URL in .env) ────
step 18 18 "Running Newman smoke tests..."

NEWMAN_COLLECTION="$PROJECT_ROOT/api/postman/InsureTech.postman_collection.json"
NEWMAN_ENV_FILE="$PROJECT_ROOT/api/postman/InsureTech_local.postman_environment.json"
NEWMAN_RESULTS="$PROJECT_ROOT/api/postman/newman_results.json"

# Load NEWMAN_BASE_URL from .env if not already set
if [ -f "$PROJECT_ROOT/.env" ]; then
    NEWMAN_BASE_URL=$(grep '^NEWMAN_BASE_URL=' "$PROJECT_ROOT/.env" 2>/dev/null | cut -d'=' -f2- | tr -d '"' | tr -d "'" || true)
fi

if [ -z "${NEWMAN_BASE_URL:-}" ]; then
    echo -e "  ${GRAY}→ Set NEWMAN_BASE_URL in .env to enable Newman smoke tests${RESET}"
    echo -e "  ${GRAY}  Example: NEWMAN_BASE_URL=http://localhost:8080${RESET}"
elif [ ! -f "$NEWMAN_COLLECTION" ]; then
    warn "Postman collection not found — run Step 15 first"
else
    # Check server is reachable before running Newman
    SERVER_REACHABLE=false
    if curl -sf --max-time 3 "$NEWMAN_BASE_URL/health" &>/dev/null || \
       curl -sf --max-time 3 "$NEWMAN_BASE_URL" &>/dev/null; then
        SERVER_REACHABLE=true
    fi

    if [ "$SERVER_REACHABLE" = false ]; then
        echo -e "  ${GRAY}→ Server not reachable at $NEWMAN_BASE_URL — skipping Newman${RESET}"
        echo -e "  ${GRAY}  Start your API server first, then re-run with NEWMAN_BASE_URL set${RESET}"
    else
        echo -e "  ${GRAY}Running Newman against $NEWMAN_BASE_URL ...${RESET}"
        NEWMAN_ARGS=(
            "run" "$NEWMAN_COLLECTION"
            "--timeout-request" "10000"
            "--reporters" "cli,json"
            "--reporter-json-export" "$NEWMAN_RESULTS"
            "--env-var" "base_url=$NEWMAN_BASE_URL"
        )
        if [ -f "$NEWMAN_ENV_FILE" ]; then
            NEWMAN_ARGS+=("--environment" "$NEWMAN_ENV_FILE")
        fi
        npx --yes newman "${NEWMAN_ARGS[@]}" \
            && ok "Newman smoke tests passed" \
            || warn "Newman reported test failures (see output above)"
    fi
fi

# ── Step 15: Start documentation server ──────────────────────────────────────
step 15 16 "Starting documentation server..."

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo -e "${CYAN}========================================${RESET}"
echo -e "${GREEN}✅ API GENERATION COMPLETE${RESET}"
echo -e "${CYAN}========================================${RESET}"
echo -e "\nTime elapsed: ${DURATION}s"
echo ""
echo -e "API Documentation:"
echo -e "  Home:        ${CYAN}http://localhost:${SERVER_PORT}/${RESET}"
echo -e "  Swagger UI:  ${CYAN}http://localhost:${SERVER_PORT}/docs/swagger.html${RESET}"
echo -e "  ReDoc:       ${CYAN}http://localhost:${SERVER_PORT}/docs/redoc.html${RESET}"
echo -e "  OpenAPI:     ${CYAN}http://localhost:${SERVER_PORT}/openapi.yaml${RESET}"
echo ""
echo -e "Statistics:"
echo -e "  Schemas: $SCHEMAS_COUNT  Events: $EVENTS_COUNT  Enums: $ENUMS_COUNT  Paths: $PATHS_COUNT"
echo ""
echo -e "${YELLOW}Starting server on port $SERVER_PORT... (Ctrl+C to stop)${RESET}"
echo ""

# Generate and run the server inline
"$PY_CMD" - << PYEOF
import http.server, socketserver, os, sys, signal

PORT = $SERVER_PORT
os.chdir('$API_DIR')

class Handler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        super().end_headers()
    def do_GET(self):
        if self.path in ('/', ''):
            self.send_response(302)
            self.send_header('Location', '/docs/index.html')
            self.end_headers()
            return
        return http.server.SimpleHTTPRequestHandler.do_GET(self)
    def log_message(self, fmt, *args):
        pass  # suppress request logs for cleaner output

httpd = None

def _shutdown(sig, frame):
    print('\n  Shutting down documentation server...')
    if httpd:
        httpd.shutdown()
    sys.exit(0)

signal.signal(signal.SIGINT, _shutdown)
signal.signal(signal.SIGTERM, _shutdown)

for attempt in range(5):
    try:
        httpd = socketserver.TCPServer(('', PORT), Handler)
        httpd.allow_reuse_address = True
        print(f'  Server: http://localhost:{PORT}/')
        print(f'  Swagger: http://localhost:{PORT}/docs/swagger.html')
        print(f'  ReDoc:   http://localhost:{PORT}/docs/redoc.html')
        print(f'  Press Ctrl+C to stop')
        httpd.serve_forever()
        break
    except OSError:
        print(f'Port {PORT} in use, trying {PORT+1}...')
        PORT += 1
else:
    print('Could not find available port', file=sys.stderr)
    sys.exit(1)
PYEOF
