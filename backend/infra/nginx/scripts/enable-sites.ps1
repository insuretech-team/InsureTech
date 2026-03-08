# Enable Nginx Sites Script (PowerShell)
# Creates symlinks from sites-available to sites-enabled

$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $PSScriptRoot
$SitesAvailable = Join-Path $ScriptDir "sites-available"
$SitesEnabled = Join-Path $ScriptDir "sites-enabled"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Nginx Sites Enabler" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Create sites-enabled directory if it doesn't exist
if (-not (Test-Path $SitesEnabled)) {
    Write-Host "Creating sites-enabled directory..." -ForegroundColor Yellow
    New-Item -ItemType Directory -Path $SitesEnabled -Force | Out-Null
}

# Function to enable a site
function Enable-Site {
    param([string]$SiteName)
    
    $Source = Join-Path $SitesAvailable $SiteName
    $Target = Join-Path $SitesEnabled $SiteName
    
    if (-not (Test-Path $Source)) {
        Write-Host "✗ Site configuration not found: $SiteName" -ForegroundColor Red
        return $false
    }
    
    if (Test-Path $Target) {
        Write-Host "✓ Site already enabled: $SiteName" -ForegroundColor Green
        return $true
    }
    
    Write-Host "→ Enabling site: $SiteName" -ForegroundColor Yellow
    
    # Create relative symlink
    $RelativePath = "..\sites-available\$SiteName"
    
    # For Docker/Linux compatibility, create a file that references the source
    # On Windows, try to create actual symlink (requires admin)
    try {
        New-Item -ItemType SymbolicLink -Path $Target -Target $Source -Force -ErrorAction Stop | Out-Null
        Write-Host "✓ Site enabled (symlink): $SiteName" -ForegroundColor Green
    } catch {
        # Fallback: Copy the file instead of symlinking
        Copy-Item -Path $Source -Destination $Target -Force
        Write-Host "✓ Site enabled (copy): $SiteName" -ForegroundColor Green
        Write-Host "  Note: Created copy instead of symlink (run as Administrator for symlinks)" -ForegroundColor Yellow
    }
    
    return $true
}

# Enable all sites
Write-Host "Enabling sites..." -ForegroundColor Yellow
Write-Host ""

$sites = @(
    "trendyco.com.bd.conf",
    "portal.trendyco.com.bd.conf",
    "mta-sts.trendyco.com.bd.conf",
    "default.conf"
)

$enabled = 0
foreach ($site in $sites) {
    if (Enable-Site -SiteName $site) {
        $enabled++
    }
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Summary" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Enabled sites: $enabled / $($sites.Count)" -ForegroundColor Green
Write-Host ""

Get-ChildItem -Path $SitesEnabled | ForEach-Object {
    Write-Host "  ✓ $($_.Name)" -ForegroundColor Green
}

Write-Host ""
Write-Host "✓ All sites enabled successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Deploy with Docker Compose:" -ForegroundColor White
Write-Host "   docker-compose -f infra/nginx/docker-compose.nginx.yml up -d" -ForegroundColor Gray
Write-Host "2. Or test configuration (if nginx installed locally):" -ForegroundColor White
Write-Host "   nginx -t" -ForegroundColor Gray
Write-Host "3. Check logs:" -ForegroundColor White
Write-Host "   docker logs -f trendico_nginx" -ForegroundColor Gray
