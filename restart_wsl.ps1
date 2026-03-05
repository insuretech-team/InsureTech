# Restart WSL Script
# Run this in PowerShell as Administrator if WSL has crashed

Write-Host "Restarting WSL..." -ForegroundColor Yellow

# Option 1: Restart WSL service
Write-Host "`n[Option 1] Restarting WSL service..." -ForegroundColor Cyan
wsl --shutdown
Start-Sleep -Seconds 3
Write-Host "WSL shutdown complete" -ForegroundColor Green

# Option 2: Restart specific distribution
Write-Host "`n[Option 2] Available distributions:" -ForegroundColor Cyan
wsl --list --verbose

Write-Host "`nTo start WSL again, run:" -ForegroundColor Yellow
Write-Host "  wsl" -ForegroundColor White
Write-Host "`nOr to start a specific distribution:" -ForegroundColor Yellow
Write-Host "  wsl -d <DistributionName>" -ForegroundColor White

Write-Host "`n✓ WSL has been shut down. Start it again with 'wsl' command." -ForegroundColor Green
