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

# Step 0: Bootstrap prerequisites
Write-Host "`n[0/7] Checking prerequisites..." -ForegroundColor Yellow
$bootstrapScript = Join-Path $PSScriptRoot "scripts" "bootstrap.ps1"
if (Test-Path $bootstrapScript) {
    & $bootstrapScript -GoOnly
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
} else {
    if (-not (Get-Command "go" -ErrorAction SilentlyContinue)) {
        Write-Host "  ERROR: 'go' not found and bootstrap.ps1 is missing." -ForegroundColor Red
        Write-Host "  Install Go 1.25 from https://go.dev/dl/" -ForegroundColor Yellow
        exit 1
    }
    Write-Host "  OK: go found" -ForegroundColor Green
}

# Step 1: Find project root
Write-Host "`n[1/7] Finding project root..." -ForegroundColor Yellow
$projectRoot = $PSScriptRoot
if (-not (Test-Path (Join-Path $projectRoot "go.mod"))) {
    Write-Host "  ERROR: go.mod not found" -ForegroundColor Red
    exit 1
}
Write-Host "  OK: Project root: $projectRoot" -ForegroundColor Green

# Step 2: Load .env
Write-Host "`n[2/7] Loading environment variables..." -ForegroundColor Yellow
$envFile = Join-Path $projectRoot ".env"
if (Test-Path $envFile) {
    Write-Host "  OK: Found .env" -ForegroundColor Green
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([A-Za-z_][A-Za-z0-9_]*)=(.*)$') {
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

# Check if GOWORK is set to off (will cause issues)
$gowork = [System.Environment]::GetEnvironmentVariable("GOWORK", "Process")
if ($gowork -eq "off") {
    Write-Host "  WARNING: GOWORK=off is set in environment. This may cause Go module resolution issues." -ForegroundColor Yellow
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
Write-Host "`n[3/7] Verifying configuration..." -ForegroundColor Yellow
$configPath = Join-Path $projectRoot "backend" "inscore" "configs" "database.yaml"
if (-not (Test-Path $configPath)) {
    Write-Host "  ERROR: database.yaml not found" -ForegroundColor Red
    exit 1
}
Write-Host "  OK: Config file found" -ForegroundColor Green

# Step 4: Checking SSL certificates...
Write-Host "`n[4/7] Checking SSL certificates..." -ForegroundColor Yellow
$certsDir = Join-Path $projectRoot "backend" "inscore" "db" "certs"
if (Test-Path $certsDir) {
    Write-Host "  OK: Certs directory exists" -ForegroundColor Green
}
else {
    Write-Host "  INFO: No certs directory (OK if not using SSL)" -ForegroundColor Gray
}

# Step 5: Regenerate proto code with GORM tags
Write-Host "`n[5/7] Regenerating proto code with GORM tags..." -ForegroundColor Yellow
$generateScript = Join-Path $projectRoot "scripts" "generate.ps1"
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

# Step 6: Skip building dbx - we will use "go run" directly
Write-Host "`n[6/7] Preparing dbx..." -ForegroundColor Yellow
$dbxDir = Join-Path $projectRoot "backend" "inscore" "cmd" "dbx"
Write-Host "  OK: Will use go run directly from project root" -ForegroundColor Green

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
    Write-Host "  cd $projectRoot" -ForegroundColor White
    $cmdArgs = "migrate --target $Target"
    if ($Prune) { $cmdArgs += " --prune" }
    if ($Strict) { $cmdArgs += " --strict" }
    Write-Host "  go run ./backend/inscore/cmd/dbx $cmdArgs" -ForegroundColor White
    exit 0
}

try {
    Push-Location $projectRoot
    $cmdArgs = @("migrate", "--target", $Target)
    if ($Prune) { $cmdArgs += "--prune" }
    if ($Strict) { $cmdArgs += "--strict" }
    
    Write-Host "Executing: go run ./backend/inscore/cmd/dbx $($cmdArgs -join ' ')`n" -ForegroundColor Cyan
    & go run ./backend/inscore/cmd/dbx $cmdArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n===========================================" -ForegroundColor Cyan
        Write-Host "SUCCESS: Migration completed!" -ForegroundColor Green
        Write-Host "===========================================" -ForegroundColor Cyan
    }
    else {
        Write-Host "`n===========================================" -ForegroundColor Red
        Write-Host "FAILED: Exit code $LASTEXITCODE" -ForegroundColor Red
        Write-Host "===========================================" -ForegroundColor Red
        exit $LASTEXITCODE
    }
}
finally {
    Pop-Location
}
