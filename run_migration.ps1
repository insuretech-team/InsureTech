# Production Migration Runner for InsureTech
param(
    [ValidateSet("primary", "backup", "both")]
    [string]$Target = "primary",
    [switch]$DryRun = $false,
    [switch]$Prune = $false,
    [switch]$Strict = $false
)

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "InsureTech Database Migration Runner" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# Step 1: Find project root
Write-Host "`n[1/6] Finding project root..." -ForegroundColor Yellow
$projectRoot = $PSScriptRoot
if (-not (Test-Path "$projectRoot\go.mod")) {
    Write-Host "  ERROR: go.mod not found" -ForegroundColor Red
    exit 1
}
Write-Host "  OK: Project root: $projectRoot" -ForegroundColor Green

# Step 2: Load .env
Write-Host "`n[2/6] Loading environment variables..." -ForegroundColor Yellow
$envFile = Join-Path $projectRoot ".env"
if (Test-Path $envFile) {
    Write-Host "  OK: Found .env" -ForegroundColor Green
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            # Remove surrounding quotes (both single and double)
            $value = $value -replace "^'|'$", ''
            $value = $value -replace '^"|"$', ''
            [System.Environment]::SetEnvironmentVariable($key, $value, "Process")
            Write-Host "    Set: $key" -ForegroundColor Gray
        }
    }
}
else {
    Write-Host "  WARNING: .env not found" -ForegroundColor Yellow
}

# Verify env vars
$requiredVars = @("PGHOST", "PGDATABASE", "PGUSER", "PGPASSWORD")
$missing = $requiredVars | Where-Object { -not [System.Environment]::GetEnvironmentVariable($_, "Process") }
if ($missing) {
    Write-Host "  ERROR: Missing env vars: $($missing -join ', ')" -ForegroundColor Red
    exit 1
}
Write-Host "  OK: All environment variables loaded" -ForegroundColor Green

# Step 3: Verify config
Write-Host "`n[3/6] Verifying configuration..." -ForegroundColor Yellow
$configPath = Join-Path $projectRoot "backend\inscore\configs\database.yaml"
if (-not (Test-Path $configPath)) {
    Write-Host "  ERROR: database.yaml not found" -ForegroundColor Red
    exit 1
}
Write-Host "  OK: Config file found" -ForegroundColor Green

# Step 4: Checking SSL certificates...
Write-Host "`n[4/7] Checking SSL certificates..." -ForegroundColor Yellow
$certsDir = Join-Path $projectRoot "backend\inscore\db\certs"
if (Test-Path $certsDir) {
    Write-Host "  OK: Certs directory exists" -ForegroundColor Green
}
else {
    Write-Host "  INFO: No certs directory (OK if not using SSL)" -ForegroundColor Gray
}

# Step 5: Regenerate proto code with GORM tags
Write-Host "`n[5/7] Regenerating proto code with GORM tags..." -ForegroundColor Yellow
$generateScript = Join-Path $projectRoot "scripts\generate.ps1"
if (Test-Path $generateScript) {
    & $generateScript
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  WARNING: Proto generation had issues (continuing...)" -ForegroundColor Yellow
    }
    else {
        Write-Host "  OK: Proto code regenerated" -ForegroundColor Green
    }
}
else {
    Write-Host "  SKIP: generate.ps1 not found" -ForegroundColor Gray
}

# Step 6: Skip building dbmanager - we will use "go run" directly
Write-Host "`n[6/7] Preparing dbmanager..." -ForegroundColor Yellow
$dbmanagerDir = Join-Path $projectRoot "backend\inscore\cmd\dbmanager"
Write-Host "  OK: Will use go run directly in $dbmanagerDir" -ForegroundColor Green

# Step 7: Run migration
Write-Host "`n[7/7] Running migration on $Target..." -ForegroundColor Yellow
if ($Prune) {
    Write-Host "  Mode: PRUNE (will remove zombie columns)" -ForegroundColor Yellow
}
if ($Strict) {
    Write-Host "  Mode: STRICT (will fail on schema drift)" -ForegroundColor Yellow
}
Write-Host "==========================================`n" -ForegroundColor Cyan

if ($DryRun) {
    Write-Host "DRY RUN - Would execute:" -ForegroundColor Yellow
    Write-Host "  cd $dbmanagerDir" -ForegroundColor White
    $cmdArgs = "migrate --target $Target"
    if ($Prune) { $cmdArgs += " --prune" }
    if ($Strict) { $cmdArgs += " --strict" }
    Write-Host "  go run . $cmdArgs" -ForegroundColor White
    exit 0
}

Push-Location $dbmanagerDir
try {
    $cmdArgs = @("migrate", "--target", $Target)
    if ($Prune) { $cmdArgs += "--prune" }
    if ($Strict) { $cmdArgs += "--strict" }
    
    Write-Host "Executing: go run . $($cmdArgs -join ' ')`n" -ForegroundColor Cyan
    & go run . $cmdArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n==========================================" -ForegroundColor Cyan
        Write-Host "SUCCESS: Migration completed!" -ForegroundColor Green
        Write-Host "==========================================" -ForegroundColor Cyan
    }
    else {
        Write-Host "`n==========================================" -ForegroundColor Red
        Write-Host "FAILED: Exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host "==========================================" -ForegroundColor Red
        exit $LASTEXITCODE
    }
}
finally {
    Pop-Location
}
