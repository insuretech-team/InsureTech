# Nginx Configuration Validation Script
# Validates all nginx configuration files and checks for completeness

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Nginx Configuration Validator" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

$baseDir = Split-Path -Parent $PSScriptRoot
$errors = @()
$warnings = @()
$passed = 0
$failed = 0

function Test-FileExists {
    param([string]$Path, [string]$Description)
    
    if (Test-Path $Path) {
        Write-Host "✓ $Description" -ForegroundColor Green
        $script:passed++
        return $true
    } else {
        Write-Host "✗ $Description - NOT FOUND" -ForegroundColor Red
        $script:errors += "$Description not found at $Path"
        $script:failed++
        return $false
    }
}

function Test-FileContent {
    param([string]$Path, [string]$Pattern, [string]$Description)
    
    if (Test-Path $Path) {
        $content = Get-Content $Path -Raw
        if ($content -match $Pattern) {
            Write-Host "✓ $Description" -ForegroundColor Green
            $script:passed++
            return $true
        } else {
            Write-Host "✗ $Description - PATTERN NOT FOUND" -ForegroundColor Red
            $script:errors += "$Description validation failed in $Path"
            $script:failed++
            return $false
        }
    } else {
        Write-Host "✗ $Description - FILE NOT FOUND" -ForegroundColor Red
        $script:errors += "File not found: $Path"
        $script:failed++
        return $false
    }
}

# Test Phase 1: Foundation
Write-Host "`nPhase 1: Foundation Files" -ForegroundColor Yellow
Write-Host "-------------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\nginx.conf" "Main nginx.conf"
Test-FileExists "$baseDir\conf.d\00-global.conf" "Global configuration"
Test-FileExists "$baseDir\conf.d\01-logging.conf" "Logging configuration"
Test-FileExists "$baseDir\conf.d\02-performance.conf" "Performance configuration"
Test-FileExists "$baseDir\conf.d\03-security.conf" "Security configuration"
Test-FileExists "$baseDir\conf.d\04-compression.conf" "Compression configuration"

# Test Phase 2: Reusable Components
Write-Host "`nPhase 2: Reusable Components" -ForegroundColor Yellow
Write-Host "----------------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\snippets\ssl-params.conf" "SSL parameters snippet"
Test-FileExists "$baseDir\snippets\proxy-params.conf" "Proxy parameters snippet"
Test-FileExists "$baseDir\snippets\security-headers.conf" "Security headers snippet"
Test-FileExists "$baseDir\snippets\websocket.conf" "WebSocket snippet"

# Test Phase 3: Backend Configuration
Write-Host "`nPhase 3: Backend Configuration" -ForegroundColor Yellow
Write-Host "------------------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\upstreams\gateway.conf" "API Gateway upstream"
Test-FileExists "$baseDir\upstreams\trendyco.conf" "Trendyco upstream"
Test-FileExists "$baseDir\upstreams\trendfront.conf" "Trendfront upstream"

# Test Phase 4: Caching Layer
Write-Host "`nPhase 4: Caching Layer" -ForegroundColor Yellow
Write-Host "----------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\cache\cache-zones.conf" "Cache zones configuration"
Test-FileExists "$baseDir\cache\cache-bypass.conf" "Cache bypass rules"
Test-FileContent "$baseDir\cache\cache-zones.conf" "static_cache" "Static cache zone defined"
Test-FileContent "$baseDir\cache\cache-zones.conf" "api_cache" "API cache zone defined"
Test-FileContent "$baseDir\cache\cache-zones.conf" "microcache" "Microcache zone defined"

# Test Phase 5: Virtual Hosts
Write-Host "`nPhase 5: Virtual Hosts" -ForegroundColor Yellow
Write-Host "----------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\sites-available\trendyco.com.bd.conf" "Trendyco site config"
Test-FileExists "$baseDir\sites-available\portal.trendyco.com.bd.conf" "Portal site config"
Test-FileExists "$baseDir\sites-available\mta-sts.trendyco.com.bd.conf" "MTA-STS site config"
Test-FileExists "$baseDir\sites-available\default.conf" "Default catch-all config"

# Test Phase 6: Additional Features
Write-Host "`nPhase 6: Additional Features" -ForegroundColor Yellow
Write-Host "----------------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\maps\bot-detection.conf" "Bot detection maps"
Test-FileExists "$baseDir\error-pages\404.html" "404 error page"
Test-FileExists "$baseDir\error-pages\502.html" "502 error page"
Test-FileContent "$baseDir\conf.d\03-security.conf" "limit_req_zone" "Rate limiting configured"

# Test Phase 7: Automation & Deployment
Write-Host "`nPhase 7: Automation & Deployment" -ForegroundColor Yellow
Write-Host "---------------------------------" -ForegroundColor Yellow
Test-FileExists "$baseDir\docker-compose.nginx.yml" "Docker Compose configuration"
Test-FileExists "$baseDir\scripts\setup.sh" "Setup script"
Test-FileExists "$baseDir\scripts\test-config.sh" "Test script"
Test-FileExists "$baseDir\scripts\clear-cache.sh" "Cache management script"
Test-FileExists "$baseDir\README.md" "README documentation"
Test-FileExists "$baseDir\DEPLOYMENT.md" "Deployment documentation"
Test-FileExists "$baseDir\IMPLEMENTATION_STATUS.md" "Implementation status"

# Test Content Validation
Write-Host "`nContent Validation" -ForegroundColor Yellow
Write-Host "------------------" -ForegroundColor Yellow
Test-FileContent "$baseDir\nginx.conf" "include /etc/nginx/conf.d/\*\.conf" "Main config includes conf.d"
Test-FileContent "$baseDir\nginx.conf" "include /etc/nginx/upstreams/\*\.conf" "Main config includes upstreams"
Test-FileContent "$baseDir\nginx.conf" "include /etc/nginx/cache/\*\.conf" "Main config includes cache"
Test-FileContent "$baseDir\nginx.conf" "include /etc/nginx/sites-enabled/\*\.conf" "Main config includes sites"

Test-FileContent "$baseDir\sites-available\trendyco.com.bd.conf" "proxy_cache static_cache" "Static caching enabled"
Test-FileContent "$baseDir\sites-available\trendyco.com.bd.conf" "proxy_cache api_cache" "API caching enabled"
Test-FileContent "$baseDir\sites-available\trendyco.com.bd.conf" "proxy_cache microcache" "Microcaching enabled"
Test-FileContent "$baseDir\sites-available\trendyco.com.bd.conf" "include snippets/ssl-params.conf" "SSL params included"
Test-FileContent "$baseDir\sites-available\trendyco.com.bd.conf" "include snippets/security-headers.conf" "Security headers included"

# Summary
Write-Host "`n==================================" -ForegroundColor Cyan
Write-Host "Validation Summary" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
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
    Write-Host "✓ All validation checks passed!" -ForegroundColor Green
    Write-Host "Configuration is ready for deployment." -ForegroundColor Green
    exit 0
} else {
    Write-Host "✗ Validation failed with $failed errors." -ForegroundColor Red
    Write-Host "Please fix the errors before deployment." -ForegroundColor Red
    exit 1
}
