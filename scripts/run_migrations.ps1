<#
.SYNOPSIS
Run database migrations on remote server

.DESCRIPTION
Executes database migrations for both primary (DigitalOcean) and backup (Neon) databases.
Assumes .env.prod already synced to remote server.

.PARAMETER Target
Database target: primary, backup, or both (default: both)

.EXAMPLE
.\scripts\run_migrations.ps1
.\scripts\run_migrations.ps1 -Target primary
#>

param (
    [string]$RemoteHost = "insureadmin@146.190.97.242",
    [string]$RemoteDir = "/home/insureadmin/insuretech",
    [ValidateSet("primary", "backup", "both")]
    [string]$Target = "both"
)

$ErrorActionPreference = "Stop"

Write-Host "🗄️  InsureTech Database Migration" -ForegroundColor Cyan
Write-Host "Target: $Target" -ForegroundColor Yellow
Write-Host ""

# Check for SSH and SCP
$sshCmd = Get-Command ssh.exe -ErrorAction SilentlyContinue
$scpCmd = Get-Command scp.exe -ErrorAction SilentlyContinue

if (-not $sshCmd -or -not $scpCmd) {
    Write-Error "SSH/SCP not found. Please install OpenSSH Client via Windows Settings or use WSL."
    exit 1
}

# Build dbops tool locally
Write-Host "📦 Building dbops tool..." -ForegroundColor Green

$ProjectRoot = Resolve-Path "$PSScriptRoot\.." | Select-Object -ExpandProperty Path
$BackendDir = Join-Path $ProjectRoot "backend\inscore"
$BuildDir = Join-Path $ProjectRoot "build_dbops"

if (Test-Path $BuildDir) {
    Remove-Item -Recurse -Force $BuildDir
}
New-Item -ItemType Directory -Path $BuildDir | Out-Null

Push-Location $BackendDir

$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"

go build -ldflags="-w -s" -o "$BuildDir\dbops" .\cmd\dbops\main.go

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to build dbops"
    Pop-Location
    exit 1
}

Pop-Location

Write-Host "  ✓ Built dbops" -ForegroundColor Green

# Transfer to remote
Write-Host "`n📤 Transferring to remote server..." -ForegroundColor Green

& ssh.exe $RemoteHost "mkdir -p $RemoteDir/bin"
& scp.exe "$BuildDir\dbops" "$RemoteHost`:$RemoteDir/bin/"

Write-Host "  ✓ Transferred" -ForegroundColor Green

# Run migrations
Write-Host "`n🚀 Running migrations on remote server..." -ForegroundColor Green

$migrationScript = @"
#!/bin/bash
set -e

cd $RemoteDir

# Load environment
if [ -f .env ]; then
    export \$(cat .env | grep -v '^#' | xargs)
fi

echo "→ Running migrations (target: $Target)..."
./bin/dbops migrate --target=$Target

echo ""
echo "→ Validating schema consistency..."
./bin/dbops validate || echo "⚠ Schema mismatch detected"

echo ""
echo "✓ Migration complete"
"@

$tempScript = [System.IO.Path]::GetTempFileName()
$migrationScript | Out-File -FilePath $tempScript -Encoding ASCII

& scp.exe $tempScript "$RemoteHost`:$RemoteDir/run_migration.sh"
& ssh.exe -t $RemoteHost "chmod +x $RemoteDir/run_migration.sh && $RemoteDir/run_migration.sh && rm $RemoteDir/run_migration.sh"

Remove-Item $tempScript

# Cleanup
Remove-Item -Recurse -Force $BuildDir

Write-Host "`n✅ Migration Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Database target: $Target" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Verify database: ssh.exe $RemoteHost 'cd $RemoteDir && ./bin/dbops validate'" -ForegroundColor White
Write-Host "  2. Check migration status: ssh.exe $RemoteHost 'cd $RemoteDir && ./bin/dbops status'" -ForegroundColor White
Write-Host ""
