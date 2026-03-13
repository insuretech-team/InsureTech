#!/usr/bin/env pwsh
# Nginx Setup Script for InsureTech Bare-Metal Deployment
# Configures nginx reverse proxy for the API gateway

param(
    [string]$RemoteHost = "insureadmin@146.190.97.242",
    [string]$Domain = "labaidinsuretech.com"
)

$ErrorActionPreference = "Stop"

Write-Host "🔧 Setting up Nginx for InsureTech..." -ForegroundColor Cyan

# Detect SSH command
$sshCmd = if (Get-Command ssh.exe -ErrorAction SilentlyContinue) { "ssh.exe" } else { "ssh" }
$scpCmd = if (Get-Command scp.exe -ErrorAction SilentlyContinue) { "scp.exe" } else { "scp" }

Write-Host "Using SSH: $sshCmd" -ForegroundColor Gray

Write-Host "`n📋 Step 1: Uploading nginx configurations..." -ForegroundColor Yellow

# Upload gateway upstream config
& $scpCmd "backend/infra/nginx/upstreams/gateway.conf" "${RemoteHost}:/tmp/gateway.conf"

# Upload API site config
& $scpCmd "backend/infra/nginx/sites-available/insuretech-api.conf" "${RemoteHost}:/tmp/insuretech-api.conf"

# Upload setup script
& $scpCmd "scripts/setup_nginx_remote.sh" "${RemoteHost}:/tmp/setup_nginx.sh"

Write-Host "✓ Configurations uploaded" -ForegroundColor Green

Write-Host "`n⚙️  Step 2: Installing and configuring nginx..." -ForegroundColor Yellow

# Convert line endings and execute
& $sshCmd -t $RemoteHost "dos2unix /tmp/setup_nginx.sh 2>/dev/null || sed -i 's/\r$//' /tmp/setup_nginx.sh && chmod +x /tmp/setup_nginx.sh && bash /tmp/setup_nginx.sh"

Write-Host "`n✅ Nginx Setup Complete!" -ForegroundColor Green
Write-Host "`nTest the API gateway:" -ForegroundColor Cyan
Write-Host "  curl http://146.190.97.242/healthz" -ForegroundColor White
Write-Host "  curl http://146.190.97.242/v1/health" -ForegroundColor White
Write-Host "`nTo enable HTTPS:" -ForegroundColor Cyan
Write-Host "  1. Ensure DNS points to 146.190.97.242" -ForegroundColor White
Write-Host "  2. Run: ssh $RemoteHost 'sudo certbot --nginx -d $Domain -d www.$Domain'" -ForegroundColor White
Write-Host "  3. Uncomment HTTPS server block in /etc/nginx/sites-available/insuretech-api.conf" -ForegroundColor White

