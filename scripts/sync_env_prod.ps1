<#
.SYNOPSIS
Syncs the local .env.prod file to the remote server's target directory as .env.

.DESCRIPTION
This script safely transports your local production variables layout directly 
to the remote VM, ensuring all microservices receive the correct configurations.
#>

param (
    [string]$RemoteHost = "insureadmin@146.190.97.242",
    [string]$RemoteDir = "/home/insureadmin/insuretech"
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 Starting Direct .env.prod Sync to $RemoteHost..." -ForegroundColor Cyan

# Resolve local paths
$ProjectRoot = Resolve-Path "$PSScriptRoot\.." | Select-Object -ExpandProperty Path
$EnvProdFile = Join-Path $ProjectRoot ".env.prod"

if (-not (Test-Path $EnvProdFile)) {
    Write-Error "❌ Could not find .env.prod at $EnvProdFile"
    exit 1
}

Write-Host "`n📁 Step 1: Ensuring remote directory exists..."
ssh $RemoteHost "mkdir -p $RemoteDir"

Write-Host "`n📤 Step 2: Transferring .env.prod to $RemoteDir/.env and $RemoteDir/.env.prod over SCP..."
scp "$EnvProdFile" "$RemoteHost`:$RemoteDir/.env"
scp "$EnvProdFile" "$RemoteHost`:$RemoteDir/.env.prod"

Write-Host "`n🎉 Environments synced successfully!" -ForegroundColor Green
