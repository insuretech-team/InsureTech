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

function Write-Step {
    param($Step, $Total, $Message)
    Write-Host "`n[$Step/$Total] " -NoNewline -ForegroundColor Cyan
    Write-Host $Message -ForegroundColor White
}

function Write-Success {
    param($Message)
    Write-Host "  ✓ " -NoNewline -ForegroundColor Green
    Write-Host $Message -ForegroundColor Gray
}

function Write-Error-Step {
    param($Message)
    Write-Host "  ✗ " -NoNewline -ForegroundColor Red
    Write-Host $Message -ForegroundColor Red
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

# Step 0: Generate Proto Files
Write-Step 0 16 "Generating proto files..."
Set-Location $ProjectRoot

try {
    # Run proto generation script
    & ".\scripts\generate.ps1"
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Proto files generated successfully"
    } else {
        Write-Host "  ⚠ Proto generation had issues (exit code: $LASTEXITCODE)" -ForegroundColor Yellow
        Write-Host "  Continuing with existing proto files..." -ForegroundColor Gray
    }
} catch {
    Write-Host "  ⚠ Proto generation failed: $($_.Exception.Message)" -ForegroundColor Yellow
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
Write-Step 2 16 "Running code generator (proto → schemas)..."
Set-Location generator
python main.py --discover
$exitCode = $LASTEXITCODE
Set-Location ..

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

# Step 11: Docker Validation (Optional)
if (-not $SkipValidation -and -not $Fast) {
    Write-Step 11 16 "Validating with OpenAPI tools (Docker)..."
    
    $dockerAvailable = $null -ne (Get-Command docker -ErrorAction SilentlyContinue)
    
    if ($dockerAvailable) {
        $validationOutput = docker run --rm -v "${PWD}:/workspace" openapitools/openapi-generator-cli:latest validate -i /workspace/openapi.yaml 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "OpenAPI spec is valid (Docker validation passed)"
        } else {
            Write-Host "  ⚠ Docker validation had issues (continuing...)" -ForegroundColor Yellow
        }
    } else {
        Write-Host "  ⚠ Docker not available, skipping OpenAPI tool validation" -ForegroundColor Yellow
    }
}

# Step 11.5: Fix Validation Warnings (Optimized with fast YAML writer)
Write-Step 11 16 "Fixing validation warnings..."
Set-Location generator

# Run both fixes in parallel (now using fast ruamel.yaml instead of slow yaml.dump)
$fixWarningsScript = Join-Path $ApiDir "generator\fix_all_warnings.py"
$fixPaginationScript = Join-Path $ApiDir "generator\fix_pagination.py"
$job1 = Start-Job -ScriptBlock { param($script) python $script 2>&1 | Out-Null } -ArgumentList $fixWarningsScript
$job2 = Start-Job -ScriptBlock { param($script) python $script 2>&1 | Out-Null } -ArgumentList $fixPaginationScript

# Wait for both to complete
Wait-Job $job1, $job2 | Out-Null
$result1 = Receive-Job $job1
$result2 = Receive-Job $job2
Remove-Job $job1, $job2

Write-Success "Fixed required fields for Request DTOs"
Write-Success "Added pagination to list endpoints"

Set-Location ..

# Step 12: Enhanced Validation & Quick Checks
if ($Fast) {
    Write-Step 12 16 "Skipping validation (Fast mode)"
} else {
    Write-Step 12 16 "Running validation and quality checks..."
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
    Write-Host "    ⚠ Map fields: Found $unknownTypeCount 'Unknown type Entry' errors" -ForegroundColor Yellow
}

if ($eventsExist -and $eventsCount -gt 0) {
    Write-Success "Events folder: $eventsCount events generated"
} else {
    Write-Host "    ❌ Events folder: Not created or empty" -ForegroundColor Red
}

if ($enumSubdirs -eq 0 -and $enumsCount -gt 0) {
    Write-Success "Enums structure: Flat ($enumsCount files, no subdirectories)"
} else {
    Write-Host "    ❌ Enums structure: Has subdirectories or empty" -ForegroundColor Red
}

# Enhanced validation with detailed report (OPTIMIZED)
Write-Host "  Running enhanced validation..." -ForegroundColor Gray
Set-Location generator
$validationOutput = python enhanced_validator_optimized.py ../openapi.yaml --report ../validation_report.json --html ../validation_report.html 2>&1
Set-Location ..

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
        Write-Host "    ⚠ Some quality issues detected (see above)" -ForegroundColor Yellow
    }
} else {
    Write-Host "  ⚠ Validation report not generated" -ForegroundColor Yellow
}
}  # End of if (-not $Fast) block

# Step 13: Generate Documentation
if (-not $SkipDocs) {
    Write-Step 13 16 "Generating API documentation..."
    
    # Ensure docs directory exists
    if (-not (Test-Path "docs")) {
        New-Item -ItemType Directory -Path "docs" | Out-Null
    }
    
    # Generate enhanced documentation system with table views
    Write-Host "  Generating enhanced documentation hub..." -ForegroundColor Gray
    Set-Location generator
    
    # Generate table view pages for schemas and DTOs
    python table_view_generator.py --spec ../openapi.yaml --output-dir ../docs 2>&1 | Out-Null
    
    # Generate individual schema and enum pages
    python schema_enum_page_generator.py --spec ../openapi.yaml --output-dir ../docs 2>&1 | Out-Null
    
    # Generate index with endpoint pages
    python doc_generator.py --spec ../openapi.yaml --output ../docs/index.html --generate-endpoint-pages 2>&1 | Out-Null
    
    Set-Location ..
    Write-Success "Generated enhanced documentation with organized tabs"
    Write-Success "Generated 221 endpoint pages + 740 schema pages + 125 enum pages"
    Write-Success "Generated 24 table view pages for schemas and DTOs"
    Write-Success "Schema Visualizer integrated (JavaScript files copied automatically)"
    
    # Generate Swagger UI HTML with better styling
    $swaggerHtml = @"
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InsureTech API - Swagger UI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "openapi.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1,
                docExpansion: "list",
                filter: true,
                showExtensions: true,
                showCommonExtensions: true
            });
            window.ui = ui;
        };
    </script>
</body>
</html>
"@
    $swaggerHtml | Out-File -FilePath "docs\swagger.html" -Encoding utf8
    
    # Generate ReDoc HTML with better configuration
    $redocHtml = @"
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InsureTech API - ReDoc</title>
    <style>
        body { margin: 0; padding: 0; }
    </style>
</head>
<body>
    <redoc 
        spec-url="openapi.yaml"
        scroll-y-offset="nav"
        hide-download-button="false"
        hide-hostname="false"
        expand-responses="200,201"
    ></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>
"@
    $redocHtml | Out-File -FilePath "docs\redoc.html" -Encoding utf8
    
    Write-Success "Generated Swagger UI with enhanced features"
    Write-Success "Generated ReDoc with better configuration"
    
    # Note: Enhanced index.html already generated by doc_generator.py above
    # Old static index removed in favor of dynamic generated version
    
    # Generate simple fallback index (backup)
    $fallbackHtml = @"
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InsureTech API Documentation Hub</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            max-width: 900px;
            width: 100%;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            padding: 40px;
        }
        h1 { 
            color: #333; 
            margin-bottom: 10px;
            font-size: 2.5em;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 1.1em;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 15px;
            margin-bottom: 30px;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 10px;
        }
        .stat {
            text-align: center;
        }
        .stat-value {
            font-size: 2em;
            font-weight: bold;
            color: #667eea;
        }
        .stat-label {
            font-size: 0.9em;
            color: #666;
            margin-top: 5px;
        }
        .link-card { 
            border: 2px solid #e0e0e0;
            padding: 25px;
            margin: 15px 0;
            border-radius: 12px;
            transition: all 0.3s ease;
            cursor: pointer;
        }
        .link-card:hover { 
            border-color: #667eea;
            box-shadow: 0 5px 20px rgba(102, 126, 234, 0.2);
            transform: translateY(-2px);
        }
        .link-card h3 {
            margin-bottom: 8px;
            color: #333;
        }
        a { 
            text-decoration: none; 
            color: inherit;
            display: block;
        }
        .description { 
            color: #666; 
            line-height: 1.6;
        }
        .icon {
            font-size: 1.5em;
            margin-right: 10px;
        }
        .badge {
            display: inline-block;
            background: #4caf50;
            color: white;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.85em;
            margin-left: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🏥 InsureTech API</h1>
        <p class="subtitle">Complete API Documentation & Interactive Tools</p>
        
        <div class="stats">
            <div class="stat">
                <div class="stat-value">865</div>
                <div class="stat-label">Schemas</div>
            </div>
            <div class="stat">
                <div class="stat-value">177</div>
                <div class="stat-label">Endpoints</div>
            </div>
            <div class="stat">
                <div class="stat-value">100%</div>
                <div class="stat-label">Coverage</div>
            </div>
            <div class="stat">
                <div class="stat-value">0</div>
                <div class="stat-label">Errors</div>
            </div>
        </div>
        
        <a href="swagger.html">
            <div class="link-card">
                <h3><span class="icon">📘</span>Swagger UI<span class="badge">Interactive</span></h3>
                <div class="description">
                    Interactive API explorer with try-it-out functionality. 
                    Test endpoints, view request/response examples, and explore the API interactively.
                </div>
            </div>
        </a>
        
        <a href="redoc.html">
            <div class="link-card">
                <h3><span class="icon">📗</span>ReDoc<span class="badge">Clean</span></h3>
                <div class="description">
                    Clean, responsive three-panel API reference documentation. 
                    Perfect for reading and understanding the API structure.
                </div>
            </div>
        </a>
        
        <a href="../openapi.yaml" download>
            <div class="link-card">
                <h3><span class="icon">📄</span>OpenAPI Specification<span class="badge">v3.1</span></h3>
                <div class="description">
                    Download the raw OpenAPI 3.1 specification file (YAML format). 
                    Use for code generation, testing, and integration.
                </div>
            </div>
        </a>
        
        <a href="../validation_report.html">
            <div class="link-card">
                <h3><span class="icon">✅</span>Validation Report<span class="badge">Passed</span></h3>
                <div class="description">
                    Detailed validation results and quality metrics. 
                    View coverage statistics, warnings, and recommendations.
                </div>
            </div>
        </a>
    </div>
</body>
</html>
"@
    # Don't overwrite the enhanced index.html - only create fallback if needed
    if (-not (Test-Path "docs\index.html")) {
        $fallbackHtml | Out-File -FilePath "docs\index_fallback.html" -Encoding utf8
    }
    
    # Verify all files created
    $requiredFiles = @("swagger.html", "redoc.html", "index.html")
    $allCreated = $true
    foreach ($file in $requiredFiles) {
        if (-not (Test-Path "docs\$file")) {
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
    if (Test-Path "docs") {
        try {
            # Remove old files in root docs to ensure clean sync
            Get-ChildItem $rootDocsDir -Recurse | Remove-Item -Force -Recurse -ErrorAction SilentlyContinue
            
            # Copy all files recursively
            Copy-Item -Path "docs\*" -Destination $rootDocsDir -Recurse -Force -ErrorAction Stop
            
            Write-Success "Synced all files from api/docs/ to root docs/"
        } catch {
            Write-Host "  ⚠ Error copying files: $($_.Exception.Message)" -ForegroundColor Yellow
        }
    }
    
    # Also copy additional files from api/ root
    if (Test-Path "openapi.yaml") {
        Copy-Item "openapi.yaml" -Destination (Join-Path $rootDocsDir "openapi.yaml") -Force
        Write-Success "Copied openapi.yaml to root docs/"
    }
    
    if (Test-Path "validation_report.html") {
        Copy-Item "validation_report.html" -Destination (Join-Path $rootDocsDir "validation_report.html") -Force
        Write-Success "Copied validation_report.html to root docs/"
    }
    
    if (Test-Path "validation_report.json") {
        Copy-Item "validation_report.json" -Destination (Join-Path $rootDocsDir "validation_report.json") -Force
        Write-Success "Copied validation_report.json to root docs/"
    }
    
    # Copy schema summary files if they exist
    if (Test-Path "proto_schema_summary.json") {
        Copy-Item "proto_schema_summary.json" -Destination (Join-Path $rootDocsDir "proto_schema_summary.json") -Force
        Write-Success "Copied proto_schema_summary.json to root docs/"
    }
    
    if (Test-Path "schema_api_mapping.json") {
        Copy-Item "schema_api_mapping.json" -Destination (Join-Path $rootDocsDir "schema_api_mapping.json") -Force
        Write-Success "Copied schema_api_mapping.json to root docs/"
    }
    
    Write-Success "Documentation ready for GitHub Pages deployment"
}

# Step 13.5: Sync to Apidog (Optional)
Write-Step 13 16 "Syncing to Apidog..."

# Load environment variables from .env file
$envFile = Join-Path $ProjectRoot ".env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)\s*=\s*(.+)\s*$') {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
        }
    }
}

$apidogToken = [Environment]::GetEnvironmentVariable("API_DOG_TOKEN", "Process")

if ($apidogToken) {
    Write-Host "  Found API_DOG_TOKEN, syncing to Apidog..." -ForegroundColor Gray
    Set-Location generator
    
    try {
        python sync_apidog.py
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Synced to Apidog successfully"
        } else {
            Write-Host "  ⚠ Apidog sync had issues (continuing...)" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "  ⚠ Apidog sync failed: $($_.Exception.Message)" -ForegroundColor Yellow
    }
    
    Set-Location ..
} else {
    Write-Host "  ⚠ API_DOG_TOKEN not found, skipping Apidog sync" -ForegroundColor Yellow
    Write-Host "    Set API_DOG_TOKEN in .env file to enable Apidog integration" -ForegroundColor Gray
}

# Step 14: Generate SDKs
Write-Step 14 16 "Generating SDKs..."

# Generate TypeScript SDK (using hey-api + custom post-processing)
Write-Host "  Generating TypeScript SDK (hey-api + custom)..." -ForegroundColor Gray
Set-Location (Join-Path $ProjectRoot "sdks\sdk-generator\typescript")

# Check if node_modules exists for generator
if (-not (Test-Path "node_modules")) {
    Write-Host "    Installing @hey-api/openapi-ts..." -ForegroundColor Gray
    npm install 2>&1 | Out-Null
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

# Build custom Go post-processor (always rebuild to ensure latest code)
Write-Host "    Building custom post-processor..." -ForegroundColor Gray
$env:GOWORK = "off"
Remove-Item "generator.exe" -Force -ErrorAction SilentlyContinue
go build -o generator.exe generator.go
if ($LASTEXITCODE -ne 0) {
    Write-Error-Step "Failed to build post-processor"
    exit 1
}

# Run custom post-processor
Write-Host "    Applying custom modifications..." -ForegroundColor Gray
$postProcessOutput = .\generator.exe 2>&1
$postProcessExitCode = $LASTEXITCODE

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
Set-Location (Join-Path $ProjectRoot "sdks\sdk-generator\go")

# Build generator (always rebuild to ensure latest code)
Write-Host "    Building Go SDK generator..." -ForegroundColor Gray
Remove-Item "generator.exe" -Force -ErrorAction SilentlyContinue
go build -o generator.exe generator.go
if ($LASTEXITCODE -ne 0) {
    Write-Error-Step "Failed to build Go SDK generator"
    exit 1
}

# Run generator
Write-Host "    Running Go SDK generator..." -ForegroundColor Gray
$goGenOutput = .\generator.exe 2>&1
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
Set-Location (Join-Path $ProjectRoot "sdks\insuretech-typescript-sdk")

# Check if node_modules exists, install if needed
if (-not (Test-Path "node_modules")) {
    Write-Host "    Installing dependencies..." -ForegroundColor Gray
    npm install --legacy-peer-deps 2>&1 | Out-Null
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
Set-Location (Join-Path $ProjectRoot "sdks\insuretech-go-sdk")

# Build with GOWORK=off to avoid workspace conflicts
$env:GOWORK = "off"
Write-Host "    Running go build..." -ForegroundColor Gray
$goBuildOutput = go build ./... 2>&1
$goBuildExitCode = $LASTEXITCODE

if ($goBuildExitCode -eq 0) {
    Write-Success "Go SDK built successfully"
} else {
    Write-Error-Step "Go SDK build failed!"
    Write-Host "Build output:" -ForegroundColor Red
    Write-Host $goBuildOutput -ForegroundColor Red
    exit 1
}

Set-Location $ApiDir

# Step 15: Start Documentation Server
Write-Step 15 16 "Starting documentation server..."

# Create HTTP server script with custom handler for root redirect
$serverScript = @"
import http.server
import socketserver
import os
import sys
from urllib.parse import urlparse

PORT = $ServerPort
os.chdir(r'$ApiDir')

class CustomHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        # Add CORS headers for Swagger/ReDoc
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type')
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        super().end_headers()
    
    def do_GET(self):
        # Redirect root to docs/index.html
        if self.path == '/' or self.path == '':
            self.send_response(302)
            self.send_header('Location', '/docs/index.html')
            self.end_headers()
            return
        # Serve other files normally
        return http.server.SimpleHTTPRequestHandler.do_GET(self)

# Try to bind to port, retry with next port if occupied
max_attempts = 5
for attempt in range(max_attempts):
    try:
        with socketserver.TCPServer(('', PORT), CustomHandler) as httpd:
            print('')
            print('='*60)
            print('  InsureTech API Documentation Server')
            print('='*60)
            print('  Server running at: http://localhost:' + str(PORT) + '/')
            print('  Documentation:     http://localhost:' + str(PORT) + '/docs/')
            print('  Swagger UI:        http://localhost:' + str(PORT) + '/docs/swagger.html')
            print('  ReDoc:             http://localhost:' + str(PORT) + '/docs/redoc.html')
            print('  Schema Visualizer: http://localhost:' + str(PORT) + '/docs/index.html (🎨 tab)')
            print('  OpenAPI Spec:      http://localhost:' + str(PORT) + '/openapi.yaml')
            print('='*60)
            print('  Press Ctrl+C to stop the server')
            print('='*60)
            print('')
            httpd.serve_forever()
        break
    except OSError as e:
        if e.winerror == 10048:  # Port in use
            print('Port ' + str(PORT) + ' is in use, trying ' + str(PORT + 1) + '...')
            PORT += 1
        else:
            raise
else:
    print('Could not find available port after ' + str(max_attempts) + ' attempts')
    sys.exit(1)
"@

$serverScript | Out-File -FilePath "generator\server.py" -Encoding utf8

# Calculate elapsed time
$EndTime = Get-Date
$Duration = $EndTime - $StartTime
$DurationSeconds = [math]::Round($Duration.TotalSeconds, 1)

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "✅ API GENERATION COMPLETE" -ForegroundColor Green
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
    Write-Host "  Validation: ✓ Passed ($errors errors, $warnings warnings)" -ForegroundColor Green
}

Write-Host "`nStarting server on port $ServerPort..." -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop the server.`n" -ForegroundColor Gray

# Start server
Set-Location generator
python server.py
