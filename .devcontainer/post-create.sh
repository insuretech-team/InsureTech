#!/usr/bin/env bash
# post-create.sh — runs once after the Codespace container is created
# Installs all Go tools and Python packages needed for the InsureTech project

set -euo pipefail

GREEN='\033[0;32m'; CYAN='\033[0;36m'; YELLOW='\033[0;33m'; RESET='\033[0m'
ok()   { echo -e "  ${GREEN}✓${RESET} $*"; }
info() { echo -e "  ${CYAN}→${RESET} $*"; }
warn() { echo -e "  ${YELLOW}⚠${RESET} $*"; }

echo ""
echo -e "${CYAN}================================================${RESET}"
echo -e "${CYAN} InsureTech Codespace post-create setup${RESET}"
echo -e "${CYAN}================================================${RESET}"
echo ""

# ── Ensure GOPATH is user-writable ────────────────────────────────────────────
# Codespaces default GOPATH is /go which is root-owned — override to ~/go
if [ ! -w "${GOPATH:-/go}" ]; then
    export GOPATH="$HOME/go"
    go env -w GOPATH="$HOME/go" 2>/dev/null || true
    info "GOPATH was root-owned — switched to $HOME/go"
else
    export GOPATH="${GOPATH:-$HOME/go}"
fi
mkdir -p "$GOPATH/bin"
export PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"

# Persist to shell profiles so every new terminal has correct PATH
for profile in "$HOME/.bashrc" "$HOME/.profile" "$HOME/.bash_profile" "$HOME/.zshrc"; do
    if [ -f "$profile" ]; then
        grep -q 'GOPATH' "$profile" 2>/dev/null || {
            echo "export GOPATH=\"\$HOME/go\"" >> "$profile"
            echo 'export PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"' >> "$profile"
        }
    fi
done

# ── Go version check ──────────────────────────────────────────────────────────
echo "[1/4] Go toolchain"
go version
ok "go $(go version | sed 's/go version go//')"

# ── Go tools ─────────────────────────────────────────────────────────────────
echo ""
echo "[2/4] Installing Go tools..."

# Pinned versions — avoids re-downloading on every Codespace rebuild
# Update these when you intentionally want to upgrade a tool.
BUF_VERSION="v1.50.0"
PROTOC_GEN_GO_VERSION="v1.36.6"
PROTOC_GEN_GO_GRPC_VERSION="v1.5.1"   # Last stable before v1.6 changed package structure

install_go_tool() {
    local name="$1" pkg="$2"
    if command -v "$name" &>/dev/null; then
        ok "$name (already installed — $(${name} --version 2>/dev/null | head -1 || echo 'version unknown'))"
        return 0
    fi
    info "Installing $name from $pkg ..."
    # GONOSUMCHECK + GOFLAGS=-mod=mod avoids sum DB timeouts in restricted networks
    GONOSUMCHECK="*" GOFLAGS="-mod=mod" GOPATH="$GOPATH" go install "$pkg" 2>&1
    # Refresh PATH in case GOPATH/bin wasn't on it yet
    export PATH="$GOPATH/bin:$PATH"
    if command -v "$name" &>/dev/null; then
        ok "$name installed"
    else
        warn "$name installed to $GOPATH/bin but not found in PATH — restart your shell or run: export PATH=\$GOPATH/bin:\$PATH"
    fi
}

install_go_tool "buf"                 "github.com/bufbuild/buf/cmd/buf@${BUF_VERSION}"
install_go_tool "protoc-gen-go"       "google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}"
install_go_tool "protoc-gen-go-grpc"  "google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}"

# ── Python packages ───────────────────────────────────────────────────────────
echo ""
echo "[3/4] Installing Python packages..."
PY_CMD="python3"
command -v python3 &>/dev/null || PY_CMD="python"

# Ensure pip is available — the devcontainer Python feature may omit it
if ! "$PY_CMD" -m pip --version &>/dev/null 2>&1; then
    info "pip not found — bootstrapping via ensurepip..."
    "$PY_CMD" -m ensurepip --upgrade 2>/dev/null || true
fi

# If pip still not available, try installing python3-pip via apt
if ! "$PY_CMD" -m pip --version &>/dev/null 2>&1; then
    info "ensurepip unavailable — installing python3-pip via apt..."
    sudo apt-get install -y python3-pip -qq 2>/dev/null || true
fi

PACKAGES="ruamel.yaml pyyaml requests protobuf"
info "Installing Python packages: $PACKAGES"

# Use --progress-bar=off to avoid hanging on non-TTY Codespace terminals.
# Try --break-system-packages first (PEP 668 / Debian/Ubuntu 23+),
# then --user, then plain install.
if "$PY_CMD" -m pip install --progress-bar=off --break-system-packages $PACKAGES 2>&1; then
    ok "Python packages installed"
elif "$PY_CMD" -m pip install --progress-bar=off --user $PACKAGES 2>&1; then
    ok "Python packages installed (--user)"
else
    "$PY_CMD" -m pip install --progress-bar=off $PACKAGES 2>&1
    ok "Python packages installed"
fi

# ── Copy .env if not present ──────────────────────────────────────────────────
echo ""
echo "[4/4] Environment setup..."
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ ! -f "$REPO_ROOT/.env" ] && [ -f "$REPO_ROOT/.env.example" ]; then
    cp "$REPO_ROOT/.env.example" "$REPO_ROOT/.env"
    info "Copied .env.example → .env (fill in your credentials)"
else
    ok ".env already exists"
fi

echo ""
echo -e "${CYAN}================================================${RESET}"
echo -e "${GREEN} Codespace ready!${RESET}"
echo -e "${CYAN}================================================${RESET}"
echo ""
echo "Quick start:"
echo "  go run ./backend/inscore/cmd/dbops/main.go migrate --target=primary"
echo "  ./run_migration.sh --target=primary"
echo "  ./scripts/generate.sh"
echo ""
