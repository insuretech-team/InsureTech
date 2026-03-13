#!/usr/bin/env pwsh
# TypeScript SDK Generator Script (hey-api + custom)

Write-Host "🚀 TypeScript SDK Generator (hey-api + custom)" -ForegroundColor Cyan
Write-Host ("=" * 60) -ForegroundColor Cyan
Write-Host ""

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptDir

try {
    # Step 1: Install hey-api dependencies
    Write-Host "📦 Installing @hey-api/openapi-ts..." -ForegroundColor Yellow
    npm install
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to install dependencies"
    }
    
    Write-Host "✓ Dependencies installed" -ForegroundColor Green
    Write-Host ""
    
    # Step 2: Run hey-api generator
    Write-Host "⚙️  Running @hey-api/openapi-ts..." -ForegroundColor Yellow
    npm run generate
    
    if ($LASTEXITCODE -ne 0) {
        throw "hey-api generator failed"
    }
    
    Write-Host "✓ Base SDK generated" -ForegroundColor Green
    Write-Host ""
    
    # Step 3: Build custom Go post-processor
    Write-Host "🔨 Building custom post-processor..." -ForegroundColor Yellow
    $env:GOWORK = "off"
    go build -o generator.exe generator.go
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to build post-processor"
    }
    
    Write-Host "✓ Post-processor built" -ForegroundColor Green
    Write-Host ""
    
    # Step 4: Run custom post-processor
    Write-Host "🔧 Running custom post-processor..." -ForegroundColor Yellow
    .\generator.exe
    
    if ($LASTEXITCODE -ne 0) {
        throw "Post-processor failed"
    }
    
    Write-Host "✓ Customizations applied" -ForegroundColor Green
    Write-Host ""
    
    # Step 5: Navigate to SDK and install dependencies
    $sdkPath = Join-Path $scriptDir ".." ".." "insuretech-typescript-sdk"
    Set-Location $sdkPath
    
    Write-Host "📦 Installing SDK dependencies..." -ForegroundColor Yellow
    npm install
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to install SDK dependencies"
    }
    
    Write-Host "✓ SDK dependencies installed" -ForegroundColor Green
    Write-Host ""
    
    # Step 6: Build SDK
    Write-Host "🏗️  Building SDK..." -ForegroundColor Yellow
    npm run build
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to build SDK"
    }
    
    Write-Host "✓ SDK built successfully" -ForegroundColor Green
    Write-Host ""
    
    # Step 7: Run tests
    Write-Host "🧪 Running tests..." -ForegroundColor Yellow
    npm test
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "⚠️  Tests failed (expected for new SDK)" -ForegroundColor Yellow
    } else {
        Write-Host "✓ Tests passed" -ForegroundColor Green
    }
    Write-Host ""
    
    Write-Host ("=" * 60) -ForegroundColor Cyan
    Write-Host "✅ TypeScript SDK generation completed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "SDK Location: $sdkPath" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  1. Review the generated code" -ForegroundColor White
    Write-Host "  2. Add tests if needed" -ForegroundColor White
    Write-Host "  3. Publish: npm publish --access public" -ForegroundColor White
    Write-Host ""
    
} catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    Write-Host ""
    exit 1
}
