<#
.SYNOPSIS
Deploys Nginx custom configurations and pages to the remote server via SCP and executes the setup script.

.DESCRIPTION
This script avoids the need for Git by manually transferring the necessary configuration files, 
custom web pages, logos, and the bash installation script to the remote server's /tmp directory, 
and then executing the setup script as root.
#>

param (
    [string]$RemoteHost = "root@146.190.97.242",
    [string]$DeployDir = "/tmp/insuretech_nginx_extend"
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 Starting Direct SCP Deployment to $RemoteHost..." -ForegroundColor Cyan

# Define local paths relative to this script's directory
$ProjectRoot = Resolve-Path "$PSScriptRoot\.." | Select-Object -ExpandProperty Path
$NginxConfDir = Join-Path $ProjectRoot "backend\infra\nginx"
$WebSharedDir = Join-Path $ProjectRoot "web_shared"
$SetupScript = Join-Path $ProjectRoot "scripts\setup_remote_nginx.sh"

# 1. Create remote deployment directory tree
Write-Host "`n📁 Step 1: Creating temporary deployment directories on remote server..."
ssh $RemoteHost "mkdir -p $DeployDir/backend/infra/nginx/sites-available $DeployDir/backend/infra/nginx/static $DeployDir/web_shared/pages $DeployDir/scripts"

# 2. Securely Copy the required assets
Write-Host "`n📤 Step 2: Transferring configurations and websites over SCP..."

# Nginx configurations
scp -r "$NginxConfDir/sites-available" "$RemoteHost`:$DeployDir/backend/infra/nginx/"
# Static Nginx artifacts (logos)
scp -r "$NginxConfDir/static" "$RemoteHost`:$DeployDir/backend/infra/nginx/"

# Web Shared Pages (404, 500, coming_soon)
scp -r "$WebSharedDir/pages" "$RemoteHost`:$DeployDir/web_shared/"

# The bash script
scp "$SetupScript" "$RemoteHost`:$DeployDir/scripts/"

# 3. Execute remote setup script
Write-Host "`n⚙️  Step 3: Executing setup script on remote server as root..."
# We pass the PROJECT_ROOT explicitly so the bash script knows where the temporal deployment is hosted.
# Since we are logging in as root, we drop the 'sudo' prefix.
$RemoteCommand = "chmod +x $DeployDir/scripts/setup_remote_nginx.sh && PROJECT_ROOT=$DeployDir $DeployDir/scripts/setup_remote_nginx.sh"
ssh -t $RemoteHost $RemoteCommand

Write-Host "`n🧹 Step 4: Cleaning up temporary deployment files on remote server..."
ssh $RemoteHost "rm -rf $DeployDir"

Write-Host "`n🎉 Direct Deployment completed successfully!" -ForegroundColor Green
