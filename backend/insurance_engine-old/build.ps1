<#
.SYNOPSIS
    Build script for InsuranceEngine.
.DESCRIPTION
    Handles cleaning, restoring, building, testing, and publishing the InsuranceEngine solution.
#>

param (
    [string]$Configuration = "Release",
    [switch]$SkipTests,
    [switch]$Docker
)

$ErrorActionPreference = "Stop"
$ProjectRoot = Get-Location
$SolutionFile = Join-Path $ProjectRoot "InsuranceEngine.sln"
$ApiHostProject = Join-Path $ProjectRoot "src\InsuranceEngine.ApiHost\InsuranceEngine.ApiHost.csproj"

Write-Host "--- InsuranceEngine Build System ---" -ForegroundColor Cyan

# 1. Clean
Write-Host "[1/5] Cleaning solution..." -ForegroundColor Gray
dotnet clean $SolutionFile -c $Configuration

# 2. Restore
Write-Host "[2/5] Restoring dependencies..." -ForegroundColor Gray
dotnet restore $SolutionFile

# 3. Build
Write-Host "[3/5] Building solution ($Configuration)..." -ForegroundColor Gray
dotnet build $SolutionFile -c $Configuration --no-restore

# 4. Test
if (-not $SkipTests) {
    Write-Host "[4/5] Running tests..." -ForegroundColor Gray
    dotnet test $SolutionFile -c $Configuration --no-build
} else {
    Write-Host "[4/5] Skipping tests." -ForegroundColor Yellow
}

# 5. Publish (Optional)
if ($Docker) {
    Write-Host "[5/5] Building Docker image..." -ForegroundColor Gray
    docker build -t insurance-engine:latest .
} else {
    Write-Host "[5/5] Skipping Docker build." -ForegroundColor Yellow
}

Write-Host "--- Build Complete ---" -ForegroundColor Green
