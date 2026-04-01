#!/bin/bash
# render-nginx-conf.sh — Renders nginx config templates from .env / .env.prod
#
# Usage:
#   ./render-nginx-conf.sh [--env <path-to-env-file>] [--out <output-dir>]
#
# Reads env file, substitutes __PLACEHOLDER__ tokens in all conf.d/*.conf files,
# writes rendered output to dist/conf.d/ (default) or --out directory.
# Rendered files are ready to be copied to /etc/nginx/conf.d/.
#
# Called by:
#   - setup.sh (local install)
#   - quickerdeploy.sh step 3d (remote deploy via SSH heredoc)
#   - CI/CD pipeline (pre-deployment-check)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NGINX_DIR="$(dirname "$SCRIPT_DIR")"           # backend/infra/nginx
INFRA_DIR="$(dirname "$NGINX_DIR")"            # backend/infra
BACKEND_DIR="$(dirname "$INFRA_DIR")"          # backend
PROJECT_ROOT="$(dirname "$BACKEND_DIR")"       # project root (InsureTech/)

# ── Defaults ──────────────────────────────────────────────────────────────────
ENV_FILE=""
OUT_DIR="$NGINX_DIR/dist/conf.d"

# ── Arg parsing ───────────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
    case "$1" in
        --env)  ENV_FILE="$2";  shift 2 ;;
        --out)  OUT_DIR="$2";   shift 2 ;;
        *)      echo "Unknown arg: $1"; exit 1 ;;
    esac
done

# ── Resolve env file: prefer explicit --env, then .env.prod, then .env ────────
if [ -z "$ENV_FILE" ]; then
    for candidate in \
        "$PROJECT_ROOT/.env.prod" \
        "$PROJECT_ROOT/.env" \
        "$NGINX_DIR/.env.prod" \
        "$NGINX_DIR/.env"; do
        if [ -f "$candidate" ]; then
            ENV_FILE="$candidate"
            break
        fi
    done
fi

if [ -z "$ENV_FILE" ]; then
    echo "⚠  No .env file found — using built-in defaults only."
else
    echo "📋 Loading env from: $ENV_FILE"
    # Export all vars from the env file (skip comments and blank lines)
    set -o allexport
    # shellcheck source=/dev/null
    source <(grep -E '^[A-Z_][A-Z0-9_]*=' "$ENV_FILE" | sed 's/\r//')
    set +o allexport
fi

# ── Built-in defaults (fallback when env var is not set) ─────────────────────
RATE_LIMIT_LOGIN_PER_MINUTE="${RATE_LIMIT_LOGIN_PER_MINUTE:-20}"
RATE_LIMIT_PASSWORD_PER_MINUTE="${RATE_LIMIT_PASSWORD_PER_MINUTE:-10}"
RATE_LIMIT_PER_MINUTE="${RATE_LIMIT_PER_MINUTE:-100}"
RATE_LIMIT_PER_DAY="${RATE_LIMIT_PER_DAY:-1000}"
OTP_RATE_LIMIT_MAX="${OTP_RATE_LIMIT_MAX:-100}"
OTP_RATE_LIMIT_WINDOW_MINUTES="${OTP_RATE_LIMIT_WINDOW_MINUTES:-60}"

# ── Output dir ────────────────────────────────────────────────────────────────
mkdir -p "$OUT_DIR"

echo ""
echo "🔧 Rendering nginx conf.d templates → $OUT_DIR"
echo "   login rate limit:    ${RATE_LIMIT_LOGIN_PER_MINUTE}r/m"
echo "   password rate limit: ${RATE_LIMIT_PASSWORD_PER_MINUTE}r/m"
echo "   OTP rate limit:      ${OTP_RATE_LIMIT_MAX}/hour"
echo ""

RENDERED=0
SKIPPED=0

for TEMPLATE in "$NGINX_DIR/conf.d/"*.conf; do
    [ -f "$TEMPLATE" ] || continue
    BASENAME="$(basename "$TEMPLATE")"
    DEST="$OUT_DIR/$BASENAME"

    # Check if file has any __PLACEHOLDER__ tokens
    if grep -q '__[A-Z_]*__' "$TEMPLATE" 2>/dev/null; then
        # Render via sed — substitute all known placeholders
        sed \
            -e "s|__RATE_LIMIT_LOGIN_PER_MINUTE__|${RATE_LIMIT_LOGIN_PER_MINUTE}|g" \
            -e "s|__RATE_LIMIT_PASSWORD_PER_MINUTE__|${RATE_LIMIT_PASSWORD_PER_MINUTE}|g" \
            -e "s|__RATE_LIMIT_PER_MINUTE__|${RATE_LIMIT_PER_MINUTE}|g" \
            -e "s|__RATE_LIMIT_PER_DAY__|${RATE_LIMIT_PER_DAY}|g" \
            -e "s|__OTP_RATE_LIMIT_MAX__|${OTP_RATE_LIMIT_MAX}|g" \
            -e "s|__OTP_RATE_LIMIT_WINDOW_MINUTES__|${OTP_RATE_LIMIT_WINDOW_MINUTES}|g" \
            "$TEMPLATE" > "$DEST"

        # Warn if any unresolved placeholders remain
        if grep -q '__[A-Z_]*__' "$DEST" 2>/dev/null; then
            REMAINING=$(grep -o '__[A-Z_]*__' "$DEST" | sort -u | tr '\n' ' ')
            echo "  ⚠  $BASENAME: unresolved placeholders: $REMAINING"
        else
            echo "  ✅ $BASENAME (rendered)"
        fi
        RENDERED=$((RENDERED + 1))
    else
        # No placeholders — copy as-is
        cp "$TEMPLATE" "$DEST"
        echo "  📄 $BASENAME (copied)"
        SKIPPED=$((SKIPPED + 1))
    fi
done

echo ""
echo "Done. Rendered: $RENDERED  Copied: $SKIPPED"
echo "Output: $OUT_DIR"
echo ""
echo "Next steps:"
echo "  Local test:   sudo cp $OUT_DIR/*.conf /etc/nginx/conf.d/ && sudo nginx -t"
echo "  Deploy:       ./scripts/quickerdeploy.sh  (uses rendered output automatically)"
