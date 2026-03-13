# Generate Go SDK from OpenAPI spec
# This script generates the InsureTech Go SDK

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  InsureTech Go SDK Generator" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$ErrorActionPreference = "Stop"

# Get the script directory and project root
$SCRIPT_DIR = Split-Path -Parent $MyInvocation.MyCommand.Path
# Go up 3 levels: go -> sdk-generator -> sdks -> InsureTech (project root)
$PROJECT_ROOT = Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $SCRIPT_DIR))

# Configuration - using relative paths from project root
$PROTO_PATH = Join-Path $PROJECT_ROOT "proto"
$API_SPEC_PATH = Join-Path $PROJECT_ROOT "api\openapi.yaml"
$OUTPUT_PATH = Join-Path $PROJECT_ROOT "sdks\insuretech-go-sdk"
$GENERATOR_PATH = $SCRIPT_DIR

Write-Host "Project Root: $PROJECT_ROOT" -ForegroundColor Cyan
Write-Host "Script Directory: $SCRIPT_DIR" -ForegroundColor Cyan
Write-Host ""

# Validate paths
Write-Host "Validating paths..." -ForegroundColor Yellow

if (-not (Test-Path $PROTO_PATH)) {
    Write-Host "✗ Proto path not found: $PROTO_PATH" -ForegroundColor Red
    exit 1
}
Write-Host "✓ Proto path exists" -ForegroundColor Green

if (-not (Test-Path $API_SPEC_PATH)) {
    Write-Host "✗ OpenAPI spec not found: $API_SPEC_PATH" -ForegroundColor Red
    exit 1
}
Write-Host "✓ OpenAPI spec exists" -ForegroundColor Green

Write-Host ""

# Create output directory if it doesn't exist
if (-not (Test-Path $OUTPUT_PATH)) {
    Write-Host "Creating output directory: $OUTPUT_PATH" -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $OUTPUT_PATH -Force | Out-Null
    Write-Host "✓ Output directory created" -ForegroundColor Green
} else {
    Write-Host "✓ Output directory exists" -ForegroundColor Green
}

Write-Host ""

# Run the generator
Write-Host "Running Go SDK generator..." -ForegroundColor Yellow
Write-Host ""

try {
    Push-Location $GENERATOR_PATH
    
    # Check if generator is built
    if (-not (Test-Path "generator.exe")) {
        Write-Host "Building generator..." -ForegroundColor Yellow
        go build -o generator.exe generator.go
        if ($LASTEXITCODE -ne 0) {
            throw "Failed to build generator"
        }
        Write-Host "✓ Generator built successfully" -ForegroundColor Green
    }
    
    # Run the generator
    ./generator.exe
    
    if ($LASTEXITCODE -ne 0) {
        throw "Generator failed with exit code $LASTEXITCODE"
    }
    
    Pop-Location
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "✓ SDK Generation Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Output location: $OUTPUT_PATH" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  1. cd $OUTPUT_PATH"
    Write-Host "  2. go mod tidy"
    Write-Host "  3. go test ./..."
    Write-Host ""
    
} catch {
    Pop-Location
    Write-Host ""
    Write-Host "✗ Generation failed: $_" -ForegroundColor Red
    exit 1
}
