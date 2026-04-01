<#
.SYNOPSIS
    Renders nginx config templates from .env / .env.prod substituting __PLACEHOLDER__ tokens.

.DESCRIPTION
    Reads env file, substitutes all __PLACEHOLDER__ tokens in conf.d/*.conf templates,
    writes rendered files to dist/conf.d/ (default) or -OutDir.
    Rendered files are ready to copy to /etc/nginx/conf.d/ on the server.

.PARAMETER EnvFile
    Path to .env or .env.prod file. Auto-detected if not specified.

.PARAMETER OutDir
    Output directory for rendered configs. Defaults to <nginx-infra>/dist/conf.d/

.EXAMPLE
    # Auto-detect .env.prod from project root
    .\scripts\render-nginx-conf.ps1

.EXAMPLE
    # Explicit env file and output dir
    .\scripts\render-nginx-conf.ps1 -EnvFile ..\..\..\..\..\.env.prod -OutDir .\dist\conf.d
#>

param(
    [string]$EnvFile = "",
    [string]$OutDir  = ""
)

$ErrorActionPreference = "Stop"

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$NginxDir   = Split-Path -Parent $ScriptDir
$ProjectRoot = Resolve-Path (Join-Path $NginxDir "..\..\..\..")

# ── Resolve env file ──────────────────────────────────────────────────────────
if (-not $EnvFile) {
    foreach ($candidate in @(
        (Join-Path $ProjectRoot ".env.prod"),
        (Join-Path $ProjectRoot ".env"),
        (Join-Path $NginxDir   ".env.prod"),
        (Join-Path $NginxDir   ".env")
    )) {
        if (Test-Path $candidate) { $EnvFile = $candidate; break }
    }
}

# ── Load env vars ─────────────────────────────────────────────────────────────
$envVars = @{}
if ($EnvFile -and (Test-Path $EnvFile)) {
    Write-Host "📋 Loading env from: $EnvFile" -ForegroundColor Cyan
    Get-Content $EnvFile | Where-Object { $_ -match '^[A-Z_][A-Z0-9_]*=' -and $_ -notmatch '^\s*#' } | ForEach-Object {
        $parts = $_ -split '=', 2
        $key   = $parts[0].Trim()
        $val   = ($parts[1] -replace '\r', '').Trim().Trim('"').Trim("'")
        $envVars[$key] = $val
    }
} else {
    Write-Host "⚠  No .env file found — using built-in defaults." -ForegroundColor Yellow
}

# ── Built-in defaults ─────────────────────────────────────────────────────────
function Get-EnvVal([string]$key, [string]$default) {
    if ($envVars.ContainsKey($key) -and $envVars[$key]) { return $envVars[$key] }
    return $default
}

$substitutions = [ordered]@{
    "__RATE_LIMIT_LOGIN_PER_MINUTE__"    = Get-EnvVal "RATE_LIMIT_LOGIN_PER_MINUTE"    "5"
    "__RATE_LIMIT_PASSWORD_PER_MINUTE__" = Get-EnvVal "RATE_LIMIT_PASSWORD_PER_MINUTE" "3"
    "__RATE_LIMIT_PER_MINUTE__"          = Get-EnvVal "RATE_LIMIT_PER_MINUTE"          "100"
    "__RATE_LIMIT_PER_DAY__"             = Get-EnvVal "RATE_LIMIT_PER_DAY"             "1000"
    "__OTP_RATE_LIMIT_MAX__"             = Get-EnvVal "OTP_RATE_LIMIT_MAX"             "100"
    "__OTP_RATE_LIMIT_WINDOW_MINUTES__"  = Get-EnvVal "OTP_RATE_LIMIT_WINDOW_MINUTES"  "60"
}

# ── Output dir ────────────────────────────────────────────────────────────────
if (-not $OutDir) { $OutDir = Join-Path $NginxDir "dist\conf.d" }
if (-not (Test-Path $OutDir)) { New-Item -ItemType Directory -Path $OutDir -Force | Out-Null }

Write-Host ""
$loginRate    = $substitutions["__RATE_LIMIT_LOGIN_PER_MINUTE__"]
$passwordRate = $substitutions["__RATE_LIMIT_PASSWORD_PER_MINUTE__"]
$otpRate      = $substitutions["__OTP_RATE_LIMIT_MAX__"]
Write-Host "Rendering nginx conf.d templates -> $OutDir" -ForegroundColor Cyan
Write-Host "   login rate limit:    ${loginRate}r/m"
Write-Host "   password rate limit: ${passwordRate}r/m"
Write-Host "   OTP rate limit:      ${otpRate}/hour"
Write-Host ""

$rendered = 0
$copied   = 0

Get-ChildItem (Join-Path $NginxDir "conf.d") -Filter "*.conf" | ForEach-Object {
    $template = $_.FullName
    $dest     = Join-Path $OutDir $_.Name
    $content  = Get-Content $template -Raw -Encoding UTF8

    $placeholderPattern = "__[A-Z_]+__"
    $hasPlaceholder = $content -match $placeholderPattern

    if ($hasPlaceholder) {
        foreach ($token in $substitutions.Keys) {
            $content = $content -replace [regex]::Escape($token), $substitutions[$token]
        }

        # Check for unresolved placeholders
        $remaining = [regex]::Matches($content, $placeholderPattern) | ForEach-Object { $_.Value } | Select-Object -Unique
        if ($remaining) {
            $remainingStr = $remaining -join ", "
            Write-Host "  warning: $($_.Name): unresolved: $remainingStr" -ForegroundColor Yellow
        } else {
            Write-Host "  OK: $($_.Name) (rendered)" -ForegroundColor Green
        }
        $rendered++
    } else {
        Write-Host "  📄 $($_.Name) (copied)" -ForegroundColor Gray
        $copied++
    }

    [System.IO.File]::WriteAllText($dest, $content, [System.Text.Encoding]::UTF8)
}

Write-Host ""
Write-Host "Done. Rendered: $rendered  Copied: $copied" -ForegroundColor Green
Write-Host "Output: $OutDir" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  Review: Get-Content $OutDir\03-security.conf"
Write-Host "  Deploy: .\scripts\quickerdeploy.sh  (uses rendered output)"
Write-Host "  Local:  sudo cp dist/conf.d/*.conf /etc/nginx/conf.d/ && sudo nginx -t"
