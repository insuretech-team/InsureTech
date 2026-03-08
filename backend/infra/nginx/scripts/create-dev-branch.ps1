# Create Dev Branch for Nginx Migration
# This script helps create and push the dev branch for testing

$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Nginx Migration - Dev Branch Creator" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Check if we're in the right directory
$currentDir = Get-Location
if (-not (Test-Path "infra\nginx\nginx.conf")) {
    Write-Host "ERROR: Please run this script from the Trendico root directory" -ForegroundColor Red
    Write-Host "Current directory: $currentDir" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ Running from correct directory" -ForegroundColor Green
Write-Host ""

# Check git status
Write-Host "Checking git status..." -ForegroundColor Yellow
$gitStatus = git status --porcelain
if ($gitStatus) {
    Write-Host ""
    Write-Host "Modified files:" -ForegroundColor Yellow
    git status --short
    Write-Host ""
    
    $response = Read-Host "You have uncommitted changes. Continue? (y/n)"
    if ($response -ne "y") {
        Write-Host "Aborted by user" -ForegroundColor Yellow
        exit 0
    }
}

Write-Host ""
Write-Host "Creating dev branch..." -ForegroundColor Yellow

# Create and checkout dev branch
try {
    git checkout -b dev/nginx-modular
    Write-Host "✓ Created branch: dev/nginx-modular" -ForegroundColor Green
} catch {
    Write-Host "Branch might already exist, checking out..." -ForegroundColor Yellow
    git checkout dev/nginx-modular
}

Write-Host ""
Write-Host "Staging changes..." -ForegroundColor Yellow

# Stage all nginx-related changes
git add infra/nginx/
git add docker-compose.prod.yml
git add .github/workflows/build_and_push.yml
git add .github/workflows/deploy.yml
git add .github/workflows/test-nginx.yml

Write-Host "✓ Changes staged" -ForegroundColor Green
Write-Host ""

# Show what will be committed
Write-Host "Files to be committed:" -ForegroundColor Cyan
git diff --cached --name-status
Write-Host ""

$response = Read-Host "Proceed with commit? (y/n)"
if ($response -ne "y") {
    Write-Host "Aborted by user" -ForegroundColor Yellow
    exit 0
}

# Commit changes
Write-Host ""
Write-Host "Committing changes..." -ForegroundColor Yellow

$commitMessage = @"
feat: Add modular nginx architecture with automated testing

- Migrate from monolithic nginx.conf to modular architecture
- Add 30+ organized configuration files
- Create Dockerfile for nginx with Alpine base
- Implement three-tier caching strategy (static, API, microcache)
- Add comprehensive automated testing workflow
- Update docker-compose.prod.yml to use new nginx image
- Update CI/CD workflows (build_and_push.yml, deploy.yml)
- Add extensive documentation (MIGRATION_GUIDE.md)
- Maintain backward compatibility with existing setup
- Include error pages in Docker image

Breaking Changes: None - Full backward compatibility maintained

Related Files:
- infra/nginx/ - Complete modular nginx configuration
- .github/workflows/test-nginx.yml - Automated testing
- docker-compose.prod.yml - Updated nginx service
- MIGRATION_GUIDE.md - Complete migration documentation
"@

git commit -m "$commitMessage"
Write-Host "✓ Changes committed" -ForegroundColor Green
Write-Host ""

# Push to remote
Write-Host "Pushing to remote..." -ForegroundColor Yellow
$response = Read-Host "Push branch to remote? (y/n)"
if ($response -eq "y") {
    try {
        git push -u origin dev/nginx-modular
        Write-Host "✓ Branch pushed to remote" -ForegroundColor Green
    } catch {
        Write-Host "Warning: Push failed. You may need to push manually:" -ForegroundColor Yellow
        Write-Host "  git push -u origin dev/nginx-modular" -ForegroundColor Gray
    }
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Dev Branch Created Successfully!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Branch: dev/nginx-modular" -ForegroundColor White
Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Yellow
Write-Host "1. Push to remote (if not done): git push -u origin dev/nginx-modular" -ForegroundColor White
Write-Host "2. GitHub Actions will automatically:" -ForegroundColor White
Write-Host "   - Run nginx tests (Ubuntu & Alpine)" -ForegroundColor Gray
Write-Host "   - Perform security scans" -ForegroundColor Gray
Write-Host "   - Run integration tests" -ForegroundColor Gray
Write-Host "   - Build and push Docker image" -ForegroundColor Gray
Write-Host "3. Check workflow results at:" -ForegroundColor White
Write-Host "   https://github.com/$((git remote get-url origin) -replace '.*github.com[:/](.+?)(.git)?$', '$1')/actions" -ForegroundColor Gray
Write-Host "4. Once tests pass, deploy to dev server" -ForegroundColor White
Write-Host "5. Monitor for 24 hours" -ForegroundColor White
Write-Host "6. Create PR to merge to main" -ForegroundColor White
Write-Host ""
Write-Host "Documentation:" -ForegroundColor Yellow
Write-Host "- MIGRATION_GUIDE.md - Complete migration steps" -ForegroundColor White
Write-Host "- DEPLOYMENT_CHECKLIST.md - Deployment procedure" -ForegroundColor White
Write-Host ""
