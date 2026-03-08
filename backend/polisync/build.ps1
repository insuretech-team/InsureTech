#!/usr/bin/env pwsh
# PoliSync Build Script

param(
    [Parameter()]
    [ValidateSet('Debug', 'Release')]
    [string]$Configuration = 'Release',
    
    [Parameter()]
    [switch]$Clean,
    
    [Parameter()]
    [switch]$Test,
    
    [Parameter()]
    [switch]$Docker
)

$ErrorActionPreference = 'Stop'

Write-Host "🔨 Building PoliSync - C# .NET 8 Insurance Engine" -ForegroundColor Cyan
Write-Host "Configuration: $Configuration" -ForegroundColor Yellow

# Clean
if ($Clean) {
    Write-Host "`n🧹 Cleaning..." -ForegroundColor Yellow
    dotnet clean -c $Configuration
    Remove-Item -Path "src/*/bin" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "src/*/obj" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "tests/*/bin" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "tests/*/obj" -Recurse -Force -ErrorAction SilentlyContinue
}

# Restore
Write-Host "`n📦 Restoring packages..." -ForegroundColor Yellow
dotnet restore
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Restore failed" -ForegroundColor Red
    exit 1
}

# Build
Write-Host "`n🔨 Building solution..." -ForegroundColor Yellow
dotnet build -c $Configuration --no-restore
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
}

# Test
if ($Test) {
    Write-Host "`n🧪 Running tests..." -ForegroundColor Yellow
    dotnet test -c $Configuration --no-build --verbosity normal
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Tests failed" -ForegroundColor Red
        exit 1
    }
}

# Docker
if ($Docker) {
    Write-Host "`n🐳 Building Docker image..." -ForegroundColor Yellow
    docker build -t polisync:latest .
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Docker build failed" -ForegroundColor Red
        exit 1
    }
    Write-Host "✅ Docker image built: polisync:latest" -ForegroundColor Green
}

Write-Host "`n✅ Build completed successfully!" -ForegroundColor Green
Write-Host "`nTo run:" -ForegroundColor Cyan
Write-Host "  cd src/PoliSync.ApiHost" -ForegroundColor White
Write-Host "  dotnet run" -ForegroundColor White
Write-Host "`nOr with Docker:" -ForegroundColor Cyan
Write-Host "  docker-compose up" -ForegroundColor White
