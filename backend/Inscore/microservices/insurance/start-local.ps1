#!/usr/bin/env pwsh
# Start Insurance Service with environment from .env

Write-Host "Starting Insurance Service..." -ForegroundColor Cyan
Write-Host ""

# Load .env file from project root
$envPath = "..\..\..\..\.env"
if (Test-Path $envPath) {
    Write-Host "Loading environment from .env..." -ForegroundColor Yellow
    Get-Content $envPath | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim().Trim('"').Trim("'")
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
    Write-Host "✓ Environment loaded" -ForegroundColor Green
} else {
    Write-Host "✗ .env file not found at $envPath" -ForegroundColor Red
    Write-Host "Service will use default/system environment variables" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Starting service on port 50115..." -ForegroundColor Green
Write-Host ""

.\insurance-service.exe
