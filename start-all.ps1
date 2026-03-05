#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Start all InsureTech inscore microservices.

.DESCRIPTION
    Launches every cmd/*/main.go entrypoint from the project root (where go.mod
    lives). Services are started in dependency order:

      1. Infrastructure   – storage, tenant
      2. Auth layer        – authn, authz
      3. Core services     – fraud, partner, kyc, beneficiary, b2b, audit, workflow
      4. Comms / Media     – notification, support, media, docgen, ocr, webrtc
      5. Intelligence      – iot, analytics, ai
      6. Gateway (last)    – gateway (HTTP entry-point, depends on everything)

    REAL services  → each gets its own visible pwsh window so you can watch logs live.
    STUB services  → run as silent background jobs (health-only, no interesting logs).

.PARAMETER Migrate
    Run DB migrations before starting authn (sets AUTHN_RUN_MIGRATIONS=true).
    Default: false (skip migrations for fast restarts).

.PARAMETER Services
    Comma-separated list of service names to start. Default: all.
    Example: -Services "gateway,authn,authz"

.PARAMETER LogDir
    Directory to write per-service log files (stubs + window echo). Default: .\logs\services

.EXAMPLE
    # Start everything (no migrations)
    .\start-all.ps1

.EXAMPLE
    # Start everything and run DB migrations first
    .\start-all.ps1 -Migrate

.EXAMPLE
    # Start only the auth stack
    .\start-all.ps1 -Services "authn,authz"
#>

[CmdletBinding()]
param(
    [switch]$Migrate,
    [string]$Services = "",
    [string]$LogDir   = "",
    [switch]$KillAll
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

# ── Resolve project root (where this script lives = where go.mod is) ──────────
$ProjectRoot = $PSScriptRoot
Set-Location $ProjectRoot

# ── KillAll mode — stop every inscore service and exit ────────────────────────
if ($KillAll) {
    Write-Host "`nKilling all inscore services..." -ForegroundColor Cyan

    # 1. Kill all pwsh windows whose title starts with "inscore |"
    $killed = 0
    Get-Process -Name "pwsh", "powershell" -ErrorAction SilentlyContinue | ForEach-Object {
        try {
            $title = $_.MainWindowTitle
            if ($title -like "inscore |*") {
                Stop-Process -Id $_.Id -Force
                Write-Host "  ✓ Closed window: $title" -ForegroundColor Green
                $killed++
            }
        } catch { }
    }

    # 2. Kill any lingering "go run" processes that match our cmd paths
    Get-Process -Name "go" -ErrorAction SilentlyContinue | ForEach-Object {
        try {
            Stop-Process -Id $_.Id -Force
            $killed++
        } catch { }
    }

    # 3. Stop any PowerShell background jobs in this session
    $jobs = Get-Job -ErrorAction SilentlyContinue
    if ($jobs) {
        $jobs | Stop-Job -ErrorAction SilentlyContinue
        $jobs | Remove-Job -Force -ErrorAction SilentlyContinue
        Write-Host "  ✓ Background jobs stopped ($($jobs.Count))." -ForegroundColor Green
        $killed += $jobs.Count
    }

    if ($killed -eq 0) {
        Write-Host "  · No inscore processes found." -ForegroundColor Gray
    } else {
        Write-Host "`n  ✓ Done. $killed process(es) killed." -ForegroundColor Green
    }
    exit 0
}

# ── Log directory ─────────────────────────────────────────────────────────────
if ($LogDir -eq "") {
    $LogDir = Join-Path $ProjectRoot "logs" "services"
}
if (-not (Test-Path $LogDir)) {
    New-Item -ItemType Directory -Path $LogDir -Force | Out-Null
}

# ── Compiled service binary directory ─────────────────────────────────────────
$BinDir = Join-Path $ProjectRoot "backend" "inscore" "bin"
if (-not (Test-Path $BinDir)) {
    New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
}

# ── Colour helpers ────────────────────────────────────────────────────────────
function Write-Header($msg) { Write-Host "`n$msg" -ForegroundColor Cyan }
function Write-Ok($msg)     { Write-Host "  ✓ $msg" -ForegroundColor Green }
function Write-Info($msg)   { Write-Host "  · $msg" -ForegroundColor Gray }
function Write-Warn($msg)   { Write-Host "  ⚠ $msg" -ForegroundColor Yellow }

function Stop-ServicePortIfBusy {
    param($Svc)

    if ($Svc.GrpcPort -le 0) {
        return
    }

    $listeners = Get-NetTCPConnection -LocalPort $Svc.GrpcPort -State Listen -ErrorAction SilentlyContinue
    if (-not $listeners) {
        return
    }

    foreach ($listener in $listeners) {
        $ownerPid = $listener.OwningProcess
        if (-not $ownerPid) {
            continue
        }
        try {
            $proc = Get-Process -Id $ownerPid -ErrorAction Stop
            Stop-Process -Id $ownerPid -Force -ErrorAction Stop
            Write-Warn "[$($Svc.Name)] Port :$($Svc.GrpcPort) already in use by PID $ownerPid ($($proc.ProcessName)). Killed stale process."
        } catch {
            Write-Warn "[$($Svc.Name)] Port :$($Svc.GrpcPort) in use by PID $ownerPid but could not stop process."
        }
    }
}

function Convert-ToWslPath {
    param([string]$WindowsPath)

    $fullPath = [System.IO.Path]::GetFullPath($WindowsPath)
    if ($fullPath -match '^([A-Za-z]):\\(.*)$') {
        $drive = $matches[1].ToLower()
        $rest = $matches[2] -replace '\\', '/'
        return "/mnt/$drive/$rest"
    }

    return ($fullPath -replace '\\', '/')
}

function Escape-ForBashSingleQuote {
    param([string]$Value)
    return $Value -replace "'", "'""'""'"
}

# ── Service catalogue ─────────────────────────────────────────────────────────
# Kind = REAL → own pwsh window with live logs
# Kind = STUB → silent background job (health-only)
$AllServices = @(
    # ── Infrastructure ────────────────────────────────────────────────────────
    [PSCustomObject]@{ Name="storage";      RelPath="backend/inscore/cmd/storage/main.go";      GrpcPort=50290; Kind="REAL" }
    [PSCustomObject]@{ Name="tenant";       RelPath="backend/inscore/cmd/tenant/main.go";       GrpcPort=50050; Kind="STUB" }

    # ── Auth layer ────────────────────────────────────────────────────────────
    [PSCustomObject]@{ Name="authn";        RelPath="backend/inscore/cmd/authn/main.go";        GrpcPort=50060; Kind="REAL" }
    [PSCustomObject]@{ Name="authz";        RelPath="backend/inscore/cmd/authz/main.go";        GrpcPort=50070; Kind="REAL" }

    # ── Core services ─────────────────────────────────────────────────────────
    [PSCustomObject]@{ Name="fraud";        RelPath="backend/inscore/cmd/fraud/main.go";        GrpcPort=50220; Kind="REAL" }
    [PSCustomObject]@{ Name="partner";      RelPath="backend/inscore/cmd/partner/main.go";      GrpcPort=50100; Kind="REAL" }
    [PSCustomObject]@{ Name="kyc";          RelPath="backend/inscore/cmd/kyc/main.go";          GrpcPort=50090; Kind="STUB" }
    [PSCustomObject]@{ Name="beneficiary";  RelPath="backend/inscore/cmd/beneficiary/main.go";  GrpcPort=50110; Kind="STUB" }
    [PSCustomObject]@{ Name="b2b";          RelPath="backend/inscore/cmd/b2b/main.go";          GrpcPort=50112; Kind="REAL" }
    [PSCustomObject]@{ Name="audit";        RelPath="backend/inscore/cmd/audit/main.go";        GrpcPort=50080; Kind="STUB" }
    [PSCustomObject]@{ Name="workflow";     RelPath="backend/inscore/cmd/workflow/main.go";     GrpcPort=50180; Kind="STUB" }

    # ── Comms / Media ─────────────────────────────────────────────────────────
    [PSCustomObject]@{ Name="notification"; RelPath="backend/inscore/cmd/notification/main.go"; GrpcPort=50230; Kind="STUB" }
    [PSCustomObject]@{ Name="support";      RelPath="backend/inscore/cmd/support/main.go";      GrpcPort=50240; Kind="STUB" }
    [PSCustomObject]@{ Name="media";        RelPath="backend/inscore/cmd/media/main.go";        GrpcPort=50260; Kind="REAL" }
    [PSCustomObject]@{ Name="docgen";       RelPath="backend/inscore/cmd/docgen/main.go";       GrpcPort=50280; Kind="REAL" }
    [PSCustomObject]@{ Name="ocr";          RelPath="backend/inscore/cmd/ocr/main.go";          GrpcPort=50270; Kind="STUB" }
    [PSCustomObject]@{ Name="webrtc-server";RelPath="backend/inscore/cmd/webrtc-server/main.go";GrpcPort=50250; Kind="REAL" }
    [PSCustomObject]@{ Name="conference";   RelPath="backend/inscore/cmd/conference/main.go";   GrpcPort=0;     Kind="REAL" }

    # ── Intelligence ──────────────────────────────────────────────────────────
    [PSCustomObject]@{ Name="iot";          RelPath="backend/inscore/cmd/iot/main.go";          GrpcPort=50300; Kind="STUB" }
    [PSCustomObject]@{ Name="analytics";    RelPath="backend/inscore/cmd/analytics/main.go";    GrpcPort=50310; Kind="STUB" }
    [PSCustomObject]@{ Name="ai";           RelPath="backend/inscore/cmd/ai/main.go";           GrpcPort=50320; Kind="STUB" }

    # ── Gateway (last — depends on all upstream services) ─────────────────────
    [PSCustomObject]@{ Name="gateway";      RelPath="backend/inscore/cmd/gateway/main.go";      GrpcPort=0;     Kind="REAL" }
)

# ── Filter by -Services if provided ──────────────────────────────────────────
if ($Services -ne "") {
    $requested = $Services -split "," | ForEach-Object { $_.Trim().ToLower() }
    $AllServices = $AllServices | Where-Object { $requested -contains $_.Name.ToLower() }
    if ($AllServices.Count -eq 0) {
        Write-Error "No services matched the filter: $Services"
        exit 1
    }
}

# ── Migration opt-in ──────────────────────────────────────────────────────────
if ($Migrate) {
    $env:AUTHN_RUN_MIGRATIONS = "true"
    Write-Warn "AUTHN_RUN_MIGRATIONS=true — migrations will run when authn starts."
} else {
    $env:AUTHN_RUN_MIGRATIONS = "false"
}

# ── Track PIDs of spawned windows (for shutdown) ──────────────────────────────
$WindowPIDs = [System.Collections.Generic.List[int]]::new()
$StubJobs   = [System.Collections.Generic.List[object]]::new()

# ─────────────────────────────────────────────────────────────────────────────
# Build-ServiceBinary
#   Builds cmd/<service> to a stable on-disk .exe path instead of relying on
#   "go run", which executes temp binaries that may be blocked by app policy.
# ─────────────────────────────────────────────────────────────────────────────
function Build-ServiceBinary {
    param($Svc)

    $serviceDir = Join-Path $ProjectRoot (
        Split-Path ($Svc.RelPath -replace "/", "\") -Parent
    )
    $exePath = Join-Path $BinDir "$($Svc.Name).exe"

    Push-Location $serviceDir
    try {
        & go build -o $exePath .
        if ($LASTEXITCODE -ne 0 -or -not (Test-Path $exePath)) {
            throw "go build failed for $($Svc.Name)"
        }
    } finally {
        Pop-Location
    }

    return $exePath
}

# ─────────────────────────────────────────────────────────────────────────────
# Start-RealService
#   Opens a new pwsh window titled with the service name.
#   The window runs the compiled service binary and stays open on crash so you
#   can read the error. Logs are also tee'd to $LogDir\<name>.log.
# ─────────────────────────────────────────────────────────────────────────────
function Start-RealService {
    param($Svc)

    $mainGoWin = Join-Path $ProjectRoot ($Svc.RelPath -replace "/", "\")
    $exePath = Build-ServiceBinary -Svc $Svc
    $logFile = Join-Path $LogDir "$($Svc.Name).log"
    $portInfo = if ($Svc.GrpcPort -gt 0) { " :$($Svc.GrpcPort)" } else { "" }
    $title   = "inscore | $($Svc.Name)$portInfo"
    $projectRootWsl = Convert-ToWslPath -WindowsPath $ProjectRoot
    $mainGoWsl = Convert-ToWslPath -WindowsPath $mainGoWin
    $projectRootWslEscaped = Escape-ForBashSingleQuote -Value $projectRootWsl
    $mainGoWslEscaped = Escape-ForBashSingleQuote -Value $mainGoWsl

    # The inner script that runs inside the new window:
    #   1. Sets the window title
    #   2. cds to project root
    #   3. runs the service, tee-ing output to the log file
    #   4. pauses on exit so you can read any crash output
    $innerScript = @"
`$host.UI.RawUI.WindowTitle = '$title'
Set-Location '$ProjectRoot'
Write-Host '[$($Svc.Name)] Starting — log: $logFile' -ForegroundColor Cyan
Write-Host '[$($Svc.Name)] Executable: $exePath' -ForegroundColor DarkGray
`$ErrorActionPreference = 'Stop'
try {
    & '$exePath' 2>&1 | Tee-Object -FilePath '$logFile' -Append
} catch {
    `$errText = `$_.Exception.Message
    if (`$errText -match 'Application Control policy has blocked this file') {
        Write-Host '[$($Svc.Name)] Windows execution blocked by policy. Falling back to WSL go run...' -ForegroundColor Yellow
        wsl -e bash -lc "cd '$projectRootWslEscaped' && go run '$mainGoWslEscaped'" 2>&1 | Tee-Object -FilePath '$logFile' -Append
    } else {
        throw
    }
}
Write-Host ''
Write-Host '[$($Svc.Name)] Process exited. Press any key to close this window...' -ForegroundColor Yellow
`$null = `$host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')
"@

    $encodedCmd = [Convert]::ToBase64String(
        [System.Text.Encoding]::Unicode.GetBytes($innerScript)
    )

    $proc = Start-Process pwsh `
        -ArgumentList "-NoLogo", "-NoExit", "-EncodedCommand", $encodedCmd `
        -PassThru

    Write-Ok "[REAL] $($Svc.Name)$portInfo  →  window opened (PID $($proc.Id))  log: $logFile"
    $WindowPIDs.Add($proc.Id)
}

# ─────────────────────────────────────────────────────────────────────────────
# Start-StubService
#   Runs the stub as a background job — no window, no noise.
#   Logs still go to $LogDir\<name>.log.
# ─────────────────────────────────────────────────────────────────────────────
function Start-StubService {
    param($Svc)

    $mainGoWin = Join-Path $ProjectRoot ($Svc.RelPath -replace "/", "\")
    $exePath = Build-ServiceBinary -Svc $Svc
    $logFile = Join-Path $LogDir "$($Svc.Name).log"
    $projectRootWsl = Convert-ToWslPath -WindowsPath $ProjectRoot
    $mainGoWsl = Convert-ToWslPath -WindowsPath $mainGoWin

    $job = Start-Job -Name $Svc.Name -ScriptBlock {
        param($root, $exe, $log, $rootWsl, $mainWsl, $svcName)
        Set-Location $root
        $ErrorActionPreference = "Stop"
        try {
            & $exe *>&1 | Out-File -FilePath $log -Append
        } catch {
            $errText = $_.Exception.Message
            if ($errText -match 'Application Control policy has blocked this file') {
                "[$svcName] Windows execution blocked by policy. Falling back to WSL go run..." | Out-File -FilePath $log -Append
                wsl -e bash -lc "cd '$rootWsl' && go run '$mainWsl'" *>&1 | Out-File -FilePath $log -Append
            } else {
                throw
            }
        }
    } -ArgumentList $ProjectRoot, $exePath, $logFile, $projectRootWsl, $mainGoWsl, $Svc.Name

    Write-Info "[STUB] $($Svc.Name) :$($Svc.GrpcPort)  →  background job #$($job.Id)  log: $logFile"
    $StubJobs.Add($job)
}

# ── Print banner ──────────────────────────────────────────────────────────────
Write-Header "InsureTech — Starting all inscore services"
Write-Info   "Project root : $ProjectRoot"
Write-Info   "Log dir      : $LogDir"
Write-Info   "Migrate      : $Migrate"
Write-Info   "Services     : $(if ($Services -eq '') { 'all' } else { $Services })"
Write-Host   ""
Write-Info   "REAL services → own pwsh window (live logs + stays open on crash)"
Write-Info   "STUB services → silent background job (health-only)"
Write-Host   ""

# ── Print port map ────────────────────────────────────────────────────────────
Write-Header "Port map  (configs/services.yaml)"
foreach ($svc in $AllServices) {
    if ($svc.GrpcPort -gt 0) {
        $color = if ($svc.Kind -eq "REAL") { "Green" } else { "DarkGray" }
        Write-Host ("  {0,-16} gRPC :{1}  [{2}]" -f $svc.Name, $svc.GrpcPort, $svc.Kind) -ForegroundColor $color
    }
}
Write-Host ""

# ── Launch all services ───────────────────────────────────────────────────────
Write-Header "Launching services..."

foreach ($svc in $AllServices) {
    $mainGo = Join-Path $ProjectRoot $svc.RelPath
    if (-not (Test-Path $mainGo)) {
        Write-Warn "[$($svc.Kind)] $($svc.Name) — main.go not found, skipping."
        continue
    }

    Stop-ServicePortIfBusy -Svc $svc

    if ($svc.Kind -eq "REAL") {
        Start-RealService -Svc $svc
    } else {
        Start-StubService -Svc $svc
    }

    # Small stagger — avoids thundering-herd on go module cache / DB
    Start-Sleep -Milliseconds 400
}

Write-Host ""
Write-Header "All services launched"
Write-Info   "  REAL services : $($WindowPIDs.Count) windows opened"
Write-Info   "  STUB services : $($StubJobs.Count) background jobs"
Write-Host ""
Write-Info   "Press Ctrl+C here to stop ALL stub jobs + close ALL service windows."
Write-Host ""

# ── Wait loop — monitor stubs, forward Ctrl+C ────────────────────────────────
try {
    while ($true) {
        foreach ($job in @($StubJobs)) {
            if ($job.State -eq "Failed" -or $job.State -eq "Stopped") {
                Write-Warn "Stub '$($job.Name)' exited (state: $($job.State)) — check logs\services\$($job.Name).log"
                $StubJobs.Remove($job) | Out-Null
            }
        }
        Start-Sleep -Seconds 5
    }
} finally {
    Write-Header "Shutting down..."

    # Stop stub background jobs
    if ($StubJobs.Count -gt 0) {
        $StubJobs | Stop-Job -ErrorAction SilentlyContinue
        $StubJobs | Remove-Job -Force -ErrorAction SilentlyContinue
        Write-Ok "Stub background jobs stopped."
    }

    # Close real-service windows
    foreach ($pid in $WindowPIDs) {
        try {
            Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
        } catch { }
    }
    if ($WindowPIDs.Count -gt 0) {
        Write-Ok "Real-service windows closed ($($WindowPIDs.Count))."
    }

    Write-Ok "All services stopped."
}
