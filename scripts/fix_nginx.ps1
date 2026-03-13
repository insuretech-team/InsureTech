#!/usr/bin/env pwsh
# Fix nginx configuration conflicts

param(
    [string]$RemoteHost = "insureadmin@146.190.97.242"
)

$ErrorActionPreference = "Stop"

Write-Host "🔧 Fixing nginx configuration conflicts..." -ForegroundColor Cyan

$sshCmd = if (Get-Command ssh.exe -ErrorAction SilentlyContinue) { "ssh.exe" } else { "ssh" }
$scpCmd = if (Get-Command scp.exe -ErrorAction SilentlyContinue) { "scp.exe" } else { "scp" }

Write-Host "Uploading fix script..." -ForegroundColor Yellow
& $scpCmd "scripts/fix_nginx_conflicts.sh" "${RemoteHost}:/tmp/fix_nginx.sh"

Write-Host "Executing fix..." -ForegroundColor Yellow
& $sshCmd -t $RemoteHost "dos2unix /tmp/fix_nginx.sh 2>/dev/null || sed -i 's/\r$//' /tmp/fix_nginx.sh && chmod +x /tmp/fix_nginx.sh && bash /tmp/fix_nginx.sh"

Write-Host "`n✅ Done!" -ForegroundColor Green
