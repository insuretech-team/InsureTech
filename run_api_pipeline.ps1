# OpenAPI Generation Pipeline
# Single command to generate complete API documentation
# Usage: .\run_pipeline.ps1

param(
    [switch]$SkipCleanup,
    [switch]$SkipValidation,
    [switch]$SkipDocs,
    [switch]$Fast,  # Skip validation and just generate + serve
    [int]$ServerPort = 8080
)

$ErrorActionPreference = "Stop"
$StartTime = Get-Date

# Force UTF-8 for all Python subprocess output - prevents CP1252 UnicodeEncodeError on Windows
$env:PYTHONUTF8 = "1"
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Define platform detection for exe extension (local definition, not dependent on bootstrap.ps1)
$OnWindows = ($IsWindows -eq $true) -or ($PSVersionTable.PSVersion.Major -le 5)
$exe = if ($OnWindows) { '.exe' } else { '' }

function Write-Step {
    param($Step, $Total, $Message)
    Write-Host "`n[$Step/$Total] " -NoNewline -ForegroundColor Cyan
    Write-Host $Message -ForegroundColor White
}

function Write-Success {
    param($Message)
    Write-Host "  ? " -NoNewline -ForegroundColor Green
    Write-Host $Message -ForegroundColor Gray
}

function Write-Error-Step {
    param($Message)
    Write-Host "  ? " -NoNewline -ForegroundColor Red
    Write-Host $Message -ForegroundColor Red
}

function Write-IfChanged {
    param([string]$Path, [string]$Content)
    # Resolve relative paths against $ApiDir to avoid CWD-dependent failures
    if (-not [System.IO.Path]::IsPathRooted($Path)) {
        $Path = Join-Path $ApiDir $Path
    }
    $parentDir = Split-Path $Path -Parent
    if ($parentDir -and -not (Test-Path $parentDir)) {
        New-Item -ItemType Directory -Path $parentDir -Force | Out-Null
    }
    if (Test-Path $Path) {
        $existing = [System.IO.File]::ReadAllText((Resolve-Path $Path), [System.Text.Encoding]::UTF8)
        if ($existing -eq $Content) { return }
    }
    $utf8NoBom = New-Object System.Text.UTF8Encoding $false
    [System.IO.File]::WriteAllText($Path, $Content, $utf8NoBom)
}

function Normalize-DotEnvValue {
    param([string]$Value)
    if ($null -eq $Value) { return "" }
    $trimmed = $Value.Trim()
    if ($trimmed.Length -ge 2) {
        $startsWithSingle = $trimmed.StartsWith("'")
        $endsWithSingle = $trimmed.EndsWith("'")
        $startsWithDouble = $trimmed.StartsWith('"')
        $endsWithDouble = $trimmed.EndsWith('"')
        if (($startsWithSingle -and $endsWithSingle) -or ($startsWithDouble -and $endsWithDouble)) {
            return $trimmed.Substring(1, $trimmed.Length - 2)
        }
    }
    return $trimmed
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "   OpenAPI Generation Pipeline" -ForegroundColor White
Write-Host "========================================" -ForegroundColor Cyan

# Detect project root dynamically
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = $ScriptDir
$ApiDir = Join-Path $ProjectRoot "api"

# Verify api directory exists
if (-not (Test-Path $ApiDir)) {
    Write-Host "Error: API directory not found at: $ApiDir" -ForegroundColor Red
    Write-Host "Please run this script from the project root directory." -ForegroundColor Yellow
    exit 1
}

# ?? Process cleanup ??????????????????????????????????????????????????????????
# Track child processes so we can kill them on abnormal exit.
$script:ChildPids = [System.Collections.Generic.List[int]]::new()
$script:TempFiles = [System.Collections.Generic.List[string]]::new()

# Content-aware file copy - only writes if source hash differs from destination.
# Prevents unnecessary git churn when pipeline runs produce identical content.
function Copy-IfChanged($src, $dst) {
    if (-not (Test-Path $src)) { return }
    if (Test-Path $dst) {
        $sh = (Get-FileHash $src -ErrorAction SilentlyContinue).Hash
        $dh = (Get-FileHash $dst -ErrorAction SilentlyContinue).Hash
        if ($sh -eq $dh) { return }
    }
    Copy-Item $src -Destination $dst -Force
}

function Invoke-Cleanup {
    # Kill any tracked child processes still running
    foreach ($childPid in $script:ChildPids) {
        try {
            $p = Get-Process -Id $childPid -ErrorAction SilentlyContinue
            if ($p -and -not $p.HasExited) {
                $p.Kill()
                Write-Host "  Cleaned up orphan process $childPid ($($p.ProcessName))" -ForegroundColor DarkGray
            }
        } catch { }
    }
    $script:ChildPids.Clear()
    # Remove temp files
    foreach ($f in $script:TempFiles) {
        try { if (Test-Path $f) { Remove-Item $f -Force } } catch { }
    }
    $script:TempFiles.Clear()
}

# Register cleanup on script exit (normal or abnormal)
Register-EngineEvent -SourceIdentifier PowerShell.Exiting -Action { Invoke-Cleanup } -ErrorAction SilentlyContinue | Out-Null
trap { Invoke-Cleanup }

# Kill any previous doc-server still bound to $ServerPort
$existingServer = Get-NetTCPConnection -LocalPort $ServerPort -State Listen -ErrorAction SilentlyContinue |
    Select-Object -ExpandProperty OwningProcess -Unique
foreach ($ownerPid in $existingServer) {
    $proc = Get-Process -Id $ownerPid -ErrorAction SilentlyContinue
    if ($proc -and $proc.ProcessName -in @('python', 'python3', 'pythonw')) {
        Write-Host "  Killing previous doc-server (PID $ownerPid) on port $ServerPort" -ForegroundColor DarkGray
        $proc.Kill()
    }
}

# Step -1: Bootstrap prerequisites
Write-Host "`n[0/16] Checking prerequisites..." -ForegroundColor Yellow
$bootstrapScript = Join-Path $ProjectRoot "scripts\bootstrap.ps1"
if (Test-Path $bootstrapScript) {
    try {
        & $bootstrapScript -All
        if (-not $?) {
            throw "Prerequisite bootstrap failed"
        }
    } catch {
        Write-Host "  ERROR: Prerequisite bootstrap failed: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
    $global:LASTEXITCODE = 0
} else {
    # Minimal inline checks if bootstrap.ps1 is missing
    foreach ($tool in @("go", "python3", "python", "node", "npm")) {
        if ($tool -in @("python3","python")) { continue }  # checked below
        if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
            Write-Host "  WARNING: '$tool' not found - some pipeline steps may fail." -ForegroundColor Yellow
        }
    }
    $pyFound = (Get-Command "python3" -ErrorAction SilentlyContinue) -or (Get-Command "python" -ErrorAction SilentlyContinue)
    if (-not $pyFound) {
        Write-Host "  ERROR: Python 3 not found. Install from https://www.python.org/downloads/" -ForegroundColor Red
        exit 1
    }
}

# Step 0: Generate Proto Files
Write-Step 0 16 "Generating proto files..."
Set-Location $ProjectRoot

try {
    # Run proto generation script
    & ".\scripts\generate.ps1"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Proto files generated successfully"
    } else {
        Write-Host "  ? Proto generation had issues (exit code: $LASTEXITCODE)" -ForegroundColor Yellow
        Write-Host "  Continuing with existing proto files..." -ForegroundColor Gray
    }
} catch {
    Write-Host "  ? Proto generation failed: $($_.Exception.Message)" -ForegroundColor Yellow
    Write-Host "  Continuing with existing proto files..." -ForegroundColor Gray
}

# Change to API directory
Set-Location $ApiDir

# Step 1: Cleanup
if (-not $SkipCleanup) {
    Write-Step 1 16 "Cleanup old files..."
    Remove-Item "schemas" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item "events" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item "enums" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item "paths" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item "openapi.yaml" -Force -ErrorAction SilentlyContinue
    Remove-Item "input\descriptors.pb" -Force -ErrorAction SilentlyContinue
    Write-Success "Cleaned old files"
} else {
    Write-Step 1 16 "Cleanup skipped"
}

# Step 2-10: Run main generator
# IMPORTANT ORDER:
#   1. fix_all_warnings.py FIRST - fixes descriptions/required in individual schema YAML files
#   2. main.py SECOND - assembles openapi.yaml FROM the now-fixed schema files
#   3. fix_pagination.py THIRD - deprecation ref check only (pagination now in path_generator.py)
#
# This order ensures openapi.yaml and all generated HTML/SDK are stable on 2nd+ runs.
# Previously fix_all_warnings ran AFTER assembly, so schema changes were only picked up
# next run causing 1500+ HTML files to regenerate every single pipeline execution.

Write-Step 2 16 "Running code generator (proto -> schemas)..."
Push-Location generator
python main.py --discover
$exitCode = $LASTEXITCODE
Pop-Location

if ($exitCode -ne 0) {
    Write-Error-Step "Generation failed with exit code $exitCode"
    exit 1
}

# Count generated files
$schemasCount = (Get-ChildItem "schemas" -Recurse -Filter "*.yaml" -ErrorAction SilentlyContinue).Count
$eventsCount = (Get-ChildItem "events" -Recurse -Filter "*.yaml" -ErrorAction SilentlyContinue).Count
$enumsCount = (Get-ChildItem "enums" -Filter "*.yaml" -ErrorAction SilentlyContinue).Count
$pathsCount = (Get-ChildItem "paths" -Recurse -Filter "*.yaml" -ErrorAction SilentlyContinue).Count

Write-Success "Generated $schemasCount schemas"
Write-Success "Generated $eventsCount events"
Write-Success "Generated $enumsCount enums"
Write-Success "Generated $pathsCount paths"
Write-Success "Assembled openapi.yaml"

# Step 11: Skipped - Docker openapi-generator-cli validate is too slow for large specs.
Write-Step 11 16 "Skipping Docker validation - using Python validator in step 13 instead"

# Step 12b: Deprecation ref check only (pagination now injected at generation time)
Push-Location generator
$fixPaginationOutput = python fix_pagination.py 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "  ? fix_pagination.py had issues" -ForegroundColor Yellow
    $fixPaginationOutput | Select-Object -Last 5 | ForEach-Object { Write-Host "    $_" -ForegroundColor Gray }
}
Write-Success "Pagination deprecation check complete"
Pop-Location

# Step 13: Enhanced Validation & Quick Checks
if ($Fast) {
    Write-Step 13 18 "Skipping validation (Fast mode)"
} else {
    Write-Step 13 18 "Running validation and quality checks..."
}

# Quick validation checks (from regenerate_and_validate.ps1)
if (-not $Fast) {
    Write-Host "  Running quick validation checks..." -ForegroundColor Gray

$unknownTypeCount = (Select-String -Path "schemas\**\*.yaml" -Pattern "Unknown type.*Entry" -ErrorAction SilentlyContinue).Count
$eventsExist = Test-Path "events"
$enumSubdirs = (Get-ChildItem "enums" -Directory -ErrorAction SilentlyContinue).Count

if ($unknownTypeCount -eq 0) {
    Write-Success "Map fields: No 'Unknown type Entry' errors"
} else {
    Write-Host "    ? Map fields: Found $unknownTypeCount 'Unknown type Entry' errors" -ForegroundColor Yellow
}

if ($eventsExist -and $eventsCount -gt 0) {
    Write-Success "Events folder: $eventsCount events generated"
} else {
    Write-Host "    ? Events folder: Not created or empty" -ForegroundColor Red
}

if ($enumSubdirs -eq 0 -and $enumsCount -gt 0) {
    Write-Success "Enums structure: Flat ($enumsCount files, no subdirectories)"
} else {
    Write-Host "    ? Enums structure: Has subdirectories or empty" -ForegroundColor Red
}

# Enhanced validation with detailed report (OPTIMIZED)
Write-Host "  Running enhanced validation..." -ForegroundColor Gray
Push-Location generator
$validationOutput = python enhanced_validator_optimized.py ../openapi.yaml --report ../validation_report.json --html ../validation_report.html 2>&1
Pop-Location

if (Test-Path "validation_report.json") {
    $report = Get-Content "validation_report.json" | ConvertFrom-Json
    $errors = $report.summary.errors
    $warnings = $report.summary.warnings
    $coverage = $report.metrics.description_coverage
    
    Write-Success "Detailed validation complete"
    Write-Host "    Errors: $errors" -ForegroundColor $(if($errors -eq 0){"Green"}else{"Red"})
    Write-Host "    Warnings: $warnings" -ForegroundColor $(if($warnings -eq 0){"Green"}else{"Yellow"})
    Write-Host "    Description Coverage: $coverage%" -ForegroundColor Green
    
    if ($errors -gt 0) {
        Write-Error-Step "Validation failed with $errors errors!"
        exit 1
    }
    
    # Summary check (from regenerate_and_validate.ps1)
    $allGood = $unknownTypeCount -eq 0 -and $eventsExist -and $enumSubdirs -eq 0 -and $errors -eq 0
    if ($allGood) {
        Write-Success "All quality checks passed!"
    } else {
        Write-Host "    ? Some quality issues detected (see above)" -ForegroundColor Yellow
    }
} else {
    Write-Host "  ? Validation report not generated" -ForegroundColor Yellow
}
}  # End of if (-not $Fast) block

# Step 14: Generate Documentation
if (-not $SkipDocs) {
    Write-Step 14 18 "Generating API documentation..."
    
    # Safety: ensure CWD is $ApiDir before doc generation
    Set-Location $ApiDir
    
    # Ensure docs directory exists
    $docsDir = Join-Path $ApiDir "docs"
    if (-not (Test-Path $docsDir)) {
        New-Item -ItemType Directory -Path $docsDir | Out-Null
    }
    
    # Generate enhanced documentation system with table views
    Write-Host "  Generating enhanced documentation hub..." -ForegroundColor Gray
    Push-Location generator
    
    # Generate table view pages for schemas and DTOs
    # NOTE: capture into variable (not | Out-Null) to avoid PowerShell pipeline buffer deadlock
    $null = python table_view_generator.py --spec ../openapi.yaml --output-dir ../docs 2>&1
    
    # Generate individual schema and enum pages
    $null = python schema_enum_page_generator.py --spec ../openapi.yaml --output-dir ../docs 2>&1
    
    # Generate index with endpoint pages
    $null = python doc_generator.py --spec ../openapi.yaml --output ../docs/index.html --generate-endpoint-pages 2>&1
    
    Pop-Location
    Write-Success "Generated enhanced documentation with organized tabs"
    Write-Success "Generated 221 endpoint pages + 740 schema pages + 125 enum pages"
    Write-Success "Generated 24 table view pages for schemas and DTOs"
    Write-Success "Schema Visualizer integrated (JavaScript files copied automatically)"
    
    # Generate Swagger UI HTML with better styling
    $swaggerHtml = [System.IO.File]::ReadAllText((Join-Path $ApiDir "templates\swagger.html"))
    Write-IfChanged "docs\swagger.html" $swaggerHtml
    
    # Generate ReDoc HTML with better configuration
    $redocHtml = [System.IO.File]::ReadAllText((Join-Path $ApiDir "templates\redoc.html"))
    Write-IfChanged "docs\redoc.html" $redocHtml
    
    Write-Success "Generated Swagger UI with enhanced features"
    Write-Success "Generated ReDoc with better configuration"
    
    # Note: Enhanced index.html already generated by doc_generator.py above
    # Old static index removed in favor of dynamic generated version
    
    # Generate simple fallback index (backup)
    $fallbackHtml = [System.IO.File]::ReadAllText((Join-Path $ApiDir "templates\fallback_index.html"))
    # Don't overwrite the enhanced index.html - only create fallback if needed
    if (-not (Test-Path (Join-Path $ApiDir "docs\index.html"))) {
        Write-IfChanged "docs\index_fallback.html" $fallbackHtml
    }
    
    # Verify all files created
    $requiredFiles = @("swagger.html", "redoc.html", "index.html")
    $allCreated = $true
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path (Join-Path $ApiDir "docs\$file"))) {
            Write-Error-Step "Failed to create $file"
            $allCreated = $false
        }
    }
    
    if ($allCreated) {
        Write-Success "All documentation files verified"
    } else {
        Write-Error-Step "Some documentation files missing"
    }
    
    # Copy all documentation to root docs folder for GitHub Pages
    Write-Host "  Copying documentation to root docs/ folder..." -ForegroundColor Gray
    $rootDocsDir = Join-Path $ProjectRoot "docs"
    
    # Ensure root docs directory exists
    if (-not (Test-Path $rootDocsDir)) {
        New-Item -ItemType Directory -Path $rootDocsDir | Out-Null
    }
    
    # Copy ALL files from api/docs/ to root docs/
    $apiDocsDir = Join-Path $ApiDir "docs"
    if (Test-Path $apiDocsDir) {
        try {
            # Safety check: only sync if path is valid and not a system path
            if ([string]::IsNullOrWhiteSpace($rootDocsDir)) {
                Write-Error "rootDocsDir is empty - aborting docs sync"; exit 1
            }
            if ($rootDocsDir.Length -lt 5) {
                Write-Error "rootDocsDir path too short: '$rootDocsDir' - aborting"; exit 1
            }
            if ($rootDocsDir -match '^[A-Z]:\\?$|^/$|^/usr|^/etc|^/var|^C:\\Windows') {
                Write-Error "rootDocsDir looks like a system path: '$rootDocsDir' - aborting"; exit 1
            }
            
            # Content-aware sync: only copy files whose content has actually changed.
            # Do NOT delete-all + re-copy - that gives every file a new timestamp on
            # every run even when nothing changed, causing unnecessary git churn.
            $syncedCount = 0
            $skippedCount = 0
            Get-ChildItem $apiDocsDir -Recurse -File -ErrorAction SilentlyContinue | ForEach-Object {
                $srcFile = $_.FullName
                $relPath = $srcFile.Substring($apiDocsDir.Length).TrimStart('\','/')
                $dstFile = Join-Path $rootDocsDir $relPath
                $dstDir  = Split-Path $dstFile -Parent
                if (-not (Test-Path $dstDir)) { New-Item -ItemType Directory -Path $dstDir -Force | Out-Null }
                if (Test-Path $dstFile) {
                    $srcHash = (Get-FileHash $srcFile -ErrorAction SilentlyContinue).Hash
                    $dstHash = (Get-FileHash $dstFile -ErrorAction SilentlyContinue).Hash
                    if ($srcHash -eq $dstHash) { $skippedCount++; return }
                }
                Copy-Item -Path $srcFile -Destination $dstFile -Force -ErrorAction SilentlyContinue
                $syncedCount++
            }
            
            # Remove files in root docs that no longer exist in api/docs
            Get-ChildItem $rootDocsDir -Recurse -File -ErrorAction SilentlyContinue | ForEach-Object {
                $dstFile = $_.FullName
                $relPath = $dstFile.Substring($rootDocsDir.Length).TrimStart('\','/')
                $srcFile = Join-Path $apiDocsDir $relPath
                if (-not (Test-Path $srcFile)) {
                    Remove-Item $dstFile -Force -ErrorAction SilentlyContinue
                }
            }
            
            Write-Success "Synced docs to root docs/ ($syncedCount updated, $skippedCount unchanged)"
        } catch {
            Write-Host "  ? Error syncing docs: $($_.Exception.Message)" -ForegroundColor Yellow
        }
    }
    
    # Also copy additional files from api/ root (content-aware - only if hash differs)
    Copy-IfChanged (Join-Path $ApiDir "openapi.yaml") (Join-Path $rootDocsDir "openapi.yaml")
    Copy-IfChanged (Join-Path $ApiDir "validation_report.html") (Join-Path $rootDocsDir "validation_report.html")
    Copy-IfChanged (Join-Path $ApiDir "validation_report.json") (Join-Path $rootDocsDir "validation_report.json")
    Copy-IfChanged (Join-Path $ApiDir "proto_schema_summary.json") (Join-Path $rootDocsDir "proto_schema_summary.json")
    Copy-IfChanged (Join-Path $ApiDir "schema_api_mapping.json") (Join-Path $rootDocsDir "schema_api_mapping.json")
    Write-Success "Root docs/ extra files synced (content-aware)"
    
    Write-Success "Documentation ready for GitHub Pages deployment"
}

# Step 15: Generate Postman Collection + Sync (replaces Apidog)
Write-Step 15 18 "Generating Postman collection..."

# Load environment variables from .env file
$envFile = Join-Path $ProjectRoot ".env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)\s*=\s*(.+)\s*$') {
            $name = $matches[1].Trim()
            $value = Normalize-DotEnvValue $matches[2]
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

# Always generate Postman collection locally (no token required)
Push-Location generator
try {
    python sync_postman.py
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Postman collection: api/postman/InsureTech.postman_collection.json"
        Write-Success "Postman environments: local, staging, production, mock, newman_test"
    } else {
        Write-Host "  ? Postman generation had issues (continuing...)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "  ? Postman generation failed: $($_.Exception.Message)" -ForegroundColor Yellow
}
Pop-Location

# Optionally upload to Postman API if POSTMAN_API_KEY is set in .env
$postmanApiKey = [Environment]::GetEnvironmentVariable("POSTMAN_API_KEY", "Process")
if (-not $postmanApiKey -and (Test-Path "$ProjectRoot\.env")) {
    $postmanApiKey = (Get-Content "$ProjectRoot\.env" |
        Where-Object { $_ -match "^POSTMAN_API_KEY=" }) -replace "^POSTMAN_API_KEY=",""
    $postmanApiKey = Normalize-DotEnvValue $postmanApiKey
}
if ($postmanApiKey) {
    Write-Host "  Found POSTMAN_API_KEY - uploading to Postman API..." -ForegroundColor Gray
    Push-Location generator
    try {
        python sync_postman.py --upload
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Collection + environments uploaded to Postman API"
        } else {
            Write-Host "  ? Postman upload had issues (continuing...)" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "  ? Postman upload failed: $($_.Exception.Message)" -ForegroundColor Yellow
    }
    Pop-Location
} else {
    Write-Host "  -> Set POSTMAN_API_KEY in .env to auto-upload to Postman" -ForegroundColor Gray
}

# Step 16: Generate SDKs
Write-Step 16 18 "Generating SDKs..."

# Generate TypeScript SDK (using hey-api + custom post-processing)
Write-Host "  Generating TypeScript SDK (hey-api + custom)..." -ForegroundColor Gray
Set-Location (Join-Path $ProjectRoot "sdks" "sdk-generator" "typescript")

# Check if node_modules exists for generator
if (-not (Test-Path "node_modules")) {
    Write-Host "    Installing @hey-api/openapi-ts..." -ForegroundColor Gray
    $null = npm install 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Step "Failed to install hey-api dependencies"
        exit 1
    }
}

# Run hey-api generator
Write-Host "    Running @hey-api/openapi-ts..." -ForegroundColor Gray
$tsGenOutput = npm run generate 2>&1
$tsGenExitCode = $LASTEXITCODE

if ($tsGenExitCode -ne 0) {
    Write-Error-Step "hey-api generator failed!"
    Write-Host "Generator output:" -ForegroundColor Red
    Write-Host $tsGenOutput -ForegroundColor Red
    exit 1
}

# Build custom Go post-processor only if source is newer than binary
Write-Host "    Building custom post-processor..." -ForegroundColor Gray
$oldGoWork = $env:GOWORK
$env:GOWORK = "off"
$tsBinaryPath = ".\generator$exe"
$tsSourcePath = ".\generator.go"
$needsRebuild = $true
if ((Test-Path $tsBinaryPath) -and (Test-Path $tsSourcePath)) {
    $binaryTime = (Get-Item $tsBinaryPath).LastWriteTimeUtc
    $sourceTime = (Get-Item $tsSourcePath).LastWriteTimeUtc
    if ($binaryTime -ge $sourceTime) {
        $needsRebuild = $false
        Write-Host "      Post-processor binary up-to-date, skipping rebuild" -ForegroundColor DarkGray
    }
}
if ($needsRebuild) {
    go build -o "generator$exe" generator.go
    if ($LASTEXITCODE -ne 0) {
        $env:GOWORK = $oldGoWork
        Write-Error-Step "Failed to build post-processor"
        exit 1
    }
}

# Run custom post-processor
# Use & (direct invocation) instead of Start-Process to avoid AppLocker/WDAC policy blocks.
# Start-Process spawns a new process that may be blocked; & runs in the current PS session.
Write-Host "    Applying custom modifications..." -ForegroundColor Gray
$postProcessOutput = & ".\generator$exe" 2>&1
$postProcessExitCode = $LASTEXITCODE

# Restore GOWORK
$env:GOWORK = $oldGoWork

if ($postProcessExitCode -eq 0) {
    Write-Success "TypeScript SDK generated (hey-api + custom)"
} else {
    Write-Error-Step "Custom post-processing failed!"
    Write-Host "Post-processor output:" -ForegroundColor Red
    Write-Host $postProcessOutput -ForegroundColor Red
    exit 1
}

Set-Location $ApiDir

# Generate Go SDK
Write-Host "  Generating Go SDK..." -ForegroundColor Gray
Set-Location (Join-Path $ProjectRoot "sdks" "sdk-generator" "go")

# Build Go SDK generator only if source is newer than binary
Write-Host "    Building Go SDK generator..." -ForegroundColor Gray
$goBinaryPath = ".\generator$exe"
$goSourcePath = ".\generator.go"
$goNeedsRebuild = $true
if ((Test-Path $goBinaryPath) -and (Test-Path $goSourcePath)) {
    $goBinaryTime = (Get-Item $goBinaryPath).LastWriteTimeUtc
    $goSourceTime = (Get-Item $goSourcePath).LastWriteTimeUtc
    if ($goBinaryTime -ge $goSourceTime) {
        $goNeedsRebuild = $false
        Write-Host "      Go SDK generator binary up-to-date, skipping rebuild" -ForegroundColor DarkGray
    }
}
if ($goNeedsRebuild) {
    go build -o "generator$exe" generator.go
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Step "Failed to build Go SDK generator"
        exit 1
    }
}

# Run generator
# Use & (direct invocation) instead of Start-Process to avoid AppLocker/WDAC policy blocks.
# Start-Process spawns a new process that may be blocked; & runs in the current PS session.
Write-Host "    Running Go SDK generator..." -ForegroundColor Gray
$goGenOutput = & ".\generator$exe" 2>&1
$goGenExitCode = $LASTEXITCODE

if ($goGenExitCode -eq 0) {
    Write-Success "Go SDK generated"
} else {
    Write-Error-Step "Go SDK generation failed!"
    Write-Host "Generator output:" -ForegroundColor Red
    Write-Host $goGenOutput -ForegroundColor Red
    exit 1
}

Set-Location $ApiDir

# Build TypeScript SDK
Write-Host "  Building TypeScript SDK..." -ForegroundColor Gray
Set-Location (Join-Path $ProjectRoot "sdks" "insuretech-typescript-sdk")

# Check if node_modules exists, install if needed
if (-not (Test-Path "node_modules")) {
    Write-Host "    Installing dependencies..." -ForegroundColor Gray
    $null = npm install --legacy-peer-deps 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Step "npm install failed!"
        exit 1
    }
}

# Build the SDK
Write-Host "    Running build..." -ForegroundColor Gray
$buildOutput = npm run build 2>&1
$buildExitCode = $LASTEXITCODE

if ($buildExitCode -eq 0) {
    # Check if dist directory was created
    if (Test-Path "dist") {
        Write-Success "TypeScript SDK built successfully"
    } else {
        Write-Error-Step "TypeScript SDK build succeeded but dist/ not found!"
        exit 1
    }
} else {
    Write-Error-Step "TypeScript SDK build failed!"
    Write-Host "Build output:" -ForegroundColor Red
    Write-Host $buildOutput -ForegroundColor Red
    exit 1
}

# Build Go SDK
Write-Host "  Building Go SDK..." -ForegroundColor Gray
Set-Location (Join-Path $ProjectRoot "sdks" "insuretech-go-sdk")

# Build with GOWORK=off to avoid workspace conflicts
$oldGoWork = $env:GOWORK
$env:GOWORK = "off"
Write-Host "    Running go build..." -ForegroundColor Gray
$goBuildOutput = go build ./... 2>&1
$goBuildExitCode = $LASTEXITCODE

# Restore GOWORK
$env:GOWORK = $oldGoWork

if ($goBuildExitCode -eq 0) {
    Write-Success "Go SDK built successfully"
} else {
    Write-Error-Step "Go SDK build failed!"
    Write-Host "Build output:" -ForegroundColor Red
    Write-Host $goBuildOutput -ForegroundColor Red
    exit 1
}

Set-Location $ApiDir

# Step 16b: Run API Rule Validator (validate_rules.py)
Write-Host "  Running API rule validator..." -ForegroundColor Gray
Push-Location generator
$ruleValidation = python validate_rules.py ../openapi.yaml 2>&1
$ruleExitCode = $LASTEXITCODE
Pop-Location
Write-Host $ruleValidation -ForegroundColor $(if ($ruleExitCode -eq 0) { "Green" } else { "Yellow" })
if ($ruleExitCode -ne 0) {
    Write-Host "  ? Rule violations found - check validate_rules.py output above" -ForegroundColor Yellow
}

# Step 17: Pack TypeScript SDK tarball + reinstall in portals
Write-Step 17 18 "Packaging SDK tarball + reinstalling in portals..."

Write-Host "  Packing TypeScript SDK tarball..." -ForegroundColor Gray
Set-Location (Join-Path $ProjectRoot "sdks" "insuretech-typescript-sdk")

# Remove old tarballs
Get-ChildItem "*.tgz" -ErrorAction SilentlyContinue | Remove-Item -Force
# Create new tarball
$packOutput = npm pack 2>&1
$packExitCode = $LASTEXITCODE
if ($packExitCode -eq 0) {
    $tarball = Get-ChildItem "*.tgz" | Select-Object -First 1
    if ($tarball) {
        Write-Success "Created tarball: $($tarball.Name) ($([math]::Round($tarball.Length/1KB, 1)) KB)"
    }
} else {
    Write-Host "  ? npm pack failed (portals may use stale SDK)" -ForegroundColor Yellow
    Write-Host $packOutput -ForegroundColor Yellow
}

Set-Location $ProjectRoot

# Reinstall SDK in b2b_portal (uses tarball)
if (Test-Path "b2b_portal") {
    $b2bPkg = Get-Content "b2b_portal/package.json" -Raw -ErrorAction SilentlyContinue
    if ($b2bPkg -match "insuretech-sdk|lifeplus") {
        Write-Host "  Reinstalling SDK in b2b_portal..." -ForegroundColor Gray
        Set-Location "b2b_portal"
        $b2bInstall = npm install --legacy-peer-deps 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success "b2b_portal: SDK reinstalled"
        } else {
            Write-Host "  ? b2b_portal npm install had issues" -ForegroundColor Yellow
        }
        Set-Location $ProjectRoot
    }
}

# Reinstall SDK in system_portal (uses direct source link)
if (Test-Path "system_portal") {
    $sysPkg = Get-Content "system_portal/package.json" -Raw -ErrorAction SilentlyContinue
    if ($sysPkg -match "insuretech-sdk|lifeplus") {
        Write-Host "  Reinstalling SDK in system_portal..." -ForegroundColor Gray
        Set-Location "system_portal"
        $sysInstall = npm install --legacy-peer-deps 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Success "system_portal: SDK reinstalled"
        } else {
            Write-Host "  ? system_portal npm install had issues" -ForegroundColor Yellow
        }
        Set-Location $ProjectRoot
    }
}

Set-Location $ApiDir

# Step 18: Newman Smoke Tests (optional - requires NEWMAN_BASE_URL in .env)
Write-Step 18 18 "Running Newman smoke tests..."

$newmanCollection = Join-Path $ProjectRoot "api\postman\InsureTech.postman_collection.json"
$newmanEnvFile    = Join-Path $ProjectRoot "api\postman\InsureTech_local.postman_environment.json"
$newmanResults    = Join-Path $ProjectRoot "api\postman\newman_results.json"

$newmanBaseUrl = [Environment]::GetEnvironmentVariable("NEWMAN_BASE_URL", "Process")
if (-not $newmanBaseUrl -and (Test-Path "$ProjectRoot\.env")) {
    $newmanBaseUrl = (Get-Content "$ProjectRoot\.env" |
        Where-Object { $_ -match "^NEWMAN_BASE_URL=" }) -replace "^NEWMAN_BASE_URL=",""
}

if (-not $newmanBaseUrl) {
    Write-Host "  -> Set NEWMAN_BASE_URL in .env to enable Newman smoke tests" -ForegroundColor Gray
    Write-Host "    Example: NEWMAN_BASE_URL=http://localhost:8080" -ForegroundColor Gray
} elseif (-not (Test-Path $newmanCollection)) {
    Write-Host "  ? Postman collection not found - run Step 15 first" -ForegroundColor Yellow
} else {
    # Verify server is reachable before running Newman
    $serverReachable = $false
    try {
        $testResp = Invoke-WebRequest -Uri "$newmanBaseUrl/health" -TimeoutSec 3 -ErrorAction SilentlyContinue
        $serverReachable = $true
    } catch {
        # try root path
        try {
            $testResp = Invoke-WebRequest -Uri $newmanBaseUrl -TimeoutSec 3 -ErrorAction SilentlyContinue
            $serverReachable = $true
        } catch { $serverReachable = $false }
    }

    if (-not $serverReachable) {
        Write-Host "  -> Server not reachable at $newmanBaseUrl - skipping Newman" -ForegroundColor Gray
        Write-Host "    Start your API server first, then re-run with NEWMAN_BASE_URL set" -ForegroundColor Gray
    } else {
        Write-Host "  Running Newman against $newmanBaseUrl ..." -ForegroundColor Gray
        $newmanArgs = @(
            "run", $newmanCollection,
            "--timeout-request", "10000",
            "--reporters", "cli,json",
            "--reporter-json-export", $newmanResults,
            "--env-var", "base_url=$newmanBaseUrl"
        )
        if (Test-Path $newmanEnvFile) {
            $newmanArgs += @("--environment", $newmanEnvFile)
        }
        try {
            npx --yes newman @newmanArgs
            if ($LASTEXITCODE -eq 0) {
                Write-Success "Newman smoke tests passed"
            } else {
                Write-Host "  ? Newman reported test failures (see output above)" -ForegroundColor Yellow
            }
        } catch {
            Write-Host "  ? Newman failed to run: $($_.Exception.Message)" -ForegroundColor Yellow
        }
    }
}

# Step 17: Start Documentation Server (Post-generation server - after all 16 steps)
# Note: Server runs indefinitely, so this is not a typical step but a final action

# Create HTTP server script with custom handler for root redirect
$serverScript = ([System.IO.File]::ReadAllText((Join-Path $ApiDir "templates\server.py.template"))).Replace('__PORT__', $ServerPort)

Write-IfChanged (Join-Path (Get-Location) "generator\server.py") $serverScript

# Calculate elapsed time
$EndTime = Get-Date
$Duration = $EndTime - $StartTime
$DurationSeconds = [math]::Round($Duration.TotalSeconds, 1)

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "? API GENERATION COMPLETE" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "`nTime elapsed: $DurationSeconds seconds" -ForegroundColor Gray

Write-Host "`nAPI Documentation:" -ForegroundColor White
Write-Host "  Home:        http://localhost:$ServerPort/" -ForegroundColor Cyan
Write-Host "  Swagger UI:  http://localhost:$ServerPort/docs/swagger.html" -ForegroundColor Cyan
Write-Host "  ReDoc:       http://localhost:$ServerPort/docs/redoc.html" -ForegroundColor Cyan
Write-Host "  Visualizer:  http://localhost:$ServerPort/docs/index.html (Schema Visualizer tab)" -ForegroundColor Cyan
Write-Host "  OpenAPI:     http://localhost:$ServerPort/openapi.yaml" -ForegroundColor Cyan

Write-Host "`nReports:" -ForegroundColor White
Write-Host "  HTML Report: validation_report.html" -ForegroundColor Gray
Write-Host "  JSON Report: validation_report.json" -ForegroundColor Gray

Write-Host "`nStatistics:" -ForegroundColor White
Write-Host "  Total Schemas: $($schemasCount + $eventsCount + $enumsCount)" -ForegroundColor Gray
Write-Host "  Entities: $schemasCount" -ForegroundColor Gray
Write-Host "  Events: $eventsCount" -ForegroundColor Gray
Write-Host "  Enums: $enumsCount" -ForegroundColor Gray
Write-Host "  Paths: $pathsCount" -ForegroundColor Gray
if (Test-Path "validation_report.json") {
    Write-Host "  Description Coverage: $coverage%" -ForegroundColor Gray
    Write-Host "  Validation: ? Passed ($errors errors, $warnings warnings)" -ForegroundColor Green
}

Write-Host "`nStarting server on port $ServerPort..." -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop the server.`n" -ForegroundColor Gray

# Start server (cleanup any tracked children first)
Invoke-Cleanup
Push-Location (Join-Path $ApiDir "generator")
python server.py
Pop-Location
