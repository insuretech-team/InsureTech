# start_stateful_mock_server.ps1 — Start the stateful Python mock server
#
# Usage:
#   .\scripts\start_stateful_mock_server.ps1 [-Port 4010]

param(
    [int]$Port = 4010
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$ServerScript = Join-Path $ProjectRoot "scripts\stateful_mock_server.py"

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  InsureTech Stateful Mock Server" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  Base URL:       http://localhost:$Port" -ForegroundColor Green
Write-Host "  Reset state:    POST http://localhost:$Port/_mock/reset" -ForegroundColor Green
Write-Host "  Inspect state:  GET  http://localhost:$Port/_mock/state" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

python $ServerScript --port $Port
