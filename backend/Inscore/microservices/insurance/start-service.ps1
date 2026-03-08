#!/usr/bin/env pwsh
# Start Insurance Service with environment variables loaded

Write-Host "=== Starting Insurance Service ===" -ForegroundColor Cyan

# Load .env file from project root
$envFile = "../../../../.env"
if (Test-Path $envFile) {
    Write-Host "Loading environment variables from .env..." -ForegroundColor Yellow
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim().Trim('"').Trim("'")
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
    Write-Host "Environment variables loaded" -ForegroundColor Green
} else {
    Write-Host "Warning: .env file not found at $envFile" -ForegroundColor Yellow
}

# Set Insurance Service specific variables
$env:INSURANCE_GRPC_PORT = "50115"
$env:INSURANCE_HTTP_PORT = "50116"

$port = $env:INSURANCE_GRPC_PORT
Write-Host "Starting Insurance Service on port $port..." -ForegroundColor Cyan

# Run the service
go run main.go
