# start_mock_server.ps1 — Start Prism mock server for local frontend development
#
# Usage:
#   .\scripts\start_mock_server.ps1 [-Port 4010] [-Dynamic]
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

param(
    [int]$Port = 4010,
    [switch]$Dynamic
)

$ErrorActionPreference = "Stop"
$ScriptDir   = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
$OpenApiSpec = Join-Path $ProjectRoot "api\openapi.yaml"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "   InsureTech Prism Mock Server" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

if (-not (Test-Path $OpenApiSpec)) {
    Write-Host "  ⚠ OpenAPI spec not found at: $OpenApiSpec" -ForegroundColor Yellow
    Write-Host "  Run .\run_api_pipeline.ps1 first to generate the spec."
    exit 1
}

Write-Host "  ✓ OpenAPI spec: $OpenApiSpec" -ForegroundColor Green
Write-Host "  ✓ Mock server:  http://localhost:$Port" -ForegroundColor Green
if ($Dynamic) {
    Write-Host "  ✓ Mode: dynamic (randomised responses)" -ForegroundColor Green
} else {
    Write-Host "  ✓ Mode: static (first example from spec)" -ForegroundColor Green
}
Write-Host ""
Write-Host "  Frontend usage:" -ForegroundColor White
Write-Host "    Set your API base URL to: http://localhost:$Port" -ForegroundColor Cyan
Write-Host "    All endpoints return example data from the OpenAPI spec."
Write-Host "    No authentication required (mock server ignores Bearer tokens)."
Write-Host ""
Write-Host "  Press Ctrl+C to stop the mock server." -ForegroundColor Gray
Write-Host ""

$prismVersion = "@stoplight/prism-cli@5.14.2"
$prismArgs = @("--yes", $prismVersion, "mock", $OpenApiSpec, "--port", $Port, "--cors")
if ($Dynamic) {
    $prismArgs += "--dynamic"
}

npx @prismArgs
