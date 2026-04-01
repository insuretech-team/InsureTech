# ==============================================================================
# reset_kyc_verified.ps1
# Resets kyc_verified = false for B2B admin and system admin users
# so they go through the eKYC flow again on next login.
#
# Uses dbx.exe sql — the InsureTech internal DB CLI (no psql required).
# DB connection is auto-resolved from .env at the project root (go.mod root).
#
# Usage (run from project root E:\Projects\InsureTech):
#   .\scripts\reset_kyc_verified.ps1                        # primary DB
#   .\scripts\reset_kyc_verified.ps1 -Target backup         # Neon backup DB
#   .\scripts\reset_kyc_verified.ps1 -Target both           # both DBs
#   .\scripts\reset_kyc_verified.ps1 -DryRun                # SELECT only
# ==============================================================================

param(
    [ValidateSet("primary", "backup", "both")]
    [string]$Target = "primary",

    [switch]$DryRun
)

# ── Resolve project root = directory containing go.mod ───────────────────────
$projectRoot = $PSScriptRoot
while ($projectRoot -ne [System.IO.Path]::GetPathRoot($projectRoot)) {
    if (Test-Path (Join-Path $projectRoot "go.mod")) { break }
    $projectRoot = Split-Path $projectRoot -Parent
}

# ── Load admin emails — try .env first (project root), then .env.dev ─────────
function Read-EnvFile([string]$path) {
    $vars = @{}
    if (Test-Path $path) {
        Get-Content $path | ForEach-Object {
            if ($_ -match "^\s*([A-Za-z_][A-Za-z0-9_]*)=(.*)$") {
                $vars[$matches[1]] = $matches[2].Trim().Trim("'").Trim('"')
            }
        }
    }
    return $vars
}

$env1 = Read-EnvFile (Join-Path $projectRoot ".env")
$env2 = Read-EnvFile (Join-Path $projectRoot ".env.dev")

# Merge: .env takes priority for DB, .env.dev as fallback for missing keys
$merged = $env2.Clone()
foreach ($k in $env1.Keys) { $merged[$k] = $env1[$k] }

$ADMIN_EMAIL = $merged["ADMIN_EMAIL"] -replace '"','' -replace "'",''
$B2B_ADMIN   = $merged["B2B_ADMIN"]  -replace '"','' -replace "'",''
if (-not $ADMIN_EMAIL) { $ADMIN_EMAIL = "faruk.hannan@gmail.com" }
if (-not $B2B_ADMIN)   { $B2B_ADMIN  = "faruk.hannan@lifeplusbd.com" }

Write-Host ""
Write-Host "👤 Target users:" -ForegroundColor Cyan
Write-Host "   System Admin : $ADMIN_EMAIL"
Write-Host "   B2B Admin    : $B2B_ADMIN"
Write-Host "   DB target    : $Target"
Write-Host "   Project root : $projectRoot"
if ($DryRun) { Write-Host "   Mode         : DRY RUN (SELECT only)" -ForegroundColor Yellow }
Write-Host ""

# ── dbx path ─────────────────────────────────────────────────────────────────
$dbxExe = Join-Path $projectRoot "dbx.exe"

function Invoke-Dbsql {
    param([string]$Sql, [string]$Label)
    Write-Host "⚡ $Label" -ForegroundColor Green

    # Write SQL to a temp file — both dbx and dbsql support --sql-file
    $tmp = [System.IO.Path]::GetTempFileName() -replace '\.tmp$', '.sql'
    $Sql | Set-Content -Path $tmp -Encoding UTF8

    Push-Location $projectRoot
    try {
        if (Test-Path $dbxExe) {
            & $dbxExe sql --sql-file $tmp --target $Target
        } else {
            Write-Host "   (dbx.exe not found, using go run dbsql)" -ForegroundColor DarkGray
            & go run ./backend/inscore/cmd/dbsql --sql-file $tmp --target $Target
        }
    } finally {
        Pop-Location
        # Give the process a moment to release the file handle, then clean up
        $null = Start-Job { param($f) Start-Sleep 1; Remove-Item $f -Force -ErrorAction SilentlyContinue } -ArgumentList $tmp
    }
}

# ── Step 1: Show current state ─────────────────────────────────────────────────
$selectSql = @"
SELECT DISTINCT ON (u.user_id)
    u.email,
    u.user_id,
    up.kyc_verified,
    up.kyc_verified_at,
    COALESCE(kv.status, 'no record') AS kyc_record_status
FROM authn_schema.users u
JOIN authn_schema.user_profiles up ON up.user_id = u.user_id
LEFT JOIN authn_schema.kyc_verifications kv
    ON kv.entity_id = u.user_id
    AND kv.entity_type = 'user'
WHERE u.email IN ('$ADMIN_EMAIL', '$B2B_ADMIN')
ORDER BY u.user_id, kv.verified_at DESC NULLS LAST;
"@

Invoke-Dbsql -Sql $selectSql -Label "Current KYC state"

if ($DryRun) {
    Write-Host ""
    Write-Host "ℹ️  Dry run — no changes made." -ForegroundColor Yellow
    exit 0
}

# ── Step 2: Reset user_profiles.kyc_verified ──────────────────────────────────
$resetProfileSql = @"
UPDATE authn_schema.user_profiles up
SET
    kyc_verified    = false,
    kyc_verified_at = NULL,
    updated_at      = NOW()
FROM authn_schema.users u
WHERE u.user_id = up.user_id
  AND u.email IN ('$ADMIN_EMAIL', '$B2B_ADMIN');
"@

Invoke-Dbsql -Sql $resetProfileSql -Label "Resetting user_profiles.kyc_verified → false"

# ── Step 3: Reset kyc_verifications record → IN_PROGRESS ──────────────────────
$resetKycSql = @"
UPDATE authn_schema.kyc_verifications
SET
    status      = 'IN_PROGRESS',
    verified_at = NULL,
    verified_by = NULL
WHERE entity_id IN (
    SELECT user_id FROM authn_schema.users
    WHERE email IN ('$ADMIN_EMAIL', '$B2B_ADMIN')
)
AND entity_type = 'user';
"@

Invoke-Dbsql -Sql $resetKycSql -Label "Resetting kyc_verifications.status → IN_PROGRESS"

# ── Step 4: Confirm ───────────────────────────────────────────────────────────
Invoke-Dbsql -Sql $selectSql -Label "Confirmed state after reset"

Write-Host ""
Write-Host "✅ Done! kyc_verified reset to false for both users." -ForegroundColor Green
Write-Host "   → Log out and log back in to be routed to /kyc." -ForegroundColor Gray
Write-Host ""
