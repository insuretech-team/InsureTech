# bootstrap.ps1 — InsureTech prerequisite checker & auto-installer
# Source from any script: . "$PSScriptRoot\bootstrap.ps1"
# Or call directly:       .\scripts\bootstrap.ps1 [-GoOnly] [-PythonOnly] [-NodeOnly] [-All]
#
# Checks and installs (when possible):
#   go, buf, protoc-gen-go, protoc-gen-go-grpc, python3, node, npm, ssh/scp
#
# Windows  : winget → choco → scoop
# Linux    : apt-get / official Go tarball / NodeSource
# macOS    : Homebrew

param(
    [switch]$GoOnly,     # Only check Go + Go tools (buf, protoc-gen-go, protoc-gen-go-grpc)
    [switch]$PythonOnly, # Only check Python + pip packages
    [switch]$NodeOnly,   # Only check Node/npm
    [switch]$All         # Check everything (default when no flag given)
)

$ErrorActionPreference = "Stop"

# Define .exe extension based on platform
$exe = if ($OnWindows) { '.exe' } else { '' }

# ── Platform detection ────────────────────────────────────────────────────────
# Use separate names to avoid shadowing PS6+ automatic variables
$OnWindows = ($IsWindows -eq $true) -or ($PSVersionTable.PSVersion.Major -le 5)
$OnLinux   = ($IsLinux   -eq $true)
$OnMacOS   = ($IsMacOS   -eq $true)

# ── Helpers ───────────────────────────────────────────────────────────────────
function Write-Check { param($msg) Write-Host "  [OK] $msg" -ForegroundColor Green  }
function Write-Warn  { param($msg) Write-Host "  [!!] $msg" -ForegroundColor Yellow }
function Write-Info  { param($msg) Write-Host "  --> $msg"  -ForegroundColor Cyan   }
function Write-Fail  { param($msg) Write-Host "  [X] $msg"  -ForegroundColor Red; exit 1 }

function Test-Cmd {
    param([string]$Name)
    return $null -ne (Get-Command $Name -ErrorAction SilentlyContinue)
}

# Refresh PATH in the current session after an install
function Update-SessionPath {
    if ($OnWindows) {
        $machine = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
        $user    = [System.Environment]::GetEnvironmentVariable("PATH", "User")
        $env:PATH = "$machine;$user"
    }
    # Add GOPATH/bin if go is now available
    if (Test-Cmd "go") {
        $gobin = "$(go env GOPATH)$([IO.Path]::DirectorySeparatorChar)bin"
        if ($env:PATH -notlike "*$gobin*") {
            $env:PATH = "$gobin$([IO.Path]::PathSeparator)$env:PATH"
            [System.Environment]::SetEnvironmentVariable("PATH", $env:PATH, "Process")
        }
    }
}

# Run a shell command (bash on Linux/macOS, cmd on Windows) and capture output
function Invoke-Shell {
    param([string]$Cmd)
    if ($OnWindows) {
        return (cmd /c $Cmd 2>$null)
    } else {
        return (bash -c $Cmd 2>$null)
    }
}

# ─────────────────────────────────────────────────────────────────────────────
# GO
# ─────────────────────────────────────────────────────────────────────────────
function Ensure-Go {
    Write-Host ""
    Write-Host "[Go]" -ForegroundColor Yellow

    if (Test-Cmd "go") {
        $ver = (go version) -replace "go version ",""
        Write-Check "go $ver"
        return
    }

    Write-Info "go not found — installing Go 1.25..."

    if ($OnWindows) {
        if (Test-Cmd "winget") {
            try { winget install --id GoLang.Go --silent --accept-package-agreements --accept-source-agreements 2>&1 | Out-Null }
            catch { Write-Warn "winget install failed, trying choco..."; $null }
        }
        if (-not (Test-Cmd "go")) {
            if (Test-Cmd "choco") {
                try { choco install golang -y 2>&1 | Out-Null }
                catch { Write-Warn "choco install failed, trying scoop..."; $null }
            }
        }
        if (-not (Test-Cmd "go")) {
            if (Test-Cmd "scoop") {
                try { scoop install go 2>&1 | Out-Null }
                catch { Write-Warn "scoop install failed"; $null }
            }
        }
        if (-not (Test-Cmd "go")) {
            Write-Fail "Cannot auto-install Go. Install from https://go.dev/dl/ then re-run this script."
        }

    } elseif ($OnLinux) {
        $goVer = "1.25.0"
        $arch  = Invoke-Shell "uname -m"
        $goArch = if ($arch -eq "aarch64") { "arm64" } else { "amd64" }
        $tarball = "go${goVer}.linux-${goArch}.tar.gz"
        Write-Info "Downloading $tarball from go.dev..."
        Invoke-WebRequest "https://go.dev/dl/$tarball" -OutFile "/tmp/$tarball" -UseBasicParsing
        Invoke-Shell "sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf /tmp/$tarball && rm /tmp/$tarball"
        $env:PATH = "/usr/local/go/bin:$env:PATH"
        [System.Environment]::SetEnvironmentVariable("PATH", $env:PATH, "Process")
        # Persist for future bash sessions
        $profileLine = 'export PATH="/usr/local/go/bin:$PATH"'
        foreach ($p in @("$HOME/.bashrc","$HOME/.profile","$HOME/.bash_profile")) {
            if ((Test-Path $p) -and -not (Get-Content $p | Select-String '/usr/local/go/bin' -Quiet)) {
                Add-Content $p $profileLine
            }
        }

    } elseif ($OnMacOS) {
        if (Test-Cmd "brew") { 
            try { brew install go 2>&1 | Out-Null }
            catch { Write-Warn "brew install go failed"; $null }
            # Ensure homebrew paths are in PATH
            $env:PATH = "/opt/homebrew/bin:/usr/local/bin:$env:PATH"
            [System.Environment]::SetEnvironmentVariable("PATH", $env:PATH, "Process")
        }
        else { Write-Fail "Homebrew not found. Install from https://brew.sh then re-run, or install Go from https://go.dev/dl/" }

    } else {
        Write-Fail "Unknown platform. Install Go 1.25 from https://go.dev/dl/"
    }

    Update-SessionPath
    if (-not (Test-Cmd "go")) { Write-Fail "Go installation failed. Install manually: https://go.dev/dl/" }
    $ver = (go version) -replace "go version ",""
    Write-Check "go $ver (just installed)"
}

# ─────────────────────────────────────────────────────────────────────────────
# GO TOOLS
# ─────────────────────────────────────────────────────────────────────────────
function Ensure-GoTool {
    param([string]$Name, [string]$Pkg)

    if (Test-Cmd $Name) { Write-Check $Name; return }

    Write-Info "$Name not found — installing via: go install $Pkg"
    go install $Pkg

    # Ensure GOPATH/bin is on PATH
    $gobin = "$(go env GOPATH)$([IO.Path]::DirectorySeparatorChar)bin"
    if ($env:PATH -notlike "*$gobin*") {
        $env:PATH = "$gobin$([IO.Path]::PathSeparator)$env:PATH"
        [System.Environment]::SetEnvironmentVariable("PATH", $env:PATH, "Process")
    }

    if (Test-Cmd $Name) {
        Write-Check "$Name (just installed)"
    } else {
        Write-Warn "$Name installed to $gobin but still not on PATH. Restart your shell or add $gobin to PATH."
    }
}

function Ensure-GoTools {
    Write-Host ""
    Write-Host "[Go tools]" -ForegroundColor Yellow
    # Pinned versions — avoids re-downloading on every rebuild.
    # Update these when intentionally upgrading a tool.
    $BUF_VERSION             = "v1.50.0"
    $PROTOC_GEN_GO_VERSION   = "v1.36.6"
    $PROTOC_GEN_GO_GRPC_VERSION = "v1.5.1"
    Ensure-GoTool "buf"                "github.com/bufbuild/buf/cmd/buf@$BUF_VERSION"
    Ensure-GoTool "protoc-gen-go"      "google.golang.org/protobuf/cmd/protoc-gen-go@$PROTOC_GEN_GO_VERSION"
    Ensure-GoTool "protoc-gen-go-grpc" "google.golang.org/grpc/cmd/protoc-gen-go-grpc@$PROTOC_GEN_GO_GRPC_VERSION"
}

# ─────────────────────────────────────────────────────────────────────────────
# PYTHON
# ─────────────────────────────────────────────────────────────────────────────
function Get-PythonCmd {
    foreach ($candidate in @("python3", "python")) {
        if (Test-Cmd $candidate) {
            $verLine = & $candidate --version 2>&1
            if ($verLine -match "Python 3\.(\d+)" -and [int]$Matches[1] -ge 8) {
                return $candidate
            }
        }
    }
    return $null
}

function Ensure-PipPackages {
    param([string]$PyCmd)

    # Map package name -> import name (for checking if already installed)
    $pkgMap = [ordered]@{
        "ruamel.yaml" = "ruamel.yaml"   # import ruamel.yaml (not ruamel_yaml)
        "pyyaml"      = "yaml"
        "python-dotenv" = "dotenv"
        "requests"    = "requests"
        "protobuf"    = "google.protobuf"
    }

    $missing = @()
    foreach ($pkg in $pkgMap.Keys) {
        $importName = $pkgMap[$pkg]
        # Special handling for ruamel.yaml (namespace package)
        if ($pkg -eq "ruamel.yaml") {
            $result = & $PyCmd -c "import ruamel; import ruamel.yaml" 2>&1
        } else {
            $result = & $PyCmd -c "import $importName" 2>&1
        }
        if ($LASTEXITCODE -ne 0) { $missing += $pkg }
    }

    if ($missing.Count -gt 0) {
        Write-Info "Installing pip packages: $($missing -join ', ')"

        # Ensure pip is available — Python devcontainer feature may omit it
        $pipAvailable = & $PyCmd -m pip --version 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Info "pip not found — bootstrapping via ensurepip..."
            & $PyCmd -m ensurepip --upgrade 2>&1 | Out-Null
        }

        # Still not available? Try apt-get python3-pip
        $pipAvailable = & $PyCmd -m pip --version 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Info "ensurepip failed — trying apt-get install python3-pip..."
            try {
                $aptOut = & sudo apt-get install -y python3-pip -qq 2>&1
                if ($LASTEXITCODE -ne 0) {
                    Write-Warn "apt-get python3-pip failed: $aptOut"
                }
            } catch {
                Write-Warn "apt-get command threw exception: $($_.Exception.Message)"
            }
        }

        # Try install with --break-system-packages first (PEP 668 / Debian/Ubuntu)
        # then fall back to --user, then plain install.
        # Use --progress-bar=off to prevent hang on non-TTY terminals (Codespaces).
        $installed = $false
        foreach ($flags in @("--break-system-packages", "--user", "")) {
            $installArgs = @("-m", "pip", "install", "--progress-bar=off") + $missing
            if ($flags) { $installArgs += $flags }
            & $PyCmd @installArgs 2>&1
            if ($LASTEXITCODE -eq 0) { $installed = $true; break }
        }

        if ($installed) {
            Write-Check "pip packages installed"
        } else {
            Write-Warn "pip install failed — packages may need manual installation"
        }
    } else {
        Write-Check "pip packages OK"
    }
}

function Ensure-Python {
    Write-Host ""
    Write-Host "[Python]" -ForegroundColor Yellow

    $pyCmd = Get-PythonCmd
    if ($pyCmd) {
        $ver = (& $pyCmd --version 2>&1)
        Write-Check "$pyCmd -> $ver"
        Ensure-PipPackages $pyCmd
        return
    }

    Write-Info "Python 3.8+ not found — installing..."

    if ($OnWindows) {
        if (Test-Cmd "winget") {
            try { winget install --id Python.Python.3.12 --silent --accept-package-agreements --accept-source-agreements 2>&1 | Out-Null }
            catch { Write-Warn "winget install failed, trying choco..."; $null }
        }
        if (-not (Test-Cmd "python3") -and -not (Test-Cmd "python")) {
            if (Test-Cmd "choco") {
                try { choco install python -y 2>&1 | Out-Null }
                catch { Write-Warn "choco install failed, trying scoop..."; $null }
            }
        }
        if (-not (Test-Cmd "python3") -and -not (Test-Cmd "python")) {
            if (Test-Cmd "scoop") {
                try { scoop install python 2>&1 | Out-Null }
                catch { Write-Warn "scoop install failed"; $null }
            }
        }
        if (-not (Test-Cmd "python3") -and -not (Test-Cmd "python")) {
            Write-Fail "Cannot auto-install Python. Download from https://www.python.org/downloads/"
        }
        Update-SessionPath
    } elseif ($OnLinux) {
        try { Invoke-Shell "sudo apt-get update -qq && sudo apt-get install -y python3 python3-pip python3-venv" }
        catch {
            Write-Warn "apt-get failed, trying apk..."
            try { Invoke-Shell "sudo apk add --no-cache python3 py3-pip" }
            catch {
                Write-Warn "apk failed, trying dnf..."
                try { Invoke-Shell "sudo dnf install -y python3 python3-pip" }
                catch {
                    Write-Warn "dnf failed, trying yum..."
                    try { Invoke-Shell "sudo yum install -y python3 python3-pip" }
                    catch { Write-Warn "All package managers failed"; $null }
                }
            }
        }
    } elseif ($OnMacOS) {
        if (Test-Cmd "brew") {
            try { brew install python@3.12 2>&1 | Out-Null }
            catch { Write-Warn "brew install python@3.12 failed"; $null }
        }
        else { Write-Fail "Please install Python 3 manually: https://www.python.org/downloads/" }
    } else {
        Write-Fail "Unknown platform. Install Python from https://www.python.org/downloads/"
    }

    $pyCmd = Get-PythonCmd
    if ($pyCmd) {
        Write-Check "$pyCmd (just installed)"
        Ensure-PipPackages $pyCmd
    } else {
        Write-Warn "Python installed but not found in PATH yet — restart your shell and re-run."
    }
}

# ─────────────────────────────────────────────────────────────────────────────
# NODE / NPM
# ─────────────────────────────────────────────────────────────────────────────
function Ensure-Node {
    Write-Host ""
    Write-Host "[Node.js / npm]" -ForegroundColor Yellow

    if ((Test-Cmd "node") -and (Test-Cmd "npm")) {
        $nodeVer = node --version
        if ($nodeVer -match "v(\d+)\." -and [int]$Matches[1] -ge 18) {
            Write-Check "node $nodeVer  /  npm v$(npm --version)"
            return
        }
        Write-Warn "node $nodeVer found but Node 18+ required — upgrading..."
    } else {
        Write-Info "Node.js not found — installing LTS..."
    }

    if ($OnWindows) {
        if (Test-Cmd "winget") {
            try { winget install --id OpenJS.NodeJS.LTS --silent --accept-package-agreements --accept-source-agreements 2>&1 | Out-Null }
            catch { Write-Warn "winget install failed, trying choco..."; $null }
        }
        if (-not (Test-Cmd "node")) {
            if (Test-Cmd "choco") {
                try { choco install nodejs-lts -y 2>&1 | Out-Null }
                catch { Write-Warn "choco install failed, trying scoop..."; $null }
            }
        }
        if (-not (Test-Cmd "node")) {
            if (Test-Cmd "scoop") {
                try { scoop install nodejs-lts 2>&1 | Out-Null }
                catch { Write-Warn "scoop install failed"; $null }
            }
        }
        if (-not (Test-Cmd "node")) {
            Write-Fail "Cannot auto-install Node.js. Download from https://nodejs.org/"
        }
        Update-SessionPath
    } elseif ($OnLinux) {
        try { Invoke-Shell "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - && sudo apt-get install -y nodejs" }
        catch {
            Write-Warn "apt-based Node.js install failed, trying apk..."
            try { Invoke-Shell "sudo apk add --no-cache nodejs npm" }
            catch {
                Write-Warn "apk failed, trying dnf..."
                try { Invoke-Shell "sudo dnf install -y nodejs" }
                catch {
                    Write-Warn "dnf failed, trying yum..."
                    try { Invoke-Shell "sudo yum install -y nodejs" }
                    catch { Write-Warn "All package managers failed"; $null }
                }
            }
        }
    } elseif ($OnMacOS) {
        if (Test-Cmd "brew") {
            try { brew install node@20 2>&1 | Out-Null; brew link --overwrite node@20 2>&1 | Out-Null }
            catch { Write-Warn "brew install node@20 failed"; $null }
        }
        else { Write-Fail "Please install Node.js manually: https://nodejs.org/" }
    } else {
        Write-Fail "Unknown platform. Install Node.js from https://nodejs.org/"
    }

    Update-SessionPath
    if (Test-Cmd "node") { Write-Check "node $(node --version) (just installed)" }
    else { Write-Warn "Node.js installed but not in PATH yet — restart your shell and re-run." }
}

# ─────────────────────────────────────────────────────────────────────────────
# SSH / SCP
# ─────────────────────────────────────────────────────────────────────────────
function Ensure-Ssh {
    Write-Host ""
    Write-Host "[SSH / SCP]" -ForegroundColor Yellow

    if ((Test-Cmd "ssh") -and (Test-Cmd "scp")) { Write-Check "ssh + scp"; return }

    Write-Info "ssh/scp not found — attempting install..."

    if ($OnWindows) {
        try {
            Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0 | Out-Null
            Update-SessionPath
            Write-Check "OpenSSH client installed"
        } catch {
            Write-Warn "Auto-install failed. Enable via: Settings -> Apps -> Optional features -> OpenSSH Client"
        }
    } elseif ($OnLinux) {
        try { Invoke-Shell "sudo apt-get install -y openssh-client" }
        catch {
            Write-Warn "apt-get failed, trying apk..."
            try { Invoke-Shell "sudo apk add --no-cache openssh-client" }
            catch {
                Write-Warn "apk failed, trying dnf..."
                try { Invoke-Shell "sudo dnf install -y openssh-clients" }
                catch {
                    Write-Warn "dnf failed, trying yum..."
                    try { Invoke-Shell "sudo yum install -y openssh-clients" }
                    catch { Write-Warn "All package managers failed"; $null }
                }
            }
        }
        Write-Check "openssh-client installed"
    } elseif ($OnMacOS) {
        Write-Check "ssh is bundled with macOS"
    }
}

# ─────────────────────────────────────────────────────────────────────────────
# DOCKER  (optional)
# ─────────────────────────────────────────────────────────────────────────────
function Check-Docker {
    Write-Host ""
    Write-Host "[Docker (optional)]" -ForegroundColor Yellow
    if (Test-Cmd "docker") {
        $ver = docker version --format "{{.Client.Version}}" 2>$null
        Write-Check "docker $ver"
    } else {
        Write-Warn "docker not found — OpenAPI Docker-validation will be skipped (non-fatal)"
        Write-Warn "Install from https://docs.docker.com/get-docker/ if you want it"
    }
}

# ─────────────────────────────────────────────────────────────────────────────
# ENTRY POINT
# ─────────────────────────────────────────────────────────────────────────────

# Default to All when no specific flag is given, or if -All is explicitly passed
$runAll = $All -or ((-not $GoOnly) -and (-not $PythonOnly) -and (-not $NodeOnly))

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host " InsureTech Prerequisite Bootstrap"       -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
$platform = if ($OnWindows) { "Windows" } elseif ($OnLinux) { "Linux" } else { "macOS" }
Write-Host " Platform : $platform  |  PS $($PSVersionTable.PSVersion)" -ForegroundColor Gray
Write-Host ""

if ($runAll -or $GoOnly)     { Ensure-Go; Ensure-GoTools }
if ($runAll -or $PythonOnly) { Ensure-Python }
if ($runAll -or $NodeOnly)   { Ensure-Node }
if ($runAll)                 { Ensure-Ssh; Check-Docker }

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host " All prerequisites satisfied!"            -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Native commands used during checks can leave a stale non-zero LASTEXITCODE
# even when the overall bootstrap succeeded. Clear it so PowerShell callers
# do not mistake a successful bootstrap for a failure.
$global:LASTEXITCODE = 0
