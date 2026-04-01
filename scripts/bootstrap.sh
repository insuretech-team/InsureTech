#!/usr/bin/env bash
# bootstrap.sh — InsureTech prerequisite checker & auto-installer
# Source from any script: source "$(dirname "$0")/bootstrap.sh"
# Or run standalone:      ./scripts/bootstrap.sh [--go-only] [--python-only] [--node-only]
#
# Checks and installs (when possible):
#   go, buf, protoc-gen-go, protoc-gen-go-grpc, python3, node, npm, ssh/scp
#
# Works on: Ubuntu/Debian (Codespaces), macOS (Homebrew), Alpine (Docker)

set -euo pipefail

# ── Parse flags ───────────────────────────────────────────────────────────────
GO_ONLY=false
PYTHON_ONLY=false
NODE_ONLY=false
ALL=true

for arg in "$@"; do
    case "$arg" in
        --go-only)     GO_ONLY=true;     ALL=false ;;
        --python-only) PYTHON_ONLY=true; ALL=false ;;
        --node-only)   NODE_ONLY=true;   ALL=false ;;
        --all)         ALL=true ;;
        *) echo "Unknown flag: $arg" >&2; exit 1 ;;
    esac
done

# ── Detect OS / package manager ───────────────────────────────────────────────
OS="$(uname -s)"
ARCH="$(uname -m)"
PKG_MGR=""

if command -v apt-get &>/dev/null; then
    PKG_MGR="apt"
elif command -v brew &>/dev/null; then
    PKG_MGR="brew"
elif command -v apk &>/dev/null; then
    PKG_MGR="apk"
elif command -v dnf &>/dev/null; then
    PKG_MGR="dnf"
elif command -v yum &>/dev/null; then
    PKG_MGR="yum"
fi

# ── Colour helpers ────────────────────────────────────────────────────────────
GREEN='\033[0;32m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; CYAN='\033[0;36m'; RESET='\033[0m'
check() { echo -e "  ${GREEN}✓${RESET} $*"; }
warn()  { echo -e "  ${YELLOW}⚠${RESET} $*"; }
info()  { echo -e "  ${CYAN}→${RESET} $*"; }
fail()  { echo -e "  ${RED}✗${RESET} $*" >&2; exit 1; }

# ── Sudo helper (no-op if already root) ───────────────────────────────────────
maybe_sudo() {
    if [ "$(id -u)" -eq 0 ]; then
        "$@"
    else
        sudo "$@"
    fi
}

# ── Ensure GOPATH is user-writable, then refresh PATH ────────────────────────
update_go_path() {
    local gopath
    gopath="$(go env GOPATH 2>/dev/null || true)"

    # If GOPATH is root-owned (e.g. /go in Codespaces), override to ~/go
    if [ -n "$gopath" ] && [ ! -w "$gopath" ]; then
        export GOPATH="$HOME/go"
        go env -w GOPATH="$HOME/go" 2>/dev/null || true
        gopath="$HOME/go"
        info "GOPATH was root-owned — switched to $GOPATH"
    fi

    # Create GOPATH/bin if missing
    mkdir -p "$gopath/bin"

    if [[ ":$PATH:" != *":$gopath/bin:"* ]]; then
        export PATH="$gopath/bin:$PATH"
    fi

    # Persist for future shells
    for profile in "$HOME/.bashrc" "$HOME/.profile" "$HOME/.bash_profile" "$HOME/.zshrc"; do
        if [ -f "$profile" ] && ! grep -q 'GOPATH' "$profile" 2>/dev/null; then
            echo "export GOPATH=\"\$HOME/go\"" >> "$profile"
            echo 'export PATH="$GOPATH/bin:$PATH"' >> "$profile"
        fi
    done
}

echo ""
echo -e "${CYAN}=========================================${RESET}"
echo -e "${CYAN} InsureTech Prerequisite Bootstrap${RESET}"
echo -e "${CYAN}=========================================${RESET}"
echo -e " Platform: ${OS} / ${ARCH} / pkg: ${PKG_MGR:-unknown}"
echo ""

# ─────────────────────────────────────────────────────────────────────────────
# GO
# ─────────────────────────────────────────────────────────────────────────────
ensure_go() {
    echo -e "\n${YELLOW}[Go]${RESET}"

    if command -v go &>/dev/null; then
        check "go $(go version | sed 's/go version go//')"
        return
    fi

    info "go not found — installing Go 1.25..."

    local GO_VER="1.25.0"
    local GO_ARCH
    case "$ARCH" in
        x86_64)  GO_ARCH="amd64" ;;
        aarch64|arm64) GO_ARCH="arm64" ;;
        armv6l)  GO_ARCH="armv6l" ;;
        *)       fail "Unsupported architecture: $ARCH — install Go manually from https://go.dev/dl/" ;;
    esac

    if [ "$OS" = "Darwin" ] && [ -n "$PKG_MGR" ] && [ "$PKG_MGR" = "brew" ]; then
        brew install go
        # Ensure homebrew paths are available
        export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH"
    elif [ "$OS" = "Linux" ]; then
        local tarball="go${GO_VER}.linux-${GO_ARCH}.tar.gz"
        info "Downloading $tarball..."
        curl -fsSL "https://go.dev/dl/$tarball" -o "/tmp/$tarball"
        maybe_sudo rm -rf /usr/local/go
        maybe_sudo tar -C /usr/local -xzf "/tmp/$tarball"
        rm -f "/tmp/$tarball"
        export PATH="/usr/local/go/bin:$PATH"
        # Persist for future shells in common profile locations
        for profile in "$HOME/.bashrc" "$HOME/.profile" "$HOME/.bash_profile"; do
            if [ -f "$profile" ] && ! grep -q '/usr/local/go/bin' "$profile"; then
                echo 'export PATH="/usr/local/go/bin:$PATH"' >> "$profile"
            fi
        done
    else
        fail "Cannot auto-install Go on $OS. Download from https://go.dev/dl/"
    fi

    if ! command -v go &>/dev/null; then
        fail "Go installation failed. Install manually: https://go.dev/dl/"
    fi
    check "go $(go version | sed 's/go version go//') (just installed)"
}

# ─────────────────────────────────────────────────────────────────────────────
# GO TOOLS
# ─────────────────────────────────────────────────────────────────────────────
ensure_go_tool() {
    local name="$1"
    local pkg="$2"

    if command -v "$name" &>/dev/null; then
        check "$name"
        return
    fi

    # Ensure GOPATH is writable before attempting install
    update_go_path

    info "$name not found — installing via go install..."
    GOPATH="${GOPATH:-$HOME/go}" go install "$pkg"

    # Refresh PATH again after install
    update_go_path

    if command -v "$name" &>/dev/null; then
        check "$name (just installed)"
    else
        warn "$name installed but not on PATH. Add \$(go env GOPATH)/bin to your PATH."
        warn "  Run: echo 'export PATH=\"\$(go env GOPATH)/bin:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
    fi
}

ensure_go_tools() {
    echo -e "\n${YELLOW}[Go tools]${RESET}"
    # Fix GOPATH first — must be writable before any go install
    update_go_path

    # Pinned versions — avoids re-downloading on every Codespace rebuild.
    # Update these when intentionally upgrading a tool.
    local BUF_VERSION="v1.50.0"
    local PROTOC_GEN_GO_VERSION="v1.36.6"
    local PROTOC_GEN_GO_GRPC_VERSION="v1.5.1"

    ensure_go_tool "buf"                "github.com/bufbuild/buf/cmd/buf@${BUF_VERSION}"
    ensure_go_tool "protoc-gen-go"      "google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}"
    ensure_go_tool "protoc-gen-go-grpc" "google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}"
}

# ─────────────────────────────────────────────────────────────────────────────
# PYTHON
# ─────────────────────────────────────────────────────────────────────────────
ensure_python() {
    echo -e "\n${YELLOW}[Python]${RESET}"

    local py_cmd=""
    for candidate in python3 python; do
        if command -v "$candidate" &>/dev/null; then
            local ver
            ver=$("$candidate" --version 2>&1 | sed 's/Python 3\.//' | sed 's/\..*//' || true)
            if [ -n "$ver" ] && [ "$ver" -ge 8 ] 2>/dev/null; then
                py_cmd="$candidate"
                break
            fi
        fi
    done

    if [ -n "$py_cmd" ]; then
        check "$py_cmd → $($py_cmd --version 2>&1)"
        ensure_pip "$py_cmd"
        ensure_pip_packages "$py_cmd"
        return
    fi

    info "Python 3 not found — installing..."

    case "$PKG_MGR" in
        apt)
            maybe_sudo apt-get update -qq
            maybe_sudo apt-get install -y python3 python3-pip python3-venv
            ;;
        brew)
            brew install python@3.12
            ;;
        apk)
            maybe_sudo apk add --no-cache python3 py3-pip
            ;;
        dnf|yum)
            maybe_sudo "$PKG_MGR" install -y python3 python3-pip
            ;;
        *)
            fail "Cannot auto-install Python. Download from https://www.python.org/downloads/"
            ;;
    esac

    for candidate in python3 python; do
        if command -v "$candidate" &>/dev/null; then py_cmd="$candidate"; break; fi
    done

    if [ -n "$py_cmd" ]; then
        check "$py_cmd (just installed)"
        ensure_pip "$py_cmd"
        ensure_pip_packages "$py_cmd"
    else
        warn "Python installed but not detected — you may need to restart your shell."
    fi
}

ensure_pip() {
    local py_cmd="$1"

    if "$py_cmd" -m pip --version &>/dev/null; then
        check "pip"
        return
    fi

    info "pip not found — bootstrapping..."

    if "$py_cmd" -m ensurepip --upgrade &>/dev/null; then
        info "pip installed via ensurepip"
    fi

    if ! "$py_cmd" -m pip --version &>/dev/null; then
        case "$PKG_MGR" in
            apt)
                maybe_sudo apt-get update -qq
                maybe_sudo apt-get install -y python3-pip python3-venv
                ;;
            brew)
                brew install python@3.12
                ;;
            apk)
                maybe_sudo apk add --no-cache py3-pip
                ;;
            dnf|yum)
                maybe_sudo "$PKG_MGR" install -y python3-pip
                ;;
        esac
    fi

    if "$py_cmd" -m pip --version &>/dev/null; then
        check "pip"
    else
        fail "pip is unavailable for $py_cmd. Install pip manually, then re-run bootstrap."
    fi
}

python_can_import() {
    local py_cmd="$1"
    local pkg="$2"

    case "$pkg" in
        ruamel.yaml)
            "$py_cmd" -c "import ruamel.yaml" &>/dev/null
            ;;
        pyyaml)
            "$py_cmd" -c "import yaml" &>/dev/null
            ;;
        protobuf)
            "$py_cmd" -c "from google import protobuf" &>/dev/null
            ;;
        *)
            "$py_cmd" -c "import $pkg" &>/dev/null
            ;;
    esac
}

install_pip_packages() {
    local py_cmd="$1"
    shift

    # Use --progress-bar=off to prevent hang on non-TTY terminals (Codespaces).
    # Try --break-system-packages first (PEP 668 / Debian/Ubuntu 23+),
    # then --user, then plain install.
    if "$py_cmd" -m pip install --progress-bar=off --break-system-packages "$@" 2>&1; then
        return 0
    fi

    if "$py_cmd" -m pip install --progress-bar=off --user "$@" 2>&1; then
        return 0
    fi

    "$py_cmd" -m pip install --progress-bar=off "$@" 2>&1
}

ensure_pip_packages() {
    local py_cmd="$1"
    # Packages required by run_api_pipeline generator scripts
    # Note: bash array syntax is used here; this script requires bash (see shebang)
    local packages=("ruamel.yaml" "pyyaml" "python-dotenv" "requests" "protobuf")
    local missing=()

    for pkg in "${packages[@]}"; do
        if ! python_can_import "$py_cmd" "$pkg"; then
            missing+=("$pkg")
        fi
    done

    if [ "${#missing[@]}" -gt 0 ]; then
        info "Installing Python packages: ${missing[*]}"
        install_pip_packages "$py_cmd" "${missing[@]}"
        check "Python packages installed"
    else
        check "Python packages OK"
    fi
}

# ─────────────────────────────────────────────────────────────────────────────
# NODE / NPM
# ─────────────────────────────────────────────────────────────────────────────
ensure_node() {
    echo -e "\n${YELLOW}[Node.js / npm]${RESET}"

    if command -v node &>/dev/null && command -v npm &>/dev/null; then
        local node_ver
        node_ver=$(node --version | sed 's/v//' | cut -d. -f1)
        if [ "$node_ver" -ge 18 ] 2>/dev/null; then
            check "node $(node --version)  /  npm v$(npm --version)"
            return
        fi
        warn "node $(node --version) found but Node 18+ required — upgrading..."
    else
        info "Node.js not found — installing..."
    fi

    case "$PKG_MGR" in
        apt)
            # Use NodeSource for a modern LTS (with curl/wget fallback for minimal Alpine)
            if command -v curl &>/dev/null; then
                curl -fsSL --max-time 30 --connect-timeout 10 https://deb.nodesource.com/setup_lts.x | maybe_sudo bash -
            elif command -v wget &>/dev/null; then
                wget --timeout=30 -qO- https://deb.nodesource.com/setup_lts.x | maybe_sudo bash -
            else
                warn "Neither curl nor wget found, skipping NodeSource setup"
            fi
            maybe_sudo apt-get install -y nodejs
            ;;
        brew)
            brew install node@20
            brew link --overwrite node@20 2>/dev/null || true
            ;;
        apk)
            maybe_sudo apk add --no-cache nodejs npm
            ;;
        dnf|yum)
            # Use NodeSource for modern version (with curl/wget fallback)
            if command -v curl &>/dev/null; then
                curl -fsSL --max-time 30 --connect-timeout 10 https://rpm.nodesource.com/setup_lts.x | maybe_sudo bash -
            elif command -v wget &>/dev/null; then
                wget --timeout=30 -qO- https://rpm.nodesource.com/setup_lts.x | maybe_sudo bash -
            else
                warn "Neither curl nor wget found, skipping NodeSource setup"
            fi
            maybe_sudo "$PKG_MGR" install -y nodejs
            ;;
        *)
            fail "Cannot auto-install Node.js. Download from https://nodejs.org/"
            ;;
    esac

    if command -v node &>/dev/null; then
        check "node $(node --version) (just installed)"
    else
        warn "Node.js installed but not detected — you may need to restart your shell."
    fi
}

# ─────────────────────────────────────────────────────────────────────────────
# SSH / SCP
# ─────────────────────────────────────────────────────────────────────────────
ensure_ssh() {
    echo -e "\n${YELLOW}[SSH / SCP]${RESET}"

    if command -v ssh &>/dev/null && command -v scp &>/dev/null; then
        check "ssh + scp"
        return
    fi

    info "ssh/scp not found — installing..."

    case "$PKG_MGR" in
        apt)   maybe_sudo apt-get install -y openssh-client ;;
        brew)  check "ssh is bundled with macOS"; return ;;
        apk)   maybe_sudo apk add --no-cache openssh-client ;;
        dnf|yum) maybe_sudo "$PKG_MGR" install -y openssh-clients ;;
        *)     fail "Cannot auto-install SSH. Please install openssh-client manually." ;;
    esac

    check "openssh-client installed"
}

# ─────────────────────────────────────────────────────────────────────────────
# DOCKER  (optional)
# ─────────────────────────────────────────────────────────────────────────────
check_docker() {
    echo -e "\n${YELLOW}[Docker (optional)]${RESET}"
    if command -v docker &>/dev/null; then
        check "docker $(docker version --format '{{.Client.Version}}' 2>/dev/null || echo 'installed')"
    else
        warn "docker not found — OpenAPI Docker-validation step will be skipped (non-fatal)"
        warn "Install from https://docs.docker.com/get-docker/ if you want full validation"
    fi
}

# ─────────────────────────────────────────────────────────────────────────────
# ENTRY POINT
# ─────────────────────────────────────────────────────────────────────────────
if [ "$ALL" = true ] || [ "$GO_ONLY" = true ]; then
    ensure_go
    ensure_go_tools
fi

if [ "$ALL" = true ] || [ "$PYTHON_ONLY" = true ]; then
    ensure_python
fi

if [ "$ALL" = true ] || [ "$NODE_ONLY" = true ]; then
    ensure_node
fi

if [ "$ALL" = true ]; then
    ensure_ssh
    check_docker
fi

echo ""
echo -e "${CYAN}=========================================${RESET}"
echo -e "${GREEN} All prerequisites satisfied!${RESET}"
echo -e "${CYAN}=========================================${RESET}"
echo ""
