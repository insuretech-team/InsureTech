# =====================================================
# Safe Database Reset (Primary/Backup)
# =====================================================
# Deletes all tables in public schema, then drops all other schemas.
# Keeps only empty public schema and system schemas.
# =====================================================

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("primary","backup","both")]
    [string]$Target = "both",

    [switch]$DryRun,
    [switch]$SkipConfirmation
)

$ErrorActionPreference = "Stop"

# SQL to drop all tables in public schema, then drop all custom schemas
$nukeSQL = @"
DO `$`$
DECLARE
    r RECORD;
BEGIN
    -- Drop all tables in public schema
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS public.' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
    
    -- Drop all non-system schemas (keeping public, pg_catalog, information_schema, pg_toast)
    FOR r IN (
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('public', 'pg_catalog', 'information_schema', 'pg_toast')
        AND schema_name NOT LIKE 'pg_%'
    ) LOOP
        EXECUTE 'DROP SCHEMA IF EXISTS ' || quote_ident(r.schema_name) || ' CASCADE';
    END LOOP;
END
`$`$;
"@

Write-Host "" 
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "SAFE RESET (DROP CUSTOM SCHEMAS ONLY)" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Target: $Target" -ForegroundColor Cyan
Write-Host "DryRun: $DryRun" -ForegroundColor Cyan
Write-Host "" 

if (-not $SkipConfirmation -and -not $DryRun) {
  $confirmation = Read-Host "Type 'RESET' to confirm"
  if ($confirmation -ne "RESET") {
    Write-Host "Aborted." -ForegroundColor Red
    exit 0
  }
}

if ($DryRun) {
  Write-Host "--- SQL Preview ---" -ForegroundColor Gray
  Write-Host $nukeSQL -ForegroundColor Gray
  exit 0
}

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location "$root\backend\inscore"

Write-Host "Resetting $Target..." -ForegroundColor Yellow

# Save SQL to temp file to avoid escaping issues
$tempSqlFile = [System.IO.Path]::GetTempFileName() + ".sql"
$nukeSQL | Out-File -FilePath $tempSqlFile -Encoding UTF8 -NoNewline

# Keep Go compile memory low on constrained Windows hosts.
$prevGoFlags = $env:GOFLAGS
$prevGoMaxProcs = $env:GOMAXPROCS

try {
  if ([string]::IsNullOrWhiteSpace($env:GOFLAGS)) {
    $env:GOFLAGS = "-p=1"
  } elseif ($env:GOFLAGS -notmatch "(^|\s)-p=1(\s|$)") {
    $env:GOFLAGS = "$($env:GOFLAGS) -p=1"
  }
  $env:GOMAXPROCS = "1"

  # Use lightweight SQL runner (avoids heavy dbmanager compile graph).
  go run ./cmd/dbsql --sql-file="$tempSqlFile" --target=$Target

  if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Reset failed for $Target" -ForegroundColor Red
    Write-Host "Exit code: $LASTEXITCODE" -ForegroundColor Red
    throw "Reset failed for $Target"
  }

  Write-Host "Reset complete for $Target" -ForegroundColor Green
}
finally {
  # Restore environment
  if ($null -eq $prevGoFlags) { Remove-Item Env:GOFLAGS -ErrorAction SilentlyContinue } else { $env:GOFLAGS = $prevGoFlags }
  if ($null -eq $prevGoMaxProcs) { Remove-Item Env:GOMAXPROCS -ErrorAction SilentlyContinue } else { $env:GOMAXPROCS = $prevGoMaxProcs }

  # Clean up temp file
  if (Test-Path $tempSqlFile) {
    Remove-Item $tempSqlFile -Force
  }
}
