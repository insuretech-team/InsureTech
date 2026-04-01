# Normalize all shell scripts (*.sh) and env files (.env*) under the repo root to LF line endings.

param(
    [Parameter(Mandatory=$false)]
    [string]$Root,

    [switch]$DryRun,

    [switch]$IncludeDependencies
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($Root)) {
    $Root = $PSScriptRoot
}

$resolvedRoot = (Resolve-Path $Root).Path
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)

Write-Host ""
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "NORMALIZE LINE ENDINGS (*.sh + .env*)" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Root: $resolvedRoot" -ForegroundColor Cyan
Write-Host "DryRun: $DryRun" -ForegroundColor Cyan
Write-Host "IncludeDependencies: $IncludeDependencies" -ForegroundColor Cyan
Write-Host ""

$files = Get-ChildItem -Path $resolvedRoot -Include @("*.sh", ".env*") -File -Recurse |
    Where-Object {
        $IncludeDependencies -or (
            $_.FullName -notmatch '[\\/]\.git[\\/]' -and
            $_.FullName -notmatch '[\\/]node_modules[\\/]'
        )
    } |
    Sort-Object FullName

if (-not $files) {
    Write-Host "No .sh or .env* files found." -ForegroundColor Yellow
    exit 0
}

$normalizedCount = 0
$alreadyLfCount = 0

foreach ($file in $files) {
    $content = Get-Content -Raw -Path $file.FullName

    if (-not $content.Contains("`r`n")) {
        $alreadyLfCount++
        Write-Host "[OK]   $($file.FullName)" -ForegroundColor DarkGray
        continue
    }

    $normalized = $content -replace "`r`n", "`n"

    if ($DryRun) {
        Write-Host "[PLAN] $($file.FullName)" -ForegroundColor Yellow
    } else {
        [System.IO.File]::WriteAllText($file.FullName, $normalized, $utf8NoBom)
        Write-Host "[FIX]  $($file.FullName)" -ForegroundColor Green
    }

    $normalizedCount++
}

Write-Host ""
Write-Host "Scanned: $($files.Count)" -ForegroundColor Cyan
Write-Host "Normalized: $normalizedCount" -ForegroundColor Green
Write-Host "Already LF: $alreadyLfCount" -ForegroundColor DarkGray
