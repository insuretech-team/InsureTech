<#
.SYNOPSIS
Deploy InsureTech services to bare-metal server (146.190.97.242)

.DESCRIPTION
Builds Go services and B2B Portal, then deploys to production server.
Assumes .env.prod already synced and nginx already configured.

.PARAMETER Services
Comma-separated list of services to deploy. Default: gateway,authn,authz,tenant,b2b_portal

.EXAMPLE
.\scripts\deploy_baremetal.ps1
.\scripts\deploy_baremetal.ps1 -Services "gateway,authn"
#>

param (
    [string]$RemoteHost = "root@146.190.97.242",
    [string]$RemoteDir = "/home/insureadmin/insuretech",
    [string]$Services = "gateway,authn,authz,tenant",
    [string]$RemoteUser = "insureadmin"
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 InsureTech Bare-Metal Deployment" -ForegroundColor Cyan
Write-Host "Target: $RemoteHost" -ForegroundColor Yellow
Write-Host "Services: $Services" -ForegroundColor Yellow
Write-Host ""

# Check for SSH and SCP (try both .exe and without extension)
$sshCmd = Get-Command ssh.exe -ErrorAction SilentlyContinue
if (-not $sshCmd) {
    $sshCmd = Get-Command ssh -ErrorAction SilentlyContinue
}

$scpCmd = Get-Command scp.exe -ErrorAction SilentlyContinue
if (-not $scpCmd) {
    $scpCmd = Get-Command scp -ErrorAction SilentlyContinue
}

if (-not $sshCmd -or -not $scpCmd) {
    Write-Error @"
SSH/SCP not found in PATH. Please install OpenSSH:

Option 1: Install via Windows Settings
  Settings > Apps > Optional Features > Add a feature > OpenSSH Client

Option 2: Install via PowerShell (Admin)
  Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0

Option 3: Use Git Bash or WSL
  - Git for Windows includes SSH: https://git-scm.com/download/win
  - Or use WSL: wsl ./scripts/quick_deploy.sh

After installation, restart PowerShell.
"@
    exit 1
}

# Use the found commands
$ssh = $sshCmd.Source
$scp = $scpCmd.Source

Write-Host "✓ SSH/SCP found: $ssh" -ForegroundColor Green
Write-Host ""

# Resolve paths
$ProjectRoot = Resolve-Path "$PSScriptRoot\.." | Select-Object -ExpandProperty Path
$BackendDir = Join-Path $ProjectRoot "backend\inscore"
$B2BPortalDir = Join-Path $ProjectRoot "b2b_portal"
$BuildDir = Join-Path $ProjectRoot "build"

# Create build directory
if (Test-Path $BuildDir) {
    Remove-Item -Recurse -Force $BuildDir
}
New-Item -ItemType Directory -Path $BuildDir | Out-Null
New-Item -ItemType Directory -Path "$BuildDir\bin" | Out-Null

$ServiceList = $Services -split ","


# ============================================================================
# STEP 1: Build Go Services
# ============================================================================
Write-Host "📦 Step 1: Building Go Services..." -ForegroundColor Green

Push-Location $BackendDir

foreach ($service in $ServiceList) {
    if ($service -eq "b2b_portal") { continue }
    
    Write-Host "  → Building $service..." -ForegroundColor Cyan
    
    $env:CGO_ENABLED = "0"
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    
    $outputPath = Join-Path $BuildDir "bin\$service"
    $mainPath = ".\cmd\$service\main.go"
    
    if (-not (Test-Path $mainPath)) {
        Write-Warning "  ⚠ Skipping $service - main.go not found at $mainPath"
        continue
    }
    
    go build -ldflags="-w -s" -o $outputPath $mainPath
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build $service"
        Pop-Location
        exit 1
    }
    
    Write-Host "  ✓ Built $service" -ForegroundColor Green
}

Pop-Location

# ============================================================================
# STEP 2: Copy Configuration Files
# ============================================================================
Write-Host "`n📋 Step 2: Copying Configuration Files..." -ForegroundColor Green

$ConfigDirs = @("configs", "secrets", "templates")

foreach ($dir in $ConfigDirs) {
    $sourcePath = Join-Path $BackendDir $dir
    $destPath = Join-Path $BuildDir $dir
    
    if (Test-Path $sourcePath) {
        Copy-Item -Recurse -Force $sourcePath $destPath
        Write-Host "  ✓ Copied $dir" -ForegroundColor Green
    } else {
        Write-Warning "  ⚠ $dir not found at $sourcePath"
    }
}

# ============================================================================
# STEP 3: Build B2B Portal (if requested)
# ============================================================================
if ($Services -match "b2b_portal") {
    Write-Host "`n🌐 Step 3: Building B2B Portal..." -ForegroundColor Green
    
    Push-Location $B2BPortalDir
    
    Write-Host "  → Installing dependencies..." -ForegroundColor Cyan
    npm ci --production=false
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to install B2B Portal dependencies"
        Pop-Location
        exit 1
    }
    
    Write-Host "  → Building Next.js application..." -ForegroundColor Cyan
    npm run build
    
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to build B2B Portal"
        Pop-Location
        exit 1
    }
    
    # Copy build output
    $b2bBuildDir = Join-Path $BuildDir "b2b_portal"
    New-Item -ItemType Directory -Path $b2bBuildDir -Force | Out-Null
    
    Copy-Item -Recurse -Force ".next\standalone\*" $b2bBuildDir
    Copy-Item -Recurse -Force ".next\static" "$b2bBuildDir\.next\"
    Copy-Item -Recurse -Force "public" $b2bBuildDir
    Copy-Item -Force "package.json" $b2bBuildDir
    
    Write-Host "  ✓ Built B2B Portal" -ForegroundColor Green
    
    Pop-Location
} else {
    Write-Host "`n⏭  Step 3: Skipping B2B Portal (not in service list)" -ForegroundColor Yellow
}


# ============================================================================
# STEP 4: Create Deployment Package
# ============================================================================
Write-Host "`n📦 Step 4: Creating Deployment Package..." -ForegroundColor Green

$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$packageName = "insuretech_deploy_$timestamp.tar.gz"
$packagePath = Join-Path $ProjectRoot $packageName

Push-Location $BuildDir

# Create tar.gz using WSL or tar.exe (Windows 10+)
if (Get-Command wsl -ErrorAction SilentlyContinue) {
    Write-Host "  → Using WSL to create tar.gz..." -ForegroundColor Cyan
    wsl tar -czf "../$packageName" .
} elseif (Get-Command tar -ErrorAction SilentlyContinue) {
    Write-Host "  → Using Windows tar to create tar.gz..." -ForegroundColor Cyan
    tar -czf "..\$packageName" .
} else {
    Write-Error "Neither WSL nor tar command found. Cannot create deployment package."
    Pop-Location
    exit 1
}

Pop-Location

if (-not (Test-Path $packagePath)) {
    Write-Error "Failed to create deployment package"
    exit 1
}

$packageSize = (Get-Item $packagePath).Length / 1MB
Write-Host "  ✓ Created package: $packageName ($([math]::Round($packageSize, 2)) MB)" -ForegroundColor Green

# ============================================================================
# STEP 5: Transfer to Remote Server
# ============================================================================
Write-Host "`n📤 Step 5: Transferring to Remote Server..." -ForegroundColor Green

Write-Host "  → Ensuring remote directory exists..." -ForegroundColor Cyan
& $ssh $RemoteHost "mkdir -p $RemoteDir/releases"

Write-Host "  → Uploading package..." -ForegroundColor Cyan
& $scp $packagePath "$RemoteHost`:$RemoteDir/releases/"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to transfer package to remote server"
    exit 1
}

Write-Host "  ✓ Package uploaded" -ForegroundColor Green

# ============================================================================
# STEP 6: Extract and Deploy on Remote
# ============================================================================
Write-Host "`n🚀 Step 6: Deploying on Remote Server..." -ForegroundColor Green

$remoteScript = @"
#!/bin/bash
set -e

RELEASE_DIR="$RemoteDir/releases"
DEPLOY_DIR="$RemoteDir"
PACKAGE_NAME="$packageName"

echo "-> Extracting package..."
cd `$RELEASE_DIR
tar -xzf `$PACKAGE_NAME

echo "-> Backing up current deployment..."
if [ -d "`$DEPLOY_DIR/bin" ]; then
    mv `$DEPLOY_DIR/bin `$DEPLOY_DIR/bin.backup.`$(date +%Y%m%d_%H%M%S)
fi
if [ -d "`$DEPLOY_DIR/configs" ]; then
    mv `$DEPLOY_DIR/configs `$DEPLOY_DIR/configs.backup.`$(date +%Y%m%d_%H%M%S)
fi

echo "-> Moving new binaries and configs..."
mv bin `$DEPLOY_DIR/
mv configs `$DEPLOY_DIR/ 2>/dev/null || true
mv secrets `$DEPLOY_DIR/ 2>/dev/null || true
mv templates `$DEPLOY_DIR/ 2>/dev/null || true

echo "-> Setting permissions..."
chmod +x `$DEPLOY_DIR/bin/*
chmod 600 `$DEPLOY_DIR/secrets/*.pem 2>/dev/null || true
chmod 600 `$DEPLOY_DIR/.env 2>/dev/null || true
chown -R $RemoteUser`:$RemoteUser `$DEPLOY_DIR/bin `$DEPLOY_DIR/configs `$DEPLOY_DIR/secrets `$DEPLOY_DIR/templates 2>/dev/null || true

echo "-> Cleaning up extraction..."
cd `$RELEASE_DIR
rm -rf bin configs secrets templates

echo "Deployment files ready"
"@

# Save script to temp file and execute
$tempScript = [System.IO.Path]::GetTempFileName()
$remoteScript | Out-File -FilePath $tempScript -Encoding ASCII

& $scp $tempScript "$RemoteHost`:$RemoteDir/deploy_temp.sh"
& $ssh $RemoteHost "chmod +x $RemoteDir/deploy_temp.sh && $RemoteDir/deploy_temp.sh && rm $RemoteDir/deploy_temp.sh"

Remove-Item $tempScript

if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to deploy on remote server"
    exit 1
}

Write-Host "  ✓ Files deployed" -ForegroundColor Green


# ============================================================================
# STEP 7: Deploy B2B Portal (if built)
# ============================================================================
if ($Services -match "b2b_portal") {
    Write-Host "`n🌐 Step 7: Deploying B2B Portal..." -ForegroundColor Green
    
    $b2bPackage = "b2b_portal_$timestamp.tar.gz"
    $b2bPackagePath = Join-Path $ProjectRoot $b2bPackage
    
    Push-Location (Join-Path $BuildDir "b2b_portal")
    
    if (Get-Command wsl -ErrorAction SilentlyContinue) {
        wsl tar -czf "../../$b2bPackage" .
    } else {
        tar -czf "..\..\$b2bPackage" .
    }
    
    Pop-Location
    
    Write-Host "  → Uploading B2B Portal..." -ForegroundColor Cyan
    & $scp $b2bPackagePath "$RemoteHost`:$RemoteDir/releases/"
    
    $b2bDeployScript = @"
#!/bin/bash
set -e

RELEASE_DIR="$RemoteDir/releases"
B2B_DIR="/var/www/insuretech-b2b"
PACKAGE_NAME="$b2bPackage"

echo "-> Creating B2B Portal directory..."
mkdir -p `$B2B_DIR

echo "-> Backing up current B2B Portal..."
if [ -d "`$B2B_DIR/.next" ]; then
    mv `$B2B_DIR `$B2B_DIR.backup.`$(date +%Y%m%d_%H%M%S)
    mkdir -p `$B2B_DIR
fi

echo "-> Extracting B2B Portal..."
cd `$RELEASE_DIR
tar -xzf `$PACKAGE_NAME
mv * `$B2B_DIR/ 2>/dev/null || true
chown -R $RemoteUser`:$RemoteUser `$B2B_DIR

echo "B2B Portal deployed"
"@
    
    $tempB2BScript = [System.IO.Path]::GetTempFileName()
    $b2bDeployScript | Out-File -FilePath $tempB2BScript -Encoding ASCII
    
    & $scp $tempB2BScript "$RemoteHost`:$RemoteDir/deploy_b2b_temp.sh"
    & $ssh $RemoteHost "chmod +x $RemoteDir/deploy_b2b_temp.sh && $RemoteDir/deploy_b2b_temp.sh && rm $RemoteDir/deploy_b2b_temp.sh"
    
    Remove-Item $tempB2BScript
    Remove-Item $b2bPackagePath
    
    Write-Host "  ✓ B2B Portal deployed" -ForegroundColor Green
} else {
    Write-Host "`n⏭  Step 7: Skipping B2B Portal deployment" -ForegroundColor Yellow
}

# ============================================================================
# STEP 8: Create/Update Systemd Services
# ============================================================================
Write-Host "`n⚙️  Step 8: Setting up Systemd Services..." -ForegroundColor Green

$serviceListStr = $ServiceList -join ' '
$systemdScript = @"
#!/bin/bash
set -e

SERVICES=($serviceListStr)

for service in "`${SERVICES[@]}"; do
    if [ "`$service" = "b2b_portal" ]; then
        # B2B Portal systemd service
        sudo tee /etc/systemd/system/insuretech-b2b-portal.service > /dev/null <<'EOFB2B'
[Unit]
Description=InsureTech B2B Portal
After=network.target

[Service]
Type=simple
User=insureadmin
Group=insureadmin
WorkingDirectory=/var/www/insuretech-b2b
Environment="NODE_ENV=production"
Environment="PORT=3000"
Environment="HOSTNAME=0.0.0.0"
ExecStart=/usr/bin/node server.js
Restart=always
RestartSec=10
StandardOutput=append:$RemoteDir/logs/b2b-portal.log
StandardError=append:$RemoteDir/logs/b2b-portal.error.log

NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOFB2B
        continue
    fi
    
    # Go service systemd template
    sudo tee /etc/systemd/system/insuretech-`$service.service > /dev/null <<EOFGO
[Unit]
Description=InsureTech `${service^} Service
After=network.target

[Service]
Type=simple
User=insureadmin
Group=insureadmin
WorkingDirectory=$RemoteDir
EnvironmentFile=$RemoteDir/.env
ExecStart=$RemoteDir/bin/`$service
Restart=always
RestartSec=10
StandardOutput=append:$RemoteDir/logs/`$service.log
StandardError=append:$RemoteDir/logs/`$service.error.log

NoNewPrivileges=true
PrivateTmp=true
ReadWritePaths=$RemoteDir/logs $RemoteDir/data

LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
EOFGO
done

echo "-> Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Systemd services configured"
"@

$tempSystemdScript = [System.IO.Path]::GetTempFileName()
$systemdScript | Out-File -FilePath $tempSystemdScript -Encoding ASCII

& $scp $tempSystemdScript "$RemoteHost`:$RemoteDir/setup_systemd.sh"
& $ssh $RemoteHost "chmod +x $RemoteDir/setup_systemd.sh && $RemoteDir/setup_systemd.sh && rm $RemoteDir/setup_systemd.sh"

Remove-Item $tempSystemdScript

Write-Host "  ✓ Systemd services configured" -ForegroundColor Green


# ============================================================================
# STEP 9: Restart Services
# ============================================================================
Write-Host "`n🔄 Step 9: Restarting Services..." -ForegroundColor Green

$serviceListStr = $ServiceList -join ' '
$restartScript = @"
#!/bin/bash
set -e

SERVICES=($serviceListStr)

echo "-> Stopping services..."
for service in "`${SERVICES[@]}"; do
    if [ "`$service" = "b2b_portal" ]; then
        sudo systemctl stop insuretech-b2b-portal 2>/dev/null || true
    else
        sudo systemctl stop insuretech-`$service 2>/dev/null || true
    fi
done

sleep 2

echo "-> Starting services..."
for service in "`${SERVICES[@]}"; do
    if [ "`$service" = "b2b_portal" ]; then
        sudo systemctl enable insuretech-b2b-portal
        sudo systemctl start insuretech-b2b-portal
    else
        sudo systemctl enable insuretech-`$service
        sudo systemctl start insuretech-`$service
    fi
done

sleep 3

echo ""
echo "-> Checking service status..."
for service in "`${SERVICES[@]}"; do
    if [ "`$service" = "b2b_portal" ]; then
        if sudo systemctl is-active --quiet insuretech-b2b-portal; then
            echo "  OK insuretech-b2b-portal RUNNING"
        else
            echo "  FAIL insuretech-b2b-portal FAILED"
        fi
    else
        if sudo systemctl is-active --quiet insuretech-`$service; then
            echo "  OK insuretech-`$service RUNNING"
        else
            echo "  FAIL insuretech-`$service FAILED"
        fi
    fi
done

echo ""
echo "-> Port status:"
sudo netstat -tlnp | grep -E '(8080|3000|5006[01]|5007[01]|5005[01])' || echo "  No services listening on expected ports"

echo ""
echo "Service restart complete"
"@

$tempRestartScript = [System.IO.Path]::GetTempFileName()
$restartScript | Out-File -FilePath $tempRestartScript -Encoding ASCII

& $scp $tempRestartScript "$RemoteHost`:$RemoteDir/restart_services.sh"
& $ssh -t $RemoteHost "chmod +x $RemoteDir/restart_services.sh && $RemoteDir/restart_services.sh && rm $RemoteDir/restart_services.sh"

Remove-Item $tempRestartScript

# ============================================================================
# STEP 10: Cleanup
# ============================================================================
Write-Host "`n🧹 Step 10: Cleaning up..." -ForegroundColor Green

Remove-Item -Recurse -Force $BuildDir
Remove-Item $packagePath

Write-Host "  ✓ Local build artifacts cleaned" -ForegroundColor Green

# ============================================================================
# COMPLETION
# ============================================================================
Write-Host "`n✅ Deployment Complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Services deployed: $Services" -ForegroundColor Cyan
Write-Host "Remote host: $RemoteHost" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Check service logs: ssh $RemoteHost 'tail -f $RemoteDir/logs/*.log'" -ForegroundColor White
Write-Host "  2. Test API Gateway: curl http://146.190.97.242:8080/healthz" -ForegroundColor White
Write-Host "  3. Test B2B Portal: curl http://146.190.97.242:3000" -ForegroundColor White
Write-Host "  4. Check nginx: ssh $RemoteHost 'sudo nginx -t && sudo systemctl reload nginx'" -ForegroundColor White
Write-Host ""
