#!/usr/bin/env pwsh
# Check service status on remote server

param(
    [string]$RemoteHost = "insureadmin@146.190.97.242"
)

$ErrorActionPreference = "Stop"

Write-Host "🔍 Checking InsureTech services..." -ForegroundColor Cyan

$sshCmd = if (Get-Command ssh.exe -ErrorAction SilentlyContinue) { "ssh.exe" } else { "ssh" }
$scpCmd = if (Get-Command scp.exe -ErrorAction SilentlyContinue) { "scp.exe" } else { "scp" }

Write-Host "Uploading check script..." -ForegroundColor Yellow
& $scpCmd "scripts/check_services_remote.sh" "${RemoteHost}:/tmp/check_services.sh"

Write-Host "Running diagnostics..." -ForegroundColor Yellow
& $sshCmd $RemoteHost "dos2unix /tmp/check_services.sh 2>/dev/null || sed -i 's/\r$//' /tmp/check_services.sh && chmod +x /tmp/check_services.sh && bash /tmp/check_services.sh"

Write-Host "`n✅ Done!" -ForegroundColor Green

