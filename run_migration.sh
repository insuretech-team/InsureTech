#!/usr/bin/env bash
# run_migration.sh — InsureTech Database Migration Runner
# Equivalent of run_migration.ps1 for Linux/macOS/Codespaces
#
# Usage:
#   ./run_migration.sh [--target=primary|backup|both] [--dry-run] [--prune] [--strict]

set -euo pipefail

# ── Parse args ────────────────────────────────────────────────────────────────
TARGET="primary"
DRY_RUN=false
PRUNE=false
STRICT=false

for arg in "$@"; do
    case "$arg" in
        --target=*)  TARGET="${arg#*=}" ;;
        --dry-run)   DRY_RUN=true ;;
        --prune)     PRUNE=true ;;
        --strict)    STRICT=true ;;
        *) echo "Unknown argument: $arg" >&2; echo "Usage: $0 [--target=primary|backup|both] [--dry-run] [--prune] [--strict]"; exit 1 ;;
    esac
done

# Validate target
case "$TARGET" in
    primary|backup|both) ;;
    *) echo "ERROR: --target must be primary, backup, or both" >&2; exit 1 ;;
esac

# ── Colours ───────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; CYAN='\033[0;36m'; GRAY='\033[0;37m'; RESET='\033[0m'
ok()   { echo -e "  ${GREEN}OK${RESET}: $*"; }
warn() { echo -e "  ${YELLOW}WARNING${RESET}: $*"; }
err()  { echo -e "  ${RED}ERROR${RESET}: $*" >&2; exit 1; }
info() { echo -e "  ${GRAY}$*${RESET}"; }

# ── Locate project root (where go.mod lives) ──────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
    # Walk up in case script is run from a subdirectory
    dir="$PROJECT_ROOT"
    while [ "$dir" != "/" ]; do
        if [ -f "$dir/go.mod" ]; then PROJECT_ROOT="$dir"; break; fi
        dir="$(dirname "$dir")"
    done
fi
if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
    err "go.mod not found — run this script from the project root"
fi

echo -e "${CYAN}==========================================${RESET}"
echo -e "${CYAN}InsureTech Database Migration Runner${RESET}"
echo -e "${CYAN}==========================================${RESET}"

# ── Step 0: Bootstrap prerequisites ──────────────────────────────────────────
echo -e "\n${YELLOW}[0/7] Checking prerequisites...${RESET}"
BOOTSTRAP="$PROJECT_ROOT/scripts/bootstrap.sh"
if [ -f "$BOOTSTRAP" ]; then
    # shellcheck source=scripts/bootstrap.sh
    bash "$BOOTSTRAP" --go-only
else
    command -v go &>/dev/null || err "'go' not found and bootstrap.sh is missing. Install Go 1.25 from https://go.dev/dl/"
    ok "go found"
fi

# ── Step 1: Project root ──────────────────────────────────────────────────────
echo -e "\n${YELLOW}[1/7] Finding project root...${RESET}"
ok "Project root: $PROJECT_ROOT"

# ── Step 2: Load .env ─────────────────────────────────────────────────────────
echo -e "\\n${YELLOW}[2/7] Loading environment variables...${RESET}"
ENV_FILE="$PROJECT_ROOT/.env"
if [ -f "$ENV_FILE" ]; then
    ok "Found .env"
    # .env contains quoted values and PEM blobs, so preserve it as shell syntax
    # and only normalize trailing CR characters when edited on Windows.
    set -o allexport
    # shellcheck disable=SC1090
    source <(sed 's/\r$//' "$ENV_FILE")
    set +o allexport
else
    warn ".env not found — using existing environment variables"
fi

# Verify required DB env vars (using explicit check for macOS bash 3.2 compatibility)
MISSING=()
for var in PGHOST PGDATABASE PGUSER PGPASSWORD; do
    val=$(printenv "$var" || true)
    [ -z "$val" ] && MISSING+=("$var")
done
if [ "${#MISSING[@]}" -gt 0 ]; then
    err "Missing required environment variables: ${MISSING[*]}"
fi
ok "All environment variables loaded"

# ── Step 3: Verify config ─────────────────────────────────────────────────────
echo -e "\n${YELLOW}[3/7] Verifying configuration...${RESET}"
CONFIG_PATH="$PROJECT_ROOT/backend/inscore/configs/database.yaml"
[ -f "$CONFIG_PATH" ] || err "database.yaml not found at $CONFIG_PATH"
ok "Config file found"

# ── Step 4: SSL certificates ──────────────────────────────────────────────────
echo -e "\n${YELLOW}[4/7] Checking SSL certificates...${RESET}"
CERTS_DIR="$PROJECT_ROOT/backend/inscore/db/certs"
if [ -d "$CERTS_DIR" ]; then
    ok "Certs directory exists"
else
    info "No certs directory (OK if not using SSL)"
fi

# ── Step 5: Regenerate proto code ─────────────────────────────────────────────
echo -e "\n${YELLOW}[5/7] Regenerating proto code with GORM tags...${RESET}"
GENERATE_SCRIPT="$PROJECT_ROOT/scripts/generate.sh"
if [ -f "$GENERATE_SCRIPT" ]; then
    bash "$GENERATE_SCRIPT" && ok "Proto code regenerated" || warn "Proto generation had issues (continuing...)"
else
    info "SKIP: generate.sh not found"
fi

# ── Step 6: Prepare dbx ─────────────────────────────────────────────────
echo -e "\n${YELLOW}[6/7] Preparing dbx...${RESET}"
DBMANAGER_DIR="$PROJECT_ROOT/backend/inscore/cmd/dbx"
[ -d "$DBMANAGER_DIR" ] || err "dbx directory not found at $DBMANAGER_DIR"
ok "Will use go run directly from project root"

# ── Step 7: Run migration ─────────────────────────────────────────────────────
echo -e "\n${YELLOW}[7/7] Running migration on $TARGET...${RESET}"
[ "$PRUNE"  = true ] && echo -e "  ${YELLOW}Mode: PRUNE (will remove zombie columns)${RESET}"
[ "$STRICT" = true ] && echo -e "  ${YELLOW}Mode: STRICT (will fail on schema drift)${RESET}"
echo -e "${CYAN}==========================================${RESET}"
echo ""

CMD_ARGS=("migrate" "--target" "$TARGET")
[ "$PRUNE"  = true ] && CMD_ARGS+=("--prune")
[ "$STRICT" = true ] && CMD_ARGS+=("--strict")

if [ "$DRY_RUN" = true ]; then
    echo -e "${YELLOW}DRY RUN — Would execute:${RESET}"
    echo "  cd $PROJECT_ROOT"
    echo "  go run ./backend/inscore/cmd/dbx ${CMD_ARGS[*]}"
    exit 0
fi

echo -e "${CYAN}Executing: go run ./backend/inscore/cmd/dbx ${CMD_ARGS[*]}${RESET}"
echo ""

cd "$PROJECT_ROOT" || err "Failed to change to project root directory"
go run ./backend/inscore/cmd/dbx "${CMD_ARGS[@]}"

echo ""
echo -e "${CYAN}==========================================${RESET}"
echo -e "${GREEN}SUCCESS: Migration completed!${RESET}"
echo -e "${CYAN}==========================================${RESET}"
