# Pre-Deployment Check Script
# Validates everything before deploying to dev/production

$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Nginx Pre-Deployment Check" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

$baseDir = "C:\_DEV\GO\Trendico"
$nginxDir = "$baseDir\infra\nginx"
$errors = @()
$warnings = @()
$passed = 0
$failed = 0

function Test-Check {
    param([string]$Description, [scriptblock]$Test)
    
    try {
        $result = & $Test
        if ($result) {
            Write-Host "✓ $Description" -ForegroundColor Green
            $script:passed++
            return $true
        } else {
            Write-Host "✗ $Description - FAILED" -ForegroundColor Red
            $script:errors += $Description
            $script:failed++
            return $false
        }
    } catch {
        Write-Host "✗ $Description - ERROR: $_" -ForegroundColor Red
        $script:errors += "$Description - $($_.Exception.Message)"
        $script:failed++
        return $false
    }
}

# Check 1: Working Directory
Write-Host "Environment Checks" -ForegroundColor Yellow
Write-Host "-----------------" -ForegroundColor Yellow
Test-Check "Working in Trendico root directory" {
    (Get-Location).Path -eq $baseDir -or (Test-Path "infra\nginx\nginx.conf")
}

Test-Check "Git repository initialized" {
    Test-Path ".git"
}

Test-Check "Docker installed and running" {
    try {
        docker ps > $null 2>&1
        return $true
    } catch {
        return $false
    }
}

# Check 2: Nginx Configuration Files
Write-Host "`nConfiguration Files" -ForegroundColor Yellow
Write-Host "-------------------" -ForegroundColor Yellow

$requiredFiles = @(
    "infra\nginx\nginx.conf",
    "infra\nginx\Dockerfile",
    "infra\nginx\conf.d\00-global.conf",
    "infra\nginx\conf.d\01-logging.conf",
    "infra\nginx\conf.d\02-performance.conf",
    "infra\nginx\conf.d\03-security.conf",
    "infra\nginx\conf.d\04-compression.conf",
    "infra\nginx\snippets\ssl-params.conf",
    "infra\nginx\snippets\proxy-params.conf",
    "infra\nginx\snippets\security-headers.conf",
    "infra\nginx\snippets\websocket.conf",
    "infra\nginx\upstreams\gateway.conf",
    "infra\nginx\upstreams\trendyco.conf",
    "infra\nginx\upstreams\trendfront.conf",
    "infra\nginx\cache\cache-zones.conf",
    "infra\nginx\cache\cache-bypass.conf",
    "infra\nginx\maps\bot-detection.conf",
    "infra\nginx\sites-enabled\trendyco.com.bd.conf",
    "infra\nginx\sites-enabled\portal.trendyco.com.bd.conf",
    "infra\nginx\sites-enabled\mta-sts.trendyco.com.bd.conf",
    "infra\nginx\sites-enabled\default.conf",
    "infra\nginx\error-pages\404.html",
    "infra\nginx\error-pages\502.html"
)

foreach ($file in $requiredFiles) {
    Test-Check "File exists: $file" {
        Test-Path $file
    }
}

# Check 3: Docker Compose Files
Write-Host "`nDocker Configuration" -ForegroundColor Yellow
Write-Host "--------------------" -ForegroundColor Yellow

Test-Check "docker-compose.prod.yml exists" {
    Test-Path "docker-compose.prod.yml"
}

Test-Check "docker-compose.prod.yml uses new nginx image" {
    $content = Get-Content "docker-compose.prod.yml" -Raw
    $content -match 'ghcr.io/\$\{GHCR_USERNAME\}/trendico-nginx:latest'
}

Test-Check "nginx cache volumes defined" {
    $content = Get-Content "docker-compose.prod.yml" -Raw
    ($content -match 'nginx_cache:') -and ($content -match 'nginx_logs:')
}

# Check 4: CI/CD Workflows
Write-Host "`nCI/CD Workflows" -ForegroundColor Yellow
Write-Host "---------------" -ForegroundColor Yellow

Test-Check "test-nginx.yml exists" {
    Test-Path ".github\workflows\test-nginx.yml"
}

Test-Check "build_and_push.yml updated" {
    $content = Get-Content ".github\workflows\build_and_push.yml" -Raw
    $content -match 'nginx:' -and $content -match 'trendico-nginx'
}

Test-Check "deploy.yml updated" {
    $content = Get-Content ".github\workflows\deploy.yml" -Raw
    $content -match 'Nginx config is now baked into the Docker image'
}

# Check 5: Documentation
Write-Host "`nDocumentation" -ForegroundColor Yellow
Write-Host "-------------" -ForegroundColor Yellow

$requiredDocs = @(
    "infra\nginx\README.md",
    "infra\nginx\DEPLOYMENT.md",
    "infra\nginx\DEPLOYMENT_CHECKLIST.md",
    "infra\nginx\QUICKSTART.md",
    "infra\nginx\IMPLEMENTATION_STATUS.md",
    "infra\nginx\PROGRESS_REPORT.md",
    "infra\nginx\MIGRATION_GUIDE.md",
    "infra\nginx\CONFIGURATION_COMPARISON.md",
    "infra\nginx\FINAL_SUMMARY.md"
)

foreach ($doc in $requiredDocs) {
    Test-Check "Documentation: $(Split-Path $doc -Leaf)" {
        Test-Path $doc
    }
}

# Check 6: Scripts
Write-Host "`nScripts" -ForegroundColor Yellow
Write-Host "-------" -ForegroundColor Yellow

$requiredScripts = @(
    "infra\nginx\scripts\validate-config.ps1",
    "infra\nginx\scripts\enable-sites.ps1",
    "infra\nginx\scripts\enable-sites.sh",
    "infra\nginx\scripts\create-dev-branch.ps1",
    "infra\nginx\scripts\setup.sh",
    "infra\nginx\scripts\test-config.sh",
    "infra\nginx\scripts\clear-cache.sh"
)

foreach ($script in $requiredScripts) {
    Test-Check "Script: $(Split-Path $script -Leaf)" {
        Test-Path $script
    }
}

# Check 7: Container Name Compatibility
Write-Host "`nBackward Compatibility" -ForegroundColor Yellow
Write-Host "----------------------" -ForegroundColor Yellow

Test-Check "Gateway upstream uses correct container name" {
    $content = Get-Content "infra\nginx\upstreams\gateway.conf" -Raw
    $content -match 'server trendico_gateway:8080'
}

Test-Check "Trendyco upstream uses correct container name" {
    $content = Get-Content "infra\nginx\upstreams\trendyco.conf" -Raw
    $content -match 'server trendico_trendyco:3000'
}

Test-Check "Trendfront upstream uses correct container name" {
    $content = Get-Content "infra\nginx\upstreams\trendfront.conf" -Raw
    $content -match 'server trendico_trendfront:3001'
}

# Check 8: Docker Build Test (Optional)
Write-Host "`nDocker Build Test (Optional)" -ForegroundColor Yellow
Write-Host "----------------------------" -ForegroundColor Yellow

$buildTest = Read-Host "Run Docker build test? This will build the image locally (y/n)"
if ($buildTest -eq "y") {
    Test-Check "Docker image builds successfully" {
        try {
            Write-Host "  Building nginx image..." -ForegroundColor Gray
            docker build -f infra\nginx\Dockerfile -t trendico-nginx:test . 2>&1 | Out-Null
            return $true
        } catch {
            return $false
        }
    }
    
    Test-Check "Nginx configuration in image is valid" {
        try {
            docker run --rm trendico-nginx:test nginx -t 2>&1 | Out-Null
            return $true
        } catch {
            return $false
        }
    }
    
    Write-Host "  Cleaning up test image..." -ForegroundColor Gray
    docker rmi trendico-nginx:test -f 2>&1 | Out-Null
}

# Summary
Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "Pre-Deployment Check Summary" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor $(if ($failed -eq 0) { "Green" } else { "Red" })

if ($errors.Count -gt 0) {
    Write-Host "`nErrors Found:" -ForegroundColor Red
    $errors | ForEach-Object { Write-Host "  - $_" -ForegroundColor Red }
}

if ($warnings.Count -gt 0) {
    Write-Host "`nWarnings:" -ForegroundColor Yellow
    $warnings | ForEach-Object { Write-Host "  - $_" -ForegroundColor Yellow }
}

Write-Host ""
if ($failed -eq 0) {
    Write-Host "✅ All pre-deployment checks passed!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Ready to proceed with deployment!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next Steps:" -ForegroundColor Yellow
    Write-Host "1. Run: .\infra\nginx\scripts\create-dev-branch.ps1" -ForegroundColor White
    Write-Host "2. Push to remote (triggers automated tests)" -ForegroundColor White
    Write-Host "3. Wait for GitHub Actions to complete" -ForegroundColor White
    Write-Host "4. Deploy to dev server" -ForegroundColor White
    Write-Host "5. Follow MIGRATION_GUIDE.md for complete steps" -ForegroundColor White
    Write-Host ""
    exit 0
} else {
    Write-Host "❌ Pre-deployment check failed with $failed errors." -ForegroundColor Red
    Write-Host "Please fix the errors before proceeding with deployment." -ForegroundColor Red
    Write-Host ""
    exit 1
}
