$ErrorActionPreference = "SilentlyContinue"

Write-Host "=============================" -ForegroundColor Cyan
Write-Host "   Disk Space Cleanup        " -ForegroundColor Cyan
Write-Host "=============================" -ForegroundColor Cyan

$objectsToDelete = @(
    # Android
    "$env:USERPROFILE\.android",
    "$env:USERPROFILE\AppData\Local\Android",
    "$env:USERPROFILE\AndroidStudioProjects",
    
    # Gradle
    "$env:USERPROFILE\.gradle",

    # Go
    "$env:USERPROFILE\AppData\Local\go-build",
    "$env:USERPROFILE\go\pkg\mod", # Only delete the module cache, not user source files

    # NPM
    "$env:USERPROFILE\AppData\Local\npm-cache"
)

$totalFreedBytes = 0

foreach ($obj in $objectsToDelete) {
    if (Test-Path $obj) {
        Write-Host "Measuring $obj..." -NoNewline
        $size = (Get-ChildItem -Path $obj -Recurse -Force -ErrorAction SilentlyContinue | Measure-Object -Property Length -Sum).Sum
        $sizeGB = [math]::Round($size / 1GB, 2)
        Write-Host " ($sizeGB GB)"

        Write-Host "Deleting $obj..." -NoNewline
        Remove-Item -Path $obj -Recurse -Force -ErrorAction SilentlyContinue
        
        # Verify deletion
        if (-not (Test-Path $obj)) {
            $totalFreedBytes += $size
            Write-Host " Done." -ForegroundColor Green
        }
        else {
            Write-Host " Partially failed (some files may be in use)." -ForegroundColor Yellow
        }
    }
    else {
        Write-Host "$obj not found. Skipping." -ForegroundColor DarkGray
    }
}

$totalFreedGB = [math]::Round($totalFreedBytes / 1GB, 2)
Write-Host "`nCleanup Complete. Successfully freed at least $totalFreedGB GB." -ForegroundColor Cyan
