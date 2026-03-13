# Generate Python protobuf files for google.api annotations
# This is needed to parse http annotations from proto descriptors

Write-Host "Generating Python protobuf files for google.api..." -ForegroundColor Yellow

# Get project root
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $ScriptDir)

# Load BUF_TOKEN from .env if it exists
$envFile = Join-Path $ProjectRoot ".env"
if (Test-Path $envFile) {
    $bufToken = Get-Content $envFile | Where-Object { $_ -match "^BUF_TOKEN=" } | ForEach-Object { $_ -replace "^BUF_TOKEN=", "" }
    if ($bufToken) {
        $env:BUF_TOKEN = $bufToken
        Write-Host "Loaded BUF_TOKEN from .env" -ForegroundColor Gray
    }
}

# Create gen directory
$genDir = Join-Path $ScriptDir "gen"
New-Item -ItemType Directory -Force -Path $genDir | Out-Null

# Change to project root to use buf
Push-Location $ProjectRoot

try {
    # Use buf to generate Python files for google.api only
    Write-Host "Running buf generate for google.api annotations..." -ForegroundColor Cyan
    
    $bufGenConfig = @"
version: v2
plugins:
  - remote: buf.build/protocolbuffers/python:v28.3
    out: $($genDir.Replace('\', '/'))
"@
    
    # Save temporary buf.gen.yaml
    $tempBufGen = Join-Path $env:TEMP "buf.gen.google-api.yaml"
    $bufGenConfig | Out-File -FilePath $tempBufGen -Encoding UTF8
    
    # Generate only for googleapis dependency
    buf generate buf.build/googleapis/googleapis `
        --template $tempBufGen `
        --path google/api/annotations.proto `
        --path google/api/http.proto `
        2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✓ Google API protos generated successfully" -ForegroundColor Green
        Write-Host "  Location: $genDir" -ForegroundColor Gray
        
        # Create __init__.py files
        Get-ChildItem -Path $genDir -Recurse -Directory | ForEach-Object {
            $initFile = Join-Path $_.FullName "__init__.py"
            if (-not (Test-Path $initFile)) {
                New-Item -ItemType File -Path $initFile -Force | Out-Null
            }
        }
        
        Write-Host "✓ Created __init__.py files" -ForegroundColor Green
    } else {
        Write-Host "✗ Failed to generate protos" -ForegroundColor Red
        exit 1
    }
    
    # Cleanup
    Remove-Item $tempBufGen -ErrorAction SilentlyContinue
    
} finally {
    Pop-Location
}

Write-Host "`n✅ Setup complete! Run the API pipeline again." -ForegroundColor Green
