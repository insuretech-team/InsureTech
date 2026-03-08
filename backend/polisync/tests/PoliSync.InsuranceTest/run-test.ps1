#!/usr/bin/env pwsh
# Run InsuranceTest with environment variables from .env

Write-Host "Loading environment from .env..." -ForegroundColor Cyan

# Load .env file from project root
$envPath = "..\..\..\..\.env"
if (Test-Path $envPath) {
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
    exit 1
}

Write-Host ""
Write-Host "Running InsuranceTest..." -ForegroundColor Cyan
Write-Host ""

dotnet run
