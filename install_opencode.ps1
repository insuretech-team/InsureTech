# Requires Run as Administrator
if (!([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Warning "This script needs to be run as Administrator!"
    Write-Warning "Please right-click PowerShell or Windows Terminal, select 'Run as Administrator', then execute this script again."
    pause
    exit
}

Write-Host "1. Uninstalling existing broken Chocolatey installation..." -ForegroundColor Cyan
choco uninstall opencode -y

Write-Host "`n2. Removing legacy standalone installation..." -ForegroundColor Cyan
$LegacyDir = "$env:LOCALAPPDATA\Programs\OpenCode"
if (Test-Path $LegacyDir) {
    # Stop any hanging instances before trying to delete
    Stop-Process -Name "opencode", "bun" -Force -ErrorAction SilentlyContinue
    Remove-Item -Recurse -Force $LegacyDir
    Write-Host "Removed legacy installation at $LegacyDir" -ForegroundColor Green
}

$InstallDir = "C:\_opencode"
Write-Host "`n3. Creating stable installation directory at $InstallDir..." -ForegroundColor Cyan
if (!(Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

Write-Host "`n4. Installing stable OpenCode via NPM..." -ForegroundColor Cyan
Set-Location $InstallDir
# Install the CLI package locally in the target folder
npm install --prefix $InstallDir opencode-ai

Write-Host "`n5. Adding $InstallDir to System PATH and Current Session..." -ForegroundColor Cyan
$NodeModulesBin = "$InstallDir\node_modules\.bin"
$CurrentMachinePath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::Machine)

if ($CurrentMachinePath -notmatch [regex]::Escape($NodeModulesBin)) {
    $NewMachinePath = $CurrentMachinePath + ";" + $NodeModulesBin
    [Environment]::SetEnvironmentVariable("PATH", $NewMachinePath, [EnvironmentVariableTarget]::Machine)
    Write-Host "Successfully added to system PATH." -ForegroundColor Green
} else {
    Write-Host "Directory already exists in system PATH." -ForegroundColor Yellow
}

# Also update the current session so a restart isn't strictly necessary.
if ($env:PATH -notmatch [regex]::Escape($NodeModulesBin)) {
    $env:PATH = "$env:PATH;$NodeModulesBin"
}

Write-Host "`nInstallation Complete!" -ForegroundColor Green
Write-Host "The current session has been updated. You can run 'opencode' immediately." -ForegroundColor Green
pause
